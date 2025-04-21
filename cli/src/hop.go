package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type HopConfig struct {
	DesktopEnvironment string
	HomeStrategy       string
}

type HopState struct {
	Timestamp     time.Time
	SourceImage   string
	DesktopEnv    string
	HomeStrategy  string
}

func hop(imageRef string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable: %w", err)
	}
	libexecPath := "/var/lib/bootc-hopper/libexec/bootc-hopper"
	if err := os.MkdirAll(filepath.Dir(libexecPath), 0755); err != nil {
		return fmt.Errorf("failed to create libexec dir: %w", err)
	}
	src, err := os.Open(exePath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(libexecPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	if err := dst.Chmod(0755); err != nil {
		return err
	}

	// 2. Parse hop config
	config := HopConfig{}
	parseConfigFile("/usr/lib/bootc-hopper/hop.conf", &config)
	parseConfigDir("/usr/lib/bootc-hopper/hop.conf.d", &config)

	// 3. bootc switch
	cmd := exec.Command("bootc", "switch", imageRef)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("bootc switch failed: %w", err)
	}

	// 4. Run scripts in hop.d
	runScriptDir("/usr/lib/bootc-hopper/hop.d")

	// 5. Write hop-state.toml
	state := HopState{
		Timestamp:     time.Now(),
		SourceImage:   imageRef,
		DesktopEnv:    config.DesktopEnvironment,
		HomeStrategy:  config.HomeStrategy,
	}
	if err := os.MkdirAll("/var/lib/bootc-hopper", 0755); err != nil {
		return err
	}
	stateFile, err := os.Create("/var/lib/bootc-hopper/hop-state.toml")
	if err != nil {
		return err
	}
	defer stateFile.Close()
	if err := toml.NewEncoder(stateFile).Encode(state); err != nil {
		return err
	}

	// 6. Create and enable systemd service
	service := `[Unit]
Description=Bootc Hopper Landing
ConditionFileExists=/var/lib/bootc-hopper/hop-state.toml

[Service]
ExecStart=/var/lib/bootc-hopper/libexec/bootc-hopper land

[Install]
WantedBy=multi-user.target`
	if err := os.WriteFile("/etc/systemd/system/bootc-hopper-land.service", []byte(service), 0644); err != nil {
		return err
	}
	if err := exec.Command("systemctl", "enable", "bootc-hopper-land.service").Run(); err != nil {
		return fmt.Errorf("failed to enable land service: %w", err)
	}

	return nil
}
