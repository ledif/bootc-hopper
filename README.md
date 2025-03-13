> [!IMPORTANT]  
> This project is in pre-alpha. Most features (including the GUI) are not implemented.

A Flatpak GUI and command-line application that allows users to rebase to entirely different OS images.

### Usage

```bash
# Hop to different images
bootc-hopper hop oci://ghcr.io/ublue-os/bazzite
bootc-hopper hop https://example.com/images/custom-image/Containerfile
bootc-hopper hop /path/to/Containerfile

# Fix up environment in new deployment (meant to be run automatically on next boot)
#  - will run scripts in /usr/share/bootc-hopper/land.d/*
#  - or generic fix up scripts if path does not exist
bootc-hopper land

# List images to hop to
bootc-hopper list

# Manage source of images
bootc-hopper remote add acme https://github.com/acme/bootc-images.git
```

### Hop Implementation
Maybe use https://systemd.io/PORTABLE_SERVICES/

The hop command will call `bootc switch` or `rpm-ostree rebase` to the new image.

After a rebase is finished, it will do the following:
- If the currently booted image is bootc-hopper aware, scripts in `/usr/share/bootc-hopper/hop.d/` will be executed
- A file named `/var/opt/bootc-hopper/var/hop-state.yml` will be created

In addition, a systemd service will be created and enabled for next boot
```bash
systemctl enable bootc-hopper-land.service
```

```
[Unit]
Exec=/var/opt/bootc-hopper/bin/bootc-hopper land
ConditionFileExists=/var/opt/bootc-hopper/var/hop-state.yml
```

The `land` command will do the following:
- If the new image is bootc-hopper aware, scripts in `/usr/share/bootc-hopper/land.d/` will be executed
- If not, a new user will be created with the same passwd as the user that started the hop
- The `/var/opt/bootc-hopper/var/hop-state.yml` file will be removed
- The systemd service will be disabled
