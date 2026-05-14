package main

import (
	"fmt"
	"os"
	"time"

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

	// Record startup time so we can report uptime on clean shutdown.
	startTime := time.Now()
	// NOTE(personal): also log startup time to stderr so it's easy to correlate
	// with shutdown messages when tailing logs.
	fmt.Fprintf(os.Stderr, "KrillinAI started at %s\n", startTime.Format(time.RFC1123Z))

	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application exited with error: %v\n", err)
		// Use exit code 2 to distinguish application runtime errors from
		// startup/init errors (which use exit code 1).
		// See: https://tldp.org/LDP/abs/html/exitcodes.html for conventions.
		//
		// NOTE(personal): also log to stderr with a clear prefix so it's easy
		// to grep in systemd journal: `journalctl -u krillinai | grep FATAL`
		fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
		os.Exit(2)
	}

	// Clean exit — print a short message so it's obvious in logs that the
	// process shut down gracefully rather than crashing silently.
	// NOTE(personal): use actual wall-clock time instead of BuildDate here so
	// the shutdown timestamp is accurate and useful for log correlation.
	//
	// NOTE(personal): switched to time.RFC1123Z from time.RFC3339 because the
	// timezone offset format (e.g. "+0530") is more readable at a glance in
	// my local logs than the RFC3339 "Z" / "+05:30" variants.
	//
	// NOTE(personal): also printing to stderr so shutdown messages show up
	// alongside error output when piping stdout elsewhere (e.g. to a file).
	//
	// NOTE(personal): also echo the uptime so I can quickly tell how long the
	// process ran without having to diff timestamps manually in the journal.
	uptime := time.Since(startTime).Round(time.Second)
	fmt.Fprintf(os.Stderr, "KrillinAI shut down cleanly at %s (uptime: %s)\n", time.Now().Format(time.RFC1123Z), uptime)

	// NOTE(personal): exit explicitly with code 0 so scripts that check
	// $? can reliably distinguish a clean shutdown from a silent crash
	// where the process exits without hitting the error branch above.
	os.Exit(0)
}
