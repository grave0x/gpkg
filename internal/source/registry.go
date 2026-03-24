package source

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// JSONRegistry implements Registry using a JSON file storage.
type JSONRegistry struct {
	filePath string
	sources  map[string]*Source
}

// NewJSONRegistry creates a new JSON-based source registry.
func NewJSONRegistry(filePath string) *JSONRegistry {
	return &JSONRegistry{
		filePath: filePath,
		sources:  make(map[string]*Source),
	}
}

// Load reads the registry from disk.
func (r *JSONRegistry) Load() error {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return nil // File doesn't exist yet
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read registry file: %w", err)
	}

	var sources []*Source
	if err := json.Unmarshal(data, &sources); err != nil {
		return fmt.Errorf("failed to parse registry JSON: %w", err)
	}

	for _, src := range sources {
		r.sources[src.ID] = src
	}
	return nil
}

// Save writes the registry to disk.
func (r *JSONRegistry) Save() error {
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create registry directory: %w", err)
	}

	sources := make([]*Source, 0, len(r.sources))
	for _, src := range r.sources {
		sources = append(sources, src)
	}

	data, err := json.MarshalIndent(sources, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry file: %w", err)
	}
	return nil
}

// AddSource adds a new source to the registry.
func (r *JSONRegistry) AddSource(ctx context.Context, src *Source) error {
	if src.ID == "" {
		return fmt.Errorf("source ID cannot be empty")
	}
	if _, exists := r.sources[src.ID]; exists {
		return fmt.Errorf("source with ID %q already exists", src.ID)
	}

	src.Enabled = true
	src.LastUpdated = time.Now().Unix()
	r.sources[src.ID] = src
	return r.Save()
}

// RemoveSource removes a source by ID or URI.
func (r *JSONRegistry) RemoveSource(ctx context.Context, idOrURI string) error {
	var toRemove string

	if _, exists := r.sources[idOrURI]; exists {
		toRemove = idOrURI
	} else {
		// Search by URI
		for id, src := range r.sources {
			if src.URI == idOrURI {
				toRemove = id
				break
			}
		}
	}

	if toRemove == "" {
		return fmt.Errorf("source not found: %q", idOrURI)
	}

	delete(r.sources, toRemove)
	return r.Save()
}

// ListSources returns all registered sources.
func (r *JSONRegistry) ListSources(ctx context.Context) ([]*Source, error) {
	sources := make([]*Source, 0, len(r.sources))
	for _, src := range r.sources {
		sources = append(sources, src)
	}
	return sources, nil
}

// GetSource retrieves a source by ID.
func (r *JSONRegistry) GetSource(ctx context.Context, id string) (*Source, error) {
	src, exists := r.sources[id]
	if !exists {
		return nil, fmt.Errorf("source not found: %q", id)
	}
	return src, nil
}

// UpdateSource updates a source's metadata.
func (r *JSONRegistry) UpdateSource(ctx context.Context, src *Source) error {
	if _, exists := r.sources[src.ID]; !exists {
		return fmt.Errorf("source not found: %q", src.ID)
	}
	src.LastUpdated = time.Now().Unix()
	r.sources[src.ID] = src
	return r.Save()
}
