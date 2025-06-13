package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pedrobarco/vibecast/internal/config"
	"github.com/pedrobarco/vibecast/internal/tui"
)

func main() {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".config", "vibecast", "config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if _, err := tui.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}
