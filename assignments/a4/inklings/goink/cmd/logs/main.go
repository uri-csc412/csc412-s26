package main

import (
	"fmt"
	"inklings/internal/config"
	"log"
	"os"
	"path/filepath"
)

/* logs program

   Read all inkling*.txt files from the logs dir
   Merge lines into a actions.txt sorted by timestamp
   Be idempotent (running logs multiple times does not duplicate entries)
   Include all log entries (no missing lines)

*/

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

	logDir := filepath.Join(".", cfg.Logs.Directory)
	if _, err := os.Stat(logDir); err != nil {
		fmt.Fprintf(os.Stderr, "logs_go: could not stat %s (run inklings first?)\n", logDir)
		os.Exit(1)
	}

	fmt.Println("logs_go stub: please implement merging logFolder/*.txt into actions.txt")
}
