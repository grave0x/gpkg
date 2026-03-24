package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/grave0x/gpkg/internal/source"
	"github.com/spf13/cobra"
)

// addSourceCmd represents the add-source subcommand
var addSourceCmd = &cobra.Command{
	Use:   "add-source <uri>",
	Short: "Add a new package source",
	Long: `Add a new package source. Sources are package indexes or repositories
where packages can be fetched from.

Example:
  gpkg add-source https://packages.example.com/index.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceURI := args[0]

		// Generate ID from URI
		sourceID := fmt.Sprintf("source-%s", sourceURI)

		src := &source.Source{
			ID:   sourceID,
			URI:  sourceURI,
			Name: sourceID,
		}

		ctx := context.Background()
		if err := GlobalRegistry.AddSource(ctx, src); err != nil {
			return fmt.Errorf("failed to add source: %w", err)
		}

		if !quiet {
			fmt.Printf("Added source: %s (%s)\n", sourceID, sourceURI)
		}
		return nil
	},
}

// removeSourceCmd represents the remove-source subcommand
var removeSourceCmd = &cobra.Command{
	Use:   "remove-source <id|uri>",
	Short: "Remove a package source",
	Long: `Remove a package source by its ID or URI.

Example:
  gpkg remove-source source-1
  gpkg remove-source https://packages.example.com/index.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idOrURI := args[0]

		ctx := context.Background()
		if err := GlobalRegistry.RemoveSource(ctx, idOrURI); err != nil {
			return fmt.Errorf("failed to remove source: %w", err)
		}

		if !quiet {
			fmt.Printf("Removed source: %s\n", idOrURI)
		}
		return nil
	},
}

// listSourcesCmd represents the list-sources subcommand
var listSourcesCmd = &cobra.Command{
	Use:   "list-sources",
	Short: "List all package sources",
	Long: `List all registered package sources.

Example:
  gpkg list-sources
  gpkg list-sources --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		sources, err := GlobalRegistry.ListSources(ctx)
		if err != nil {
			return fmt.Errorf("failed to list sources: %w", err)
		}

		if jsonOutput {
			return printJSON(sources)
		}

		if len(sources) == 0 {
			fmt.Println("No sources configured.")
			return nil
		}

		fmt.Printf("%-20s %-40s %-10s\n", "ID", "URI", "Status")
		fmt.Println("----" + "----------" + "-----" + "--------" + "-----")

		for _, src := range sources {
			status := "enabled"
			if !src.Enabled {
				status = "disabled"
			}
			fmt.Printf("%-20s %-40s %-10s\n", src.ID, src.URI, status)
		}
		return nil
	},
}

func init() {
	// Add source management commands to root
	rootCmd.AddCommand(addSourceCmd)
	rootCmd.AddCommand(removeSourceCmd)
	rootCmd.AddCommand(listSourcesCmd)

	// Initialize global registry if needed
	if GlobalRegistry == nil {
		// This will be properly initialized with config path
		GlobalRegistry = source.NewJSONRegistry("~/.gpkg/sources.json")
		if err := GlobalRegistry.Load(); err != nil {
			log.Printf("Warning: failed to load sources: %v\n", err)
		}
	}
}
