package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/grave0x/gpkg/internal/config"
	"github.com/grave0x/gpkg/internal/download"
	pkgmod "github.com/grave0x/gpkg/internal/package"
	"github.com/grave0x/gpkg/internal/manifest"
	"github.com/spf13/cobra"
)

var (
	fromRelease bool
	fromSource  bool
	prefixFlag  string
)

var installCmd = &cobra.Command{
	Use:   "install <pkg|manifest> [--from-release|--from-source]",
	Short: "Install a package",
	Long: `Install a package from a release binary or build from source.

Can install from:
- Package name (requires source configuration)
- Local manifest file path
- GitHub repo (with proper manifest)

Examples:
  gpkg install owner/repo --from-release
  gpkg install ./examples/manifest.yaml --from-source
  gpkg install my-package --prefix=/custom/path`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgOrManifest := args[0]

		// Load config
		loader := config.NewYAMLLoader("")
		cfg, err := loader.Load()
		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
			}
			cfg = config.DefaultConfig()
		}
		cfg = loader.MergeDefaults(cfg)

		if prefixFlag != "" {
			cfg.Prefix = prefixFlag
		}

		// Parse manifest
		var mf *manifest.Manifest
		if filepath.IsAbs(pkgOrManifest) || filepath.Base(pkgOrManifest) != pkgOrManifest {
			// Looks like a file path
			parser := manifest.NewYAMLParser()
			mf, err = parser.Parse(pkgOrManifest)
			if err != nil {
				return fmt.Errorf("failed to parse manifest: %w", err)
			}
		} else {
			// Package name - would need to fetch from sources
			return fmt.Errorf("package name resolution not implemented yet")
		}

		// Determine installation mode
		if !fromRelease && !fromSource {
			// Auto-detect: prefer release over source
		if mf.Install.Type != "" {
				fromRelease = true
			} else if mf.BuildSource != nil {
				fromSource = true
			} else {
				return fmt.Errorf("no installation method available in manifest")
			}
		}

		if fromRelease && fromSource {
			return fmt.Errorf("cannot specify both --from-release and --from-source")
		}

		// Create installer with timeout
		timeout := time.Duration(cfg.NetworkTimeout) * time.Second
		downloader := download.NewHTTPDownloader(timeout, offline)
		atomicInstaller := download.NewAtomicInstaller(filepath.Join(cfg.Prefix, "backups"))
		installer := pkgmod.NewInstaller(downloader, atomicInstaller, cfg.Prefix)

		ctx := context.Background()
		if dryRun {
			fmt.Printf("[DRY RUN] Would install %s v%s\n", mf.Package.Name, mf.Package.Version)
			if fromRelease {
				fmt.Printf("[DRY RUN] From release: %s\n", mf.Install.Source)
			} else {
				fmt.Printf("[DRY RUN] From source: %s\n", mf.BuildSource.Source)
			}
			return nil
		}

		var installed *pkgmod.InstalledPackage
		if fromRelease {
			installed, err = installer.InstallFromRelease(ctx, mf)
		} else {
			installed, err = installer.InstallFromSource(ctx, mf)
		}

		if err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		if jsonOutput {
			return printJSON(installed)
		}

		if !quiet {
			fmt.Printf("✓ Installed %s v%s to %s\n", installed.Package.Name, installed.Package.Version, cfg.Prefix)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&fromRelease, "from-release", false, "install from release binary")
	installCmd.Flags().BoolVar(&fromSource, "from-source", false, "install by building from source")
	installCmd.Flags().StringVar(&prefixFlag, "prefix", "", "installation prefix (overrides config)")
}
