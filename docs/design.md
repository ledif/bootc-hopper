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

The hop command will essentially call `bootc switch` or `rpm-ostree rebase` to the new image with a prologue and epilogue.

Before the rebase starts:
- Copy the currently running `bootc-hopper` executable to `/var/lib/bootc-hopper/libexec/bootc-hopper`

After a rebase is finished:
- If the currently booted image is bootc-hopper aware, scripts in `/usr/lib/bootc-hopper/hop.d/` will be executed
- Create a file named `/var/lib/bootc-hopper/hop-state.yml` with information about the current deployment and future deployment

In addition, a systemd service will be created and enabled for next boot

```
# /etc/systemd/system/bootc-hopper-land.service
[Unit]
Exec=/var/opt/bootc-hopper/bin/bootc-hopper land
ConditionFileExists=/var/lib/bootc-hopper/hop-state.yml

[Install]
WantedBy=multi-user.target
```
and

```bash
systemctl enable bootc-hopper-land.service
```

### Land Implementation

The `land` command will do the following:
- If the new image is bootc-hopper aware, scripts in `/usr/lib/bootc-hopper/land.d/` will be executed
- If not, a new user will be created with the same passwd as the user that started the hop
- The `/var/lib/bootc-hopper/hop-state.yml` file will be removed
- The systemd service will be disabled
- The database at `/var/lib/bootc-hopper/history.sqlite` will be updated with information about the hop

## Image Repo
- Try to use Artifacthub if possible
- Define a remote for images called "Pond"
  - A git repo with a pond.yml in its root
  - Organization / Product Line / Product
    - For example, Universal Blue / Bazzite / bazzite-deck
