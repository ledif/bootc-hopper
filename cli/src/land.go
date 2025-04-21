package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type LandConfig struct {
	// Add more options as needed
}

func land() error {
	// 1. Parse hop-state
	var state HopState
	if _, err := toml.DecodeFile("/var/lib/bootc-hopper/hop-state.toml", &state); err != nil {
		return fmt.Errorf("failed to parse hop-state: %w", err)
	}

	// 2. Parse land config (if needed)
	landConfig := LandConfig{}
	parseConfigFile("/usr/lib/bootc-hopper/land.conf", &landConfig)
	parseConfigDir("/usr/lib/bootc-hopper/land.conf.d", &landConfig)

	// 3. HomeStrategy handling
	switch state.HomeStrategy {
	case "NewUser":
		// Create new user same as old one, assume env USER
		username := os.Getenv("USER")
		if err := exec.Command("useradd", "-m", username).Run(); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	case "BestEffort":
		// Attempt to reset known config dirs (example)
		home := os.Getenv("HOME")
		dirs := []string{
			".config", ".local", ".kde", ".gnome2", ".cache",
		}
		for _, dir := range dirs {
			os.RemoveAll(filepath.Join(home, dir))
		}
	}

	// 4. Run scripts in land.d
	runScriptDir("/usr/lib/bootc-hopper/land.d")

	// 5. Cleanup hop state
	os.Remove("/var/lib/bootc-hopper/hop-state.toml")

	// 6. Disable systemd service
	exec.Command("systemctl", "disable", "bootc-hopper-land.service").Run()

	// 7. Update history DB (skipped here, needs sqlite setup)
	fmt.Println("Hop landed successfully.")
	return nil
}

