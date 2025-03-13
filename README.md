> [!IMPORTANT]  
> This project is in pre-alpha. Most features (including the GUI) are not implemented.

A Flatpak GUI and command-line application that allows users to rebase to entirely different OS images.

### Usage

```bash
# Hop to different images
bootc-hopper hop oci://ghcr.io/ublue-os/bazzite
bootc-hopper hop https://example.com/images/custom-image/Containerfile
bootc-hopper hop /path/to/Containerfile

# List images to hop to
bootc-hopper list

# Manage source of images
bootc-hopper remote add acme https://github.com/acme/bootc-images.git
```

