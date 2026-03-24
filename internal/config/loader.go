package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// YAMLLoader implements the Loader interface using YAML files.
type YAMLLoader struct {
	configPath string
}

// NewYAMLLoader creates a new YAML-based config loader.
func NewYAMLLoader(configPath string) *YAMLLoader {
	return &YAMLLoader{
		configPath: configPath,
	}
}

// Load reads config from default locations and environment with precedence.
// Precedence: CLI flags > environment > user config > system config
func (l *YAMLLoader) Load() (*Config, error) {
	cfg := DefaultConfig()

	// Start with system config if it exists
	systemConfigPaths := []string{
		"/etc/gpkg/config.yaml",
		"/etc/gpkg.yaml",
	}
	for _, path := range systemConfigPaths {
		if userCfg, err := l.LoadFrom(path); err == nil {
			cfg = l.mergeConfigs(cfg, userCfg)
			break
		}
	}

	// Then user config
	userHome, _ := os.UserHomeDir()
	userConfigPaths := []string{
		filepath.Join(userHome, ".gpkg", "config.yaml"),
		filepath.Join(userHome, ".gpkg.yaml"),
		filepath.Join(userHome, ".config", "gpkg", "config.yaml"),
	}
	for _, path := range userConfigPaths {
		if userCfg, err := l.LoadFrom(path); err == nil {
			cfg = l.mergeConfigs(cfg, userCfg)
			break
		}
	}

	// Environment variables override
	if prefix := os.Getenv("GPKG_PREFIX"); prefix != "" {
		cfg.Prefix = prefix
	}
	if cacheDir := os.Getenv("GPKG_CACHE_DIR"); cacheDir != "" {
		cfg.CacheDir = cacheDir
	}
	if logLevel := os.Getenv("GPKG_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	return cfg, nil
}

// LoadFrom reads config from a specific file.
func (l *YAMLLoader) LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return &cfg, nil
}

// MergeDefaults applies default values for unset fields.
func (l *YAMLLoader) MergeDefaults(cfg *Config) *Config {
	defaults := DefaultConfig()

	if cfg.Prefix == "" {
		cfg.Prefix = defaults.Prefix
	}
	if cfg.CacheDir == "" {
		cfg.CacheDir = defaults.CacheDir
	}
	if cfg.SourcesFile == "" {
		cfg.SourcesFile = defaults.SourcesFile
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = defaults.LogLevel
	}
	if cfg.NetworkTimeout == 0 {
		cfg.NetworkTimeout = defaults.NetworkTimeout
	}

	return cfg
}

// mergeConfigs merges source config into base config.
func (l *YAMLLoader) mergeConfigs(base, source *Config) *Config {
	if source.Prefix != "" {
		base.Prefix = source.Prefix
	}
	if source.CacheDir != "" {
		base.CacheDir = source.CacheDir
	}
	if source.SourcesFile != "" {
		base.SourcesFile = source.SourcesFile
	}
	if source.LogLevel != "" {
		base.LogLevel = source.LogLevel
	}
	if source.NetworkTimeout != 0 {
		base.NetworkTimeout = source.NetworkTimeout
	}
	return base
}
