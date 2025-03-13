> [!IMPORTANT]  
> This project is in pre-alpha. Most features (including the GUI) are not implemented.

A Flatpak GUI and command-line application that allows users to rebase to entirely different OS images.

### Usage

```bash
# List products that can be hopped to
$ bootc-hopper list-products --desktop-environment=KDE
● ublue-os/bazzite - "The next generation of Linux gaming" (4.9k stars)
    Source: https://github.com/ublue-os/bazzite

● ublue-os/aurora - "The ultimate productivity workstation" (144 stars)
    Source: https://github.com/ublue-os/aurora/

● winblues/blutiger - "Frutiger Aero in an atomic fashion" (1 star)
    Source: https://github.com/winblues/blutiger

# Hop to specific images
$ bootc-hopper hop oci://ghcr.io/ublue-os/bazzite:stable
$ bootc-hopper hop https://example.com/images/custom-image/Containerfile
$ bootc-hopper hop /path/to/Containerfile

# List images provided by products
$ bootc-hopper list-images ublue-os/aurora
- oci://ghcr.io/ublue-os/aurora:stable
- oci://ghcr.io/ublue-os/aurora-dx:stable

# Manage source of images
bootc-hopper remote add artifacthub https://artifacthub.io
bootc-hopper remote add acme git@github.com:acme/bootc-images.git

```

