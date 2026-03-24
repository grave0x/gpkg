package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion",
	Long: `Generate completion script for bash, zsh, or fish shell.

Examples:
  # bash
  gpkg completion bash | sudo tee /usr/share/bash-completion/completions/gpkg

  # zsh
  gpkg completion zsh | sudo tee /usr/share/zsh/site-functions/_gpkg

  # fish
  gpkg completion fish | sudo tee /usr/share/fish/vendor_completions.d/gpkg.fish`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]

		switch shell {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		default:
			return fmt.Errorf("unsupported shell: %s", shell)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
