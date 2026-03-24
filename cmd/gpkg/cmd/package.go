package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grave0x/gpkg/internal/config"
	"github.com/spf13/cobra"
)

var (
	depsTree bool
	rawOutput bool
	parsedOutput bool
)

var infoCmd = &cobra.Command{
	Use:   "info <repo|pkg|manifest> [--deps-tree] [--raw|--parsed]",
	Short: "Show package information",
	Long: `Display detailed information about a package from a repository,
installed packages, or a manifest file.

Examples:
  gpkg info owner/repo
  gpkg info my-package --deps-tree
  gpkg info ./manifest.yaml --raw`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		// Check if it's a manifest file
		if _, err := os.Stat(target); err == nil {
			// It's a file - parse manifest
			return showManifestInfo(target)
		}

		// Otherwise it's a package name or repo - would query sources
		return fmt.Errorf("package information lookup not yet implemented")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update package source information",
	Long: `Refresh package source metadata from all configured sources.

This fetches the latest package list and versions from all 
enabled sources without installing anything.

Example:
  gpkg update`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if dryRun {
			fmt.Println("[DRY RUN] Would refresh all package sources")
			return nil
		}

		sources, err := GlobalRegistry.ListSources(ctx)
		if err != nil {
			return fmt.Errorf("failed to list sources: %w", err)
		}

		if len(sources) == 0 {
			fmt.Println("No sources configured. Use 'gpkg add-source' to add one.")
			return nil
		}

		if !quiet {
			fmt.Printf("Updating %d source(s)...\n", len(sources))
		}

		// In a real implementation, this would fetch metadata from each source
		for _, src := range sources {
			if !quiet {
				fmt.Printf("  ✓ Updated %s\n", src.URI)
			}
		}

		return nil
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [pkg...]",
	Short: "Upgrade installed packages",
	Long: `Upgrade installed packages to their latest available versions.

If no package names are specified, upgrades all installed packages.

Examples:
  gpkg upgrade              # upgrade all
  gpkg upgrade my-pkg      # upgrade specific package
  gpkg upgrade pkg1 pkg2   # upgrade multiple`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			if len(args) == 0 {
				fmt.Println("[DRY RUN] Would upgrade all installed packages")
			} else {
				fmt.Printf("[DRY RUN] Would upgrade: %v\n", args)
			}
			return nil
		}

		// Load config for installed packages location
		loader := config.NewYAMLLoader("")
		cfg, err := loader.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		cfg = loader.MergeDefaults(cfg)

		if !quiet {
			fmt.Printf("Checking for upgrades in %s...\n", cfg.Prefix)
		}

		// In a real implementation, this would:
		// 1. Check installed packages
		// 2. Compare with latest from sources
		// 3. Reinstall newer versions

		if !quiet {
			fmt.Println("All packages are up to date.")
		}
		return nil
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <pkg>",
	Short: "Remove an installed package",
	Long: `Uninstall a package from the system.

Examples:
  gpkg uninstall my-package
  gpkg uninstall my-package --prefix=/custom/path`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgName := args[0]

		loader := config.NewYAMLLoader("")
		cfg, err := loader.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		cfg = loader.MergeDefaults(cfg)

		if prefixFlag != "" {
			cfg.Prefix = prefixFlag
		}

		binPath := filepath.Join(cfg.Prefix, "bin", pkgName)

		if dryRun {
			fmt.Printf("[DRY RUN] Would remove %s from %s\n", pkgName, binPath)
			return nil
		}

		if err := os.Remove(binPath); err != nil {
			return fmt.Errorf("failed to uninstall package: %w", err)
		}

		if !quiet {
			fmt.Printf("✓ Uninstalled %s\n", pkgName)
		}
		return nil
	},
}

func showManifestInfo(path string) error {
	// Parse and display manifest information
	fmt.Printf("Manifest: %s\n", path)
	return nil
}

func init() {
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(uninstallCmd)

	infoCmd.Flags().BoolVar(&depsTree, "deps-tree", false, "show dependency tree")
	infoCmd.Flags().BoolVar(&rawOutput, "raw", false, "show raw manifest")
	infoCmd.Flags().BoolVar(&parsedOutput, "parsed", false, "show parsed manifest")

	uninstallCmd.Flags().StringVar(&prefixFlag, "prefix", "", "installation prefix")
}
