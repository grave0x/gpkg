package pkg

import "time"

// Package represents a software package that can be installed.
type Package struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Author      string            `json:"author,omitempty"`
	URL         string            `json:"url,omitempty"`
	License     string            `json:"license,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// InstalledPackage represents a package with installation metadata.
type InstalledPackage struct {
	Package      *Package
	InstalledAt  time.Time         `json:"installed_at"`
	Source       string            `json:"source"`
	Prefix       string            `json:"prefix"`
	Checksums    map[string]string `json:"checksums,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
}

// VersionInfo represents version information for a package.
type VersionInfo struct {
	Latest      string   `json:"latest"`
	Current     string   `json:"current,omitempty"`
	Available   []string `json:"available,omitempty"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
}
