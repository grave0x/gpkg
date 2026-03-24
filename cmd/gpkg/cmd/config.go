package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grave0x/gpkg/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `Manage gpkg configuration.

Subcommands:
  get <key>       Get a configuration value
  set <key> <val> Set a configuration value
  show            Display merged configuration

Examples:
  gpkg config get prefix
  gpkg config set log_level debug
  gpkg config show`,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		loader := config.NewYAMLLoader("")
		cfg, err := loader.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		cfg = loader.MergeDefaults(cfg)

		value := getConfigValue(cfg, key)
		if value == "" {
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		fmt.Println(value)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		val := args[1]

		loader := config.NewYAMLLoader("")
		cfg, err := loader.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		cfg = loader.MergeDefaults(cfg)

		// Update config
		switch key {
		case "prefix":
			cfg.Prefix = val
		case "cache_dir":
			cfg.CacheDir = val
		case "log_level":
			cfg.LogLevel = val
		case "network_timeout":
			fmt.Sscanf(val, "%d", &cfg.NetworkTimeout)
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		// Write to user config file
		homeDir, _ := os.UserHomeDir()
		configFile := filepath.Join(homeDir, ".gpkg", "config.yaml")

		if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		if err := os.WriteFile(configFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		if !quiet {
			fmt.Printf("✓ Set %s = %s\n", key, val)
		}
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display merged configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		loader := config.NewYAMLLoader("")
		cfg, err := loader.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		cfg = loader.MergeDefaults(cfg)

		if jsonOutput {
			return printJSON(cfg)
		}

		fmt.Println("Configuration (merged):")
		fmt.Printf("  Prefix:          %s\n", cfg.Prefix)
		fmt.Printf("  Cache Dir:       %s\n", cfg.CacheDir)
		fmt.Printf("  Sources File:    %s\n", cfg.SourcesFile)
		fmt.Printf("  Log Level:       %s\n", cfg.LogLevel)
		fmt.Printf("  Color:           %v\n", cfg.Color)
		fmt.Printf("  Strict Checksum: %v\n", cfg.StrictChecksum)
		fmt.Printf("  Network Timeout: %d seconds\n", cfg.NetworkTimeout)
		return nil
	},
}

func getConfigValue(cfg *config.Config, key string) string {
	switch key {
	case "prefix":
		return cfg.Prefix
	case "cache_dir":
		return cfg.CacheDir
	case "sources_file":
		return cfg.SourcesFile
	case "log_level":
		return cfg.LogLevel
	case "network_timeout":
		return fmt.Sprintf("%d", cfg.NetworkTimeout)
	default:
		return ""
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
}
