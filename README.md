# Fast-WireGuard

Quickly set up a WireGuard VPN server on your Linux machine with minimal configuration.

This repository is used just for learning purposes.

> [!WARNING]
> Current version can not uninstall Fast-WireGuard properly.
> Please wait for the next version.

## Quick Start
Run the following commands to install Fast-WireGuard:

```bash
wget https://github.com/MugeTong/fast-wireguard/releases/latest/download/fast-wireguard-Linux-amd64.sh
chmod +x fast-wireguard-Linux-amd64.sh
sudo ./fast-wireguard-Linux-amd64.sh
rm -f fast-wireguard-Linux-amd64.sh
```
After installation, you can use the `fwg` command to manage your WireGuard VPN server.

For example, to start the WireGuard server, run:
```bash
sudo fwg create
```

For more information on usage and configuration, refer to the documentation in the `docs` directory.

## Building from Source
To build Fast-WireGuard from source, ensure you have Go installed and run:
```bash
make build
```
This will create the binary in the `releases/bin` directory.
To create a self-extracting installer, run:
```bash
make pack
```
This will generate the installer script in the `releases` directory.
The installer can then be used to set up Fast-WireGuard on your system.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
