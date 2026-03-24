package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grave0x/gpkg/internal/source"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"

	// Global flags
	configPath string
	jsonOutput bool
	assumeYes  bool
	dryRun     bool
	logLevel   string
	verbose    int
	quiet      bool
	noColor    bool
	offline    bool
)

var rootCmd = &cobra.Command{
	Use:   "gpkg",
	Short: "gpkg - a simple package manager for GitHub releases",
	Long: `gpkg is a user-focused package manager that installs either release binaries
or builds from source into a configurable prefix. It aims to be simple, secure,
and script-friendly.`,
	Version: version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file (overrides default locations)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "machine-readable JSON output where supported")
	rootCmd.PersistentFlags().BoolVarP(&assumeYes, "yes", "y", false, "assume yes for all prompts (non-interactive)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "plan actions without modifying filesystem/pkgdb")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set logging level (error,warn,info,debug)")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "increase verbosity (stackable)")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "minimal output (conflicts with --verbose)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colorized output")
	rootCmd.PersistentFlags().BoolVar(&offline, "offline", false, "disallow network fetches")

	// Initialize global registry
	homeDir, _ := os.UserHomeDir()
	sourcesFile := filepath.Join(homeDir, ".gpkg", "sources.json")
	GlobalRegistry = source.NewJSONRegistry(sourcesFile)
	if err := GlobalRegistry.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load sources registry: %v\n", err)
	}
}

func printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}
