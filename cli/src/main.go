package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information overriden using ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)


func hopCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hop [container]",
		Short: "Hop to the specified image",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			imageRef := args[0]

			fmt.Fprintf(os.Stderr, "Hopping to: %v\n", imageRef)
			err := hop(imageRef)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to hop: %v\n", err)
				os.Exit(1)
			} else {
				fmt.Fprintf(os.Stderr, "Hop succeeded! System will switch to %v on next boot.\n", imageRef)
			}
		},
	}

	return cmd
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Use:   "bootc-hopper",
		Short: "Hop to bootable container images",
	}

	rootCmd.AddCommand(hopCmd(config))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
