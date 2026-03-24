package source

import "context"

// Source represents a package source (e.g., package index).
type Source struct {
	ID           string `json:"id"`
	URI          string `json:"uri"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	LastUpdated  int64  `json:"last_updated,omitempty"`
	Enabled      bool   `json:"enabled"`
	SourceType   string `json:"type,omitempty"` // "http", "local", "github"
}

// Registry manages available package sources.
type Registry interface {
	// AddSource adds a new source to the registry.
	AddSource(ctx context.Context, src *Source) error

	// RemoveSource removes a source by ID or URI.
	RemoveSource(ctx context.Context, idOrURI string) error

	// ListSources returns all registered sources.
	ListSources(ctx context.Context) ([]*Source, error)

	// GetSource retrieves a source by ID.
	GetSource(ctx context.Context, id string) (*Source, error)

	// UpdateSource updates a source's metadata.
	UpdateSource(ctx context.Context, src *Source) error
}

// Fetcher handles retrieving package information from a source.
type Fetcher interface {
	// Fetch retrieves package metadata from the source.
	Fetch(ctx context.Context, src *Source) (interface{}, error)

	// FetchPackage retrieves specific package info.
	FetchPackage(ctx context.Context, src *Source, pkgName string) (interface{}, error)
}
