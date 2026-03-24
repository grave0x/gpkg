package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grave0x/gpkg/internal/config"
	"github.com/spf13/cobra"
)

var rollbackVersion string

var rollbackCmd = &cobra.Command{
	Use:   "rollback <pkg>",
	Short: "Rollback package to a previous version",
	Long: `Rollback an installed package to a previous version.

Requires version history to be available in the package database.

Examples:
  gpkg rollback my-tool --to-version 1.0.0
  gpkg rollback my-tool --to-version 1.0.0 --prefix /opt/gpkg`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgName := args[0]

		if rollbackVersion == "" {
			return fmt.Errorf("--to-version flag is required")
		}

		loader := config.NewYAMLLoader("")
		cfg, err := loader.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		cfg = loader.MergeDefaults(cfg)

		if prefixFlag != "" {
			cfg.Prefix = prefixFlag
		}

		// In a real implementation, this would query the package database
		backupDir := filepath.Join(cfg.Prefix, "backups")
		versionPath := filepath.Join(backupDir, pkgName, rollbackVersion, "bin", pkgName)

		if dryRun {
			fmt.Printf("[DRY RUN] Would rollback %s to v%s\n", pkgName, rollbackVersion)
			fmt.Printf("[DRY RUN] Restore from: %s\n", versionPath)
			return nil
		}

		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			return fmt.Errorf("version %s not found in backups for package %s", rollbackVersion, pkgName)
		}

		// Perform rollback
		currentBin := filepath.Join(cfg.Prefix, "bin", pkgName)
		if err := os.Rename(versionPath, currentBin); err != nil {
			return fmt.Errorf("failed to restore version: %w", err)
		}

		if !quiet {
			fmt.Printf("✓ Rolled back %s to v%s\n", pkgName, rollbackVersion)
		}

		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list [--installed|--available]",
	Short: "List packages",
	Long: `List installed or available packages.

Options:
  --installed   Show installed packages (default)
  --available   Show available packages from sources
  --filter      Filter by name pattern
  --sort        Sort by (name, version, date)

Examples:
  gpkg list
  gpkg list --available
  gpkg list --filter "tool"
  gpkg list --sort version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _ = cmd.Flags().GetBool("installed")
		availableFlag, _ := cmd.Flags().GetBool("available")
		filterPattern, _ := cmd.Flags().GetString("filter")
		sortBy, _ := cmd.Flags().GetString("sort")

		if availableFlag {
			return listAvailablePackages(filterPattern, sortBy)
		}

		// Default to installed
		return listInstalledPackages(filterPattern, sortBy)
	},
}

func listInstalledPackages(filter, sortBy string) error {
	loader := config.NewYAMLLoader("")
	cfg, err := loader.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}
	cfg = loader.MergeDefaults(cfg)

	binDir := filepath.Join(cfg.Prefix, "bin")

	entries, err := os.ReadDir(binDir)
	if err != nil {
		if os.IsNotExist(err) {
			if !quiet {
				fmt.Println("No installed packages.")
			}
			return nil
		}
		return fmt.Errorf("failed to read bin directory: %w", err)
	}

	type PkgEntry struct {
		Name        string
		Installed   string
		Version     string
	}

	var packages []PkgEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if filter == "" || contains(name, filter) {
				packages = append(packages, PkgEntry{
					Name:      name,
					Installed: "~/.gpkg/bin",
					Version:   "unknown", // Would come from pkgdb
				})
			}
		}
	}

	if jsonOutput {
		return printJSON(packages)
	}

	if len(packages) == 0 {
		if !quiet {
			fmt.Println("No installed packages.")
		}
		return nil
	}

	fmt.Printf("%-20s %-20s %-15s\n", "Package", "Location", "Version")
	fmt.Println("----" + "------" + "-----" + "------")
	for _, pkg := range packages {
		fmt.Printf("%-20s %-20s %-15s\n", pkg.Name, pkg.Installed, pkg.Version)
	}

	return nil
}

func listAvailablePackages(filter, sortBy string) error {
	ctx := context.Background()

	sources, err := GlobalRegistry.ListSources(ctx)
	if err != nil {
		return fmt.Errorf("failed to list sources: %w", err)
	}

	if len(sources) == 0 {
		return fmt.Errorf("no sources configured; use 'gpkg add-source' to add one")
	}

	if !quiet {
		fmt.Printf("Fetching packages from %d source(s)...\n", len(sources))
	}

	type AvailablePkg struct {
		Package string `json:"package"`
		Version string `json:"version"`
		Source  string `json:"source"`
	}

	var packages []AvailablePkg

	// In real implementation, would query sources
	for _, src := range sources {
		if !quiet {
			fmt.Printf("  %s: no packages fetched (mock)\n", src.URI)
		}
	}

	if jsonOutput {
		return printJSON(packages)
	}

	if len(packages) == 0 {
		if !quiet {
			fmt.Println("No packages found in sources.")
		}
	}

	return nil
}

func contains(str, substr string) bool {
	return len(substr) == 0 || len(str) >= len(substr) && str[:len(substr)] == substr || len(str) >= len(substr)
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(listCmd)

	rollbackCmd.Flags().StringVar(&rollbackVersion, "to-version", "", "version to rollback to")
	rollbackCmd.Flags().StringVar(&prefixFlag, "prefix", "", "installation prefix")

	listCmd.Flags().Bool("installed", true, "show installed packages")
	listCmd.Flags().Bool("available", false, "show available packages from sources")
	listCmd.Flags().String("filter", "", "filter by name pattern")
	listCmd.Flags().String("sort", "name", "sort by (name, version, date)")
}
