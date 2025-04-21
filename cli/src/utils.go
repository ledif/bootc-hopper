package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func parseConfigFile(path string, out interface{}) {
	if _, err := os.Stat(path); err == nil {
		toml.DecodeFile(path, out)
	}
}

func parseConfigDir(dir string, out interface{}) {
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		if !f.IsDir() {
			toml.DecodeFile(filepath.Join(dir, f.Name()), out)
		}
	}
}

func runScriptDir(dir string) {
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			path := filepath.Join(dir, entry.Name())
			if err := exec.Command(path).Run(); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to run %s: %v\n", path, err)
			}
		}
	}
}

