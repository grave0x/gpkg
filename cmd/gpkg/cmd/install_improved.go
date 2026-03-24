package cmd

import (
	"context"
	"fmt"

	"github.com/grave0x/gpkg/internal/manifest"
	"github.com/grave0x/gpkg/internal/planner"
	"github.com/spf13/cobra"
)

// UpdatedInstallCmd is an improved install command integrating planner and error handling
var updatedInstallCmd = &cobra.Command{
	Use:   "install <pkg|manifest> [--from-release|--from-source]",
	Short: "Install a package",
	Long: `Install a package from a release binary or build from source.

Examples:
  gpkg install owner/repo
  gpkg install ./manifest.yaml --from-source`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgOrManifest := args[0]

		// Determine if fromRelease or fromSource
		if !fromRelease && !fromSource {
			fromRelease = true // Default to release
		}

		if fromRelease && fromSource {
			return &GPkgError{
				Code:    ExitUsageError,
				Message: "Cannot specify both --from-release and --from-source",
			}
		}

		// Parse manifest
		mfInterface, err := parseManifestOrPackage(pkgOrManifest)
		if err != nil {
			return &GPkgError{
				Code:    ExitManifestValidation,
				Message: "Failed to resolve package",
				Detail:  err.Error(),
			}
		}

		mf, ok := mfInterface.(*manifest.Manifest)
		if !ok {
			return &GPkgError{
				Code:    ExitManifestValidation,
				Message: "Invalid manifest format",
				Detail:  "Could not parse manifest",
			}
		}

		// Create planner
		p := planner.NewDefaultPlanner(offline, dryRun)

		// Generate installation plan
		ctx := context.Background()
		plan, err := p.PlanInstallation(ctx, mf, fromRelease)
		if err != nil {
			return &GPkgError{
				Code:    ExitGeneralFailure,
				Message: "Failed to create installation plan",
				Detail:  err.Error(),
			}
		}

		// Validate plan
		if err := p.ValidatePlan(plan); err != nil {
			return &GPkgError{
				Code:    ExitManifestValidation,
				Message: "Invalid installation plan",
				Detail:  err.Error(),
			}
		}

		// Show plan for dry-run
		if dryRun {
			fmt.Println("[DRY RUN] Installation Plan:")
			fmt.Printf("Package: %s v%s\n", plan.Package, plan.Version)
			fmt.Printf("Estimated Time: %d seconds\n", plan.EstimatedTime)
			fmt.Printf("Total Size: %.2f MB\n", float64(plan.TotalSize)/(1024*1024))
			fmt.Println("\nActions:")
			for i, action := range plan.Actions {
				fmt.Printf("  %d. [%s] %s\n", i+1, action.Type, action.Description)
			}

			if len(plan.Warnings) > 0 {
				fmt.Println("\nWarnings:")
				for _, w := range plan.Warnings {
					fmt.Printf("  - %s\n", w)
				}
			}

			if jsonOutput {
				return printJSON(plan)
			}
			return nil
		}

		// Would execute plan here
		fmt.Printf("Installing %s v%s...\n", plan.Package, plan.Version)

		return nil
	},
}

func parseManifestOrPackage(identifier string) (interface{}, error) {
	// This is a placeholder - would implement actual parsing
	return nil, fmt.Errorf("not implemented")
}

func init() {
	// Would add updatedInstallCmd to rootCmd
}
