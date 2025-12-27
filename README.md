# Fast-WireGuard

Quickly set up a WireGuard VPN server on your Linux machine with minimal configuration.

This repository is used just for learning purposes.

> [!WARNING]
> Current version has lots of problem about the place to install and the uninstallation.
> Do not use this repo before the official release.

## Quick Start
```bash
wget https://github.com/MugeTong/fast-wireguard/releases/download/v0.1.0/fast-wireguard-0.1.0-Linux-x86_64.sh
chmod +x fast-wireguard-0.1.0-Linux-x86_64.sh
sudo ./fast-wireguard-0.1.0-Linux-x86_64.sh
rm -f ./fast-wireguard-0.1.0-Linux-x86_64.sh
sudo /root/.fast-wireguard/bin/fwg check
sudo /root/.fast-wireguard/bin/fwg setup
```


## Steps to Install
1. Select the latest release from the [Releases](https://github.com/MugeTong/fast-wireguard/releases) page.
2. run the following command in your terminal:
   ```bash
