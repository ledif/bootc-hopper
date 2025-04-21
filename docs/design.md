A Flatpak GUI and command-line application that allows users to rebase to entirely different OS images.
### Usage

```bash
# Hop to different images
bootc-hopper hop oci://ghcr.io/ublue-os/bazzite
bootc-hopper hop https://example.com/images/custom-image/Containerfile
bootc-hopper hop /path/to/Containerfile

# Fix up environment in new deployment (meant to be run automatically on next boot)
#  - will run scripts in /usr/lib/bootc-hopper/land.d/*
#  - or generic fix up scripts if path does not exist
bootc-hopper land

# List images to hop to
bootc-hopper list

# Manage source of images
bootc-hopper remote add acme https://github.com/acme/bootc-images.git
```

### Hop Implementation

The hop command will call `bootc switch` to the new image with a prologue and epilogue.

Before the rebase starts:
- Copy the currently running `bootc-hopper` executable to `/var/lib/bootc-hopper/libexec/bootc-hopper`
- Parse the configuration files in `/usr/lib/bootc-hopper/hop.conf` and `/usr/lib/bootc-hopper/hop.conf.d/` if they exist

Example configuration:
```toml
# The desktop environment of the booted image. Available options:
#   - autodetect, plasma, gnome, xfce
dektop_environment = "autodetect"
# How to handle DE-specific configurations in the user's home directory
# Options:
#  - new_user: Create a new user on boot into the new image
#  - best_effort: Attempt to remove DE-specific configs from the user's home directory
#  - none: Do not modify any of the user's files
migration_strategy = "new_user"
# Comma separated list of groups for whom all users in that group will be migrated
migration_group = "wheel"
# Comma separated list of users to migrate
migration_user = ""
```

After a rebase is finished:
- Execute scripts in `/usr/lib/bootc-hopper/hop.d/` if they exist
- Create a file named `/var/lib/bootc-hopper/hop-state.toml` with information about the current deployment and future deployment
  - Current deployment information will include parsed configuration from `hop.conf`/`hop.conf.d`.

In addition, a systemd service will be created and enabled for next boot

```
# /etc/systemd/system/bootc-hopper-land.service
[Unit]
Exec=/var/opt/bootc-hopper/bin/bootc-hopper land
ConditionFileExists=/var/lib/bootc-hopper/hop-state.toml

[Install]
WantedBy=multi-user.target
```
and

```bash
systemctl enable bootc-hopper-land.service
```

### Land Implementation

The `land` command will do the following:
- Parse the configuration files in `/usr/lib/bootc-hopper/land.conf` and `/usr/lib/bootc-hopper/land.conf.d/` if they exist
- Execute built-in functions based on combination of hop configuration (saved in `hop-state.toml`) and land configuration 
   - `HomeStrategy=NewUser`: a new user will be created with the same passwd as the user that started the hop
   - `HomeStrategy=BestEffort`: based on the desktop environment, try to remove all known configuration files in the directory of the user that started the hop
- Execute scripts in `/usr/lib/bootc-hopper/land.d/` if they exist
- The `/var/lib/bootc-hopper/hop-state.toml` file will be removed
- The systemd service will be disabled
- The database at `/var/lib/bootc-hopper/history.sqlite` will be updated with information about the hop

## Image Repo
- Try to use Artifacthub if possible
- Define a remote for images called "Pond"
  - A git repo with a pond.yml in its root
  - Organization / Product Line / Product
    - For example, Universal Blue / Bazzite / bazzite-deck
