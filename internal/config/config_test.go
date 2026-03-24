package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grave0x/gpkg/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg == nil {
		t.Fatal("expected config but got nil")
	}

	if cfg.LogLevel != "info" {
		t.Errorf("expected log level 'info', got %s", cfg.LogLevel)
	}

	if cfg.Color != true {
		t.Errorf("expected color enabled, got %v", cfg.Color)
	}

	if cfg.StrictChecksum != true {
		t.Errorf("expected strict checksum enabled, got %v", cfg.StrictChecksum)
	}

	if cfg.NetworkTimeout != 30 {
		t.Errorf("expected network timeout 30, got %d", cfg.NetworkTimeout)
	}
}

func TestYAMLConfigLoading(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Create test config
	configContent := `
prefix: /custom/prefix
cache_dir: /custom/cache
log_level: debug
color: false
strict_checksum: true
network_timeout: 60
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	loader := config.NewYAMLLoader("")
	cfg, err := loader.LoadFrom(configFile)

	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Prefix != "/custom/prefix" {
		t.Errorf("expected prefix '/custom/prefix', got %s", cfg.Prefix)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("expected log level 'debug', got %s", cfg.LogLevel)
	}

	if cfg.NetworkTimeout != 60 {
		t.Errorf("expected network timeout 60, got %d", cfg.NetworkTimeout)
	}
}

func TestMergeDefaults(t *testing.T) {
	loader := config.NewYAMLLoader("")

	cfg := &config.Config{
		Prefix:   "/custom/prefix",
		LogLevel: "debug",
		// Other fields unset
	}

	merged := loader.MergeDefaults(cfg)

	if merged.Prefix != "/custom/prefix" {
		t.Errorf("expected prefix '/custom/prefix', got %s", merged.Prefix)
	}

	if merged.LogLevel != "debug" {
		t.Errorf("expected log level 'debug', got %s", merged.LogLevel)
	}

	if merged.NetworkTimeout != 30 {
		t.Errorf("expected network timeout default 30, got %d", merged.NetworkTimeout)
	}
}
