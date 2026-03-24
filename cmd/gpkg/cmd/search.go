package cmd

import (
	"context"
	"fmt"

	"github.com/grave0x/gpkg/internal/source"
	"github.com/spf13/cobra"
)

type Source = source.Source

var searchSource string

var searchCmd = &cobra.Command{
	Use:   "search <term>",
	Short: "Search package indices",
	Long: `Search across configured package sources for packages matching a term.

Examples:
  gpkg search golang
  gpkg search --source example golang
  gpkg search --json curl`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		term := args[0]
		ctx := context.Background()

		sources, err := GlobalRegistry.ListSources(ctx)
		if err != nil {
			return fmt.Errorf("failed to list sources: %w", err)
		}

		if len(sources) == 0 {
			return fmt.Errorf("no sources configured; use 'gpkg add-source' to add one")
		}

		// Filter by source if specified
		var sourceList []*Source
		if searchSource != "" {
			found := false
			for _, src := range sources {
				if src.ID == searchSource || src.URI == searchSource {
					sourceList = append(sourceList, src)
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("source not found: %s", searchSource)
			}
		} else {
			sourceList = sources
		}

		type SearchResult struct {
			Package string `json:"package"`
			Version string `json:"version"`
			Source  string `json:"source"`
			Summary string `json:"summary,omitempty"`
		}

		var results []SearchResult

		// In a real implementation, this would query each source
		// For now, mock results
		if !quiet {
			fmt.Printf("Searching %d source(s) for '%s'...\n", len(sourceList), term)
		}

		// This would typically fetch from remote indices
		for _, src := range sourceList {
			if !quiet {
				fmt.Printf("  Searching %s...\n", src.URI)
			}
			// Would query source here
		}

		if jsonOutput {
			return printJSON(results)
		}

		if len(results) == 0 {
			if !quiet {
				fmt.Printf("No packages found matching '%s'\n", term)
			}
			return nil
		}

		fmt.Printf("%-20s %-15s %-30s\n", "Package", "Version", "Source")
		fmt.Println("----" + "------" + "-----" + "------")
		for _, r := range results {
			fmt.Printf("%-20s %-15s %-30s\n", r.Package, r.Version, r.Source)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&searchSource, "source", "", "search in specific source")
}
