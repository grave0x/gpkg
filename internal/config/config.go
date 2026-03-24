package config

import "os"

// Config holds application configuration with proper precedence.
type Config struct {
	// Installation prefix (default: ~/.gpkg)
	Prefix string `yaml:"prefix" json:"prefix"`

	// Cache directory
	CacheDir string `yaml:"cache_dir" json:"cache_dir"`

	// Sources registry file
	SourcesFile string `yaml:"sources_file" json:"sources_file"`

	// Logging level
	LogLevel string `yaml:"log_level" json:"log_level"`

	// Enable color output
	Color bool `yaml:"color" json:"color"`

	// Checksum verification strict mode
	StrictChecksum bool `yaml:"strict_checksum" json:"strict_checksum"`

	// Network timeout (seconds)
	NetworkTimeout int `yaml:"network_timeout" json:"network_timeout"`
}

// Loader interface for config loading with precedence.
type Loader interface {
	// Load reads config from default locations and environment.
	// Precedence: CLI flags > environment > user config > system config
	Load() (*Config, error)

	// LoadFrom reads config from a specific file.
	LoadFrom(path string) (*Config, error)

	// MergeDefaults applies default values for unset fields.
	MergeDefaults(cfg *Config) *Config
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	prefix := homeDir + "/.gpkg"

	return &Config{
		Prefix:         prefix,
		CacheDir:       prefix + "/cache",
		SourcesFile:    prefix + "/sources.json",
		LogLevel:       "info",
		Color:          true,
		StrictChecksum: true,
		NetworkTimeout: 30,
	}
}
