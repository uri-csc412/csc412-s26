package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"inklings/internal/asciiart"
	"inklings/internal/config"
	"inklings/internal/sim"
)

/*
* so this works across different makefiles
 */
func loadConfigOrGoByeBye() *config.Config {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("ERROR: unable to get executable path: %v\n", err)
	}
	exeDir := filepath.Dir(exePath)

	cfgPath := filepath.Join(exeDir, "..", "..", "config", "inklings_config.yaml")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to load config from %s: %v\n", cfgPath, err)
		if cwd, cwdErr := os.Getwd(); cwdErr == nil {
			fmt.Println("Current working directory:", cwd)
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: failed to get cwd: %v\n", cwdErr)
		}
		fmt.Println("Executable directory:", exeDir)

		os.Exit(1)
	}

	return cfg
}

func main() {
	cfg := loadConfigOrGoByeBye()

	// TODO: override cfg.Simulation rows/cols/threads using os.Args if desired.

	state := sim.NewState(cfg)

	// TODO: start goroutines for inklings and ink producers.

	asciiart.EventLoop(state)
}
