package source_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/grave0x/gpkg/internal/source"
)

func TestJSONRegistry(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	registryFile := filepath.Join(tmpDir, "sources.json")

	registry := source.NewJSONRegistry(registryFile)

	ctx := context.Background()

	// Test AddSource
	src := &source.Source{
		ID:  "test-source",
		URI: "https://example.com/packages",
		Name: "Test Source",
	}

	err := registry.AddSource(ctx, src)
	if err != nil {
		t.Fatalf("failed to add source: %v", err)
	}

	// Test ListSources
	sources, err := registry.ListSources(ctx)
	if err != nil {
		t.Fatalf("failed to list sources: %v", err)
	}

	if len(sources) != 1 {
		t.Errorf("expected 1 source, got %d", len(sources))
	}

	// Test GetSource
	retrieved, err := registry.GetSource(ctx, "test-source")
	if err != nil {
		t.Fatalf("failed to get source: %v", err)
	}

	if retrieved.URI != src.URI {
		t.Errorf("expected URI %s, got %s", src.URI, retrieved.URI)
	}

	// Verify file was created
	if _, err := os.Stat(registryFile); os.IsNotExist(err) {
		t.Errorf("registry file was not created")
	}

	// Test RemoveSource
	err = registry.RemoveSource(ctx, "test-source")
	if err != nil {
		t.Fatalf("failed to remove source: %v", err)
	}

	sources, err = registry.ListSources(ctx)
	if err != nil {
		t.Fatalf("failed to list sources after removal: %v", err)
	}

	if len(sources) != 0 {
		t.Errorf("expected 0 sources after removal, got %d", len(sources))
	}
}

func TestRegistryPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	registryFile := filepath.Join(tmpDir, "sources.json")

	ctx := context.Background()

	// Create and populate registry
	{
		registry := source.NewJSONRegistry(registryFile)
		src := &source.Source{
			ID:  "persistent-source",
			URI: "https://persistent.example.com",
			Name: "Persistent Source",
		}
		err := registry.AddSource(ctx, src)
		if err != nil {
			t.Fatalf("failed to add source: %v", err)
		}
	}

	// Load registry from disk
	{
		registry := source.NewJSONRegistry(registryFile)
		err := registry.Load()
		if err != nil {
			t.Fatalf("failed to load registry: %v", err)
		}

		sources, err := registry.ListSources(ctx)
		if err != nil {
			t.Fatalf("failed to list sources: %v", err)
		}

		if len(sources) != 1 {
			t.Errorf("expected 1 source after reload, got %d", len(sources))
		}

		if sources[0].ID != "persistent-source" {
			t.Errorf("expected source ID 'persistent-source', got %s", sources[0].ID)
		}
	}
}
