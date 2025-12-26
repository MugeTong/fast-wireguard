#!/bin/bash

# Check if the script is running with root privileges
core::ensure_root() {
    # EUID 0 is the root user
    if [[ ${EUID} -ne 0 ]]; then
        echo "Error: This command must be run as root." >&2
        echo "Please try again using 'sudo fwg <command>'" >&2
        exit 1
    fi
}

core::generate_keys() {
    local dir="/etc/wireguard"
    mkdir -p "$dir" && cd "$dir" || return 1
    umask 077

    # Generate the private key and save it
    local priv=$(wg genkey)
    echo "$priv" > privatekey

    # Generate the public key based on the private key
    local pub=$(echo "$priv" | wg pubkey)
    echo "$pub" > publickey

    # Return the keys
    echo "$priv $pub"
}

# Get the Ethernet name
core::get_eth() {
    ip route list default | awk '{for(i=1;i<=NF;i++) if($i=="dev") print $(i+1)}'
}

core::setup() {
    # Ensure the privilege
    core::ensure_root

    # Default port of WireGuard
    local port="51820"

    while [[ $# -gt 0 ]]; do
        case $1 in
            --port)
                if [[ -n $2 && $2 != -* ]]; then
                    port="$2"
                    shift 2
                else
                    echo "Error: --port requires a value." >&2
                    return 1
                fi
                ;;
            *)
                echo "Unknown option: $1"
                shift
                ;;
        esac
    done

    echo "Setting up WireGuard on port: $port"
    echo "Please ensure that you allow the port in your firewall."

    # Get the network interface name
    local eth=$(core::get_eth)
    # Unpack multiple values: get the generated keys
    local private_key public_key
    read -r private_key public_key <<< "$(core::generate_keys)"

    # Ask the user to enter the client public key
    local client_pubkey
    read -p "Enter the client public key: " client_pubkey
    # Basic validation to ensure it's not empty
    if [[ -z "$client_pubkey" ]]; then
        echo "Error: Client public key cannot be empty." >&2
        return 1
    fi

    # Generate the WireGuard configuration file
    local config_path="/etc/wireguard/wg0.conf"

    echo "Generating configuration at ${config_path}..."

    sudo tee "${config_path}" > /dev/null <<EOF
[Interface]
# Server internal IP
Address = 10.0.0.1/24
# Listening port
ListenPort = ${port}
# Server private key
PrivateKey = ${private_key}

# --- Core Network Forwarding Rules (iptables) ---
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o ${eth} -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o ${eth} -j MASQUERADE

[Peer]
# Client public key
PublicKey = ${client_pubkey}
# Client internal IP
AllowedIPs = 10.0.0.2/32
EOF

    echo "Configuration generated successfully."
    if ! sudo wg-quick up wg0; then
        echo "Failed to start WireGuard"
        exit 1
    fi

    # Display the client configuration
    local server_ip=$(curl -s https://api.ipify.org || curl -s https://ifconfig.me)
    local green="\033[0;32m"
    local cyan="\033[0;36m"
    local reset="\033[0m"

    echo -e "${green}---------------------------------------------------------------${reset}"
    echo -e "${green}Client Configuration (Copy the content below):${reset}"
    echo -e "${green}---------------------------------------------------------------${reset}"

    echo -e "${cyan}"
    cat <<EOF
[Interface]
# Client's private key (The one corresponding to the public key you entered)
PrivateKey = <INSERT_CLIENT_PRIVATE_KEY_HERE>
Address = 10.0.0.2/32
DNS = 8.8.8.8

[Peer]
# Server's public key
PublicKey = ${public_key}
# Server's public IP and port
Endpoint = ${server_ip}:${port}
# Forward all traffic through VPN
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
EOF
    echo -e "${reset}"

    echo -e "${green}---------------------------------------------------------------${reset}"
    echo "Please replace <INSERT_CLIENT_PRIVATE_KEY_HERE> with the actual private key."
    echo -e "${green}---------------------------------------------------------------${reset}"
}
