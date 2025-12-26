#!/bin/bash

DEPENDENCIES=("wg" "ip" "wg-quick")

# Define the Dependence list
sys::check_requirements() {
    for cmd in "${DEPENDENCIES[@]}"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            echo "Error: Package '$cmd' is not installed. Please install Wireguard by:"
            echo ""
            echo "\t sudo apt install wireguard"
            echo ""
            exit 1
        fi
    done
}
