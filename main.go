package main

import (
	"fmt"
	"os"

	"github.com/krillinai/KrillinAI/internal/app"
	"github.com/krillinai/KrillinAI/internal/config"
)

// Version information injected at build time via ldflags
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func main() {
	// Print version banner
	fmt.Printf("KrillinAI %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)

	// Load application configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize and run the application
	application, err := app.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application exited with error: %v\n", err)
		os.Exit(1)
	}
}
