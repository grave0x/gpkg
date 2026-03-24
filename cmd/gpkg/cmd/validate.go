package cmd

import (
	"fmt"
	"os"

	"github.com/grave0x/gpkg/internal/manifest"
	"github.com/spf13/cobra"
)

var fixManifest bool

var validateCmd = &cobra.Command{
	Use:   "validate <manifest-path>",
	Short: "Validate a manifest file",
	Long: `Validate a manifest against the schema.

Checks for:
- Required fields (name, version)
- Valid install or build_source specification
- Proper checksum format
- Valid build commands

Examples:
  gpkg validate ./my-tool.yaml
  gpkg validate ./manifest.yaml --fix`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestPath := args[0]

		parser := manifest.NewYAMLParser()
		mf, err := parser.Parse(manifestPath)
		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stderr, "✗ Validation failed:\n")
				fmt.Fprintf(os.Stderr, "  %v\n", err)
			}
			return err
		}

		// Additional validation
		warnings := validateManifestStrict(mf)
		if len(warnings) > 0 {
			if !quiet {
				fmt.Fprintf(os.Stderr, "⚠ Warnings:\n")
				for _, w := range warnings {
					fmt.Fprintf(os.Stderr, "  - %s\n", w)
				}
			}
		}

		if jsonOutput {
			result := map[string]interface{}{
				"valid":    true,
				"warnings": warnings,
			}
			return printJSON(result)
		}

		if !quiet {
			fmt.Printf("✓ Manifest is valid\n")
			if len(warnings) > 0 {
				fmt.Printf("  (%d warnings)\n", len(warnings))
			}
		}

		return nil
	},
}

func validateManifestStrict(mf *manifest.Manifest) []string {
	var warnings []string

	if mf.Package.Author == "" {
		warnings = append(warnings, "package author not specified")
	}
	if mf.Package.URL == "" {
		warnings = append(warnings, "package URL not specified")
	}
	if mf.Package.License == "" {
		warnings = append(warnings, "package license not specified")
	}

	if mf.Install.Type != "" && len(mf.Install.Checksum) == 0 {
		warnings = append(warnings, "no checksums specified for release install")
	}

	if len(mf.Dependencies) > 0 {
		warnings = append(warnings, fmt.Sprintf("manifest has %d dependencies (dependency resolution not yet implemented)", len(mf.Dependencies)))
	}

	return warnings
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().BoolVar(&fixManifest, "fix", false, "attempt to fix trivial issues")
}
