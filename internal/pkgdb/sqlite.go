package pkgdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteManager implements Manager using SQLite.
type SQLiteManager struct {
	db *sql.DB
}

// NewSQLiteManager creates a new SQLite-based package database.
func NewSQLiteManager(dbPath string) (*SQLiteManager, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	manager := &SQLiteManager{db: db}

	// Initialize schema
	if err := manager.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return manager, nil
}

// initSchema creates tables if they don't exist.
func (m *SQLiteManager) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS packages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		version TEXT NOT NULL,
		installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		source TEXT,
		prefix TEXT,
		author TEXT,
		url TEXT,
		license TEXT,
		checksums TEXT,
		dependencies TEXT,
		build_metadata TEXT
	);

	CREATE TABLE IF NOT EXISTS versions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		package_id INTEGER NOT NULL,
		version TEXT NOT NULL,
		installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		checksums TEXT,
		files TEXT,
		FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		package_id INTEGER NOT NULL,
		file_path TEXT NOT NULL,
		checksum TEXT,
		file_size INTEGER,
		FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE,
		UNIQUE(package_id, file_path)
	);

	CREATE INDEX IF NOT EXISTS idx_packages_name ON packages(name);
	CREATE INDEX IF NOT EXISTS idx_versions_package_id ON versions(package_id);
	CREATE INDEX IF NOT EXISTS idx_files_package_id ON files(package_id);
	`

	_, err := m.db.Exec(schema)
	return err
}

// AddPackage records a newly installed package.
func (m *SQLiteManager) AddPackage(p *PackageRecord) (int64, error) {
	checksums, _ := json.Marshal(p.Checksums)
	deps, _ := json.Marshal(p.Dependencies)
	metadata, _ := json.Marshal(p.BuildMetadata)

	result, err := m.db.Exec(`
		INSERT INTO packages (name, version, source, prefix, author, url, license, checksums, dependencies, build_metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, p.Name, p.Version, p.Source, p.Prefix, p.Author, p.URL, p.License, checksums, deps, metadata)

	if err != nil {
		return 0, fmt.Errorf("failed to add package: %w", err)
	}

	return result.LastInsertId()
}

// GetPackage retrieves a package by name.
func (m *SQLiteManager) GetPackage(name string) (*PackageRecord, error) {
	var p PackageRecord
	var checksums, deps, metadata []byte

	err := m.db.QueryRow(`
		SELECT id, name, version, installed_at, updated_at, source, prefix, author, url, license, checksums, dependencies, build_metadata
		FROM packages WHERE name = ?
	`, name).Scan(
		&p.ID, &p.Name, &p.Version, &p.InstalledAt, &p.UpdatedAt, &p.Source, &p.Prefix,
		&p.Author, &p.URL, &p.License, &checksums, &deps, &metadata,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("package not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get package: %w", err)
	}

	json.Unmarshal(checksums, &p.Checksums)
	json.Unmarshal(deps, &p.Dependencies)
	json.Unmarshal(metadata, &p.BuildMetadata)

	// Load files
	files, _ := m.GetFiles(p.ID)
	p.Files = files

	return &p, nil
}

// UpdatePackage updates package metadata.
func (m *SQLiteManager) UpdatePackage(p *PackageRecord) error {
	checksums, _ := json.Marshal(p.Checksums)
	deps, _ := json.Marshal(p.Dependencies)
	metadata, _ := json.Marshal(p.BuildMetadata)

	_, err := m.db.Exec(`
		UPDATE packages
		SET version = ?, source = ?, author = ?, url = ?, license = ?, checksums = ?, dependencies = ?, build_metadata = ?, updated_at = CURRENT_TIMESTAMP
		WHERE name = ?
	`, p.Version, p.Source, p.Author, p.URL, p.License, checksums, deps, metadata, p.Name)

	return err
}

// DeletePackage removes a package record.
func (m *SQLiteManager) DeletePackage(name string) error {
	_, err := m.db.Exec("DELETE FROM packages WHERE name = ?", name)
	return err
}

// ListPackages returns all installed packages.
func (m *SQLiteManager) ListPackages() ([]*PackageRecord, error) {
	rows, err := m.db.Query(`
		SELECT id, name, version, installed_at, updated_at, source, prefix, author, url, license, checksums, dependencies, build_metadata
		FROM packages ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query packages: %w", err)
	}
	defer rows.Close()

	var packages []*PackageRecord

	for rows.Next() {
		var p PackageRecord
		var checksums, deps, metadata []byte

		err := rows.Scan(
			&p.ID, &p.Name, &p.Version, &p.InstalledAt, &p.UpdatedAt, &p.Source, &p.Prefix,
			&p.Author, &p.URL, &p.License, &checksums, &deps, &metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan package: %w", err)
		}

		json.Unmarshal(checksums, &p.Checksums)
		json.Unmarshal(deps, &p.Dependencies)
		json.Unmarshal(metadata, &p.BuildMetadata)

		files, _ := m.GetFiles(p.ID)
		p.Files = files

		packages = append(packages, &p)
	}

	return packages, nil
}

// AddFiles records files installed by a package.
func (m *SQLiteManager) AddFiles(packageID int64, files []string) error {
	stmt, err := m.db.Prepare(`
		INSERT OR REPLACE INTO files (package_id, file_path)
		VALUES (?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, file := range files {
		_, err := stmt.Exec(packageID, file)
		if err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}
	}

	return nil
}

// GetFiles retrieves files for a package.
func (m *SQLiteManager) GetFiles(packageID int64) ([]string, error) {
	rows, err := m.db.Query(`
		SELECT file_path FROM files WHERE package_id = ? ORDER BY file_path
	`, packageID)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []string

	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

// AddVersion records a new version in history.
func (m *SQLiteManager) AddVersion(packageID int64, v *VersionRecord) (int64, error) {
	checksums, _ := json.Marshal(v.Checksums)
	files, _ := json.Marshal(v.Files)

	result, err := m.db.Exec(`
		INSERT INTO versions (package_id, version, checksums, files)
		VALUES (?, ?, ?, ?)
	`, packageID, v.Version, checksums, files)

	if err != nil {
		return 0, fmt.Errorf("failed to add version: %w", err)
	}

	return result.LastInsertId()
}

// GetVersionHistory retrieves version history for a package.
func (m *SQLiteManager) GetVersionHistory(packageID int64) ([]*VersionRecord, error) {
	rows, err := m.db.Query(`
		SELECT id, package_id, version, installed_at, checksums, files
		FROM versions WHERE package_id = ? ORDER BY installed_at DESC
	`, packageID)
	if err != nil {
		return nil, fmt.Errorf("failed to query versions: %w", err)
	}
	defer rows.Close()

	var versions []*VersionRecord

	for rows.Next() {
		var v VersionRecord
		var checksums, files []byte

		err := rows.Scan(&v.ID, &v.PackageID, &v.Version, &v.InstalledAt, &checksums, &files)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}

		json.Unmarshal(checksums, &v.Checksums)
		json.Unmarshal(files, &v.Files)

		versions = append(versions, &v)
	}

	return versions, nil
}

// Close closes the database connection.
func (m *SQLiteManager) Close() error {
	return m.db.Close()
}
