package pkgdb

import "time"

// PackageRecord represents an installed package in the database.
type PackageRecord struct {
	ID            int64
	Name          string
	Version       string
	InstalledAt   time.Time
	UpdatedAt     time.Time
	Source        string // "release" or "source"
	Prefix        string
	Author        string
	URL           string
	License       string
	Checksums     map[string]string
	Dependencies  []string
	Files         []string // List of installed files
	BuildMetadata map[string]string
}

// VersionRecord represents a historical version of a package.
type VersionRecord struct {
	ID          int64
	PackageID   int64
	Version     string
	InstalledAt time.Time
	Checksums   map[string]string
	Files       []string
}

// FileRecord represents a file installed by a package.
type FileRecord struct {
	ID        int64
	PackageID int64
	FilePath  string
	Checksum  string
	FileSize  int64
}

// Manager provides database operations for package records.
type Manager interface {
	// AddPackage records a newly installed package.
	AddPackage(p *PackageRecord) (int64, error)

	// GetPackage retrieves a package by name.
	GetPackage(name string) (*PackageRecord, error)

	// UpdatePackage updates package metadata.
	UpdatePackage(p *PackageRecord) error

	// DeletePackage removes a package record.
	DeletePackage(name string) error

	// ListPackages returns all installed packages.
	ListPackages() ([]*PackageRecord, error)

	// AddFiles records files installed by a package.
	AddFiles(packageID int64, files []string) error

	// GetFiles retrieves files for a package.
	GetFiles(packageID int64) ([]string, error)

	// AddVersion records a new version in history.
	AddVersion(packageID int64, v *VersionRecord) (int64, error)

	// GetVersionHistory retrieves version history for a package.
	GetVersionHistory(packageID int64) ([]*VersionRecord, error)

	// Close closes the database connection.
	Close() error
}
