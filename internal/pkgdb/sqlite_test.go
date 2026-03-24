package pkgdb_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/grave0x/gpkg/internal/pkgdb"
)

func TestSQLiteManager(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	manager, err := pkgdb.NewSQLiteManager(dbPath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	defer manager.Close()

	// Test AddPackage
	pkg := &pkgdb.PackageRecord{
		Name:         "test-pkg",
		Version:      "1.0.0",
		Source:       "release",
		Prefix:       "/home/user/.gpkg",
		Author:       "Test Author",
		URL:          "https://example.com",
		License:      "MIT",
		Checksums:    map[string]string{"sha256": "abc123"},
		Dependencies: []string{"dep1"},
	}

	id, err := manager.AddPackage(pkg)
	if err != nil {
		t.Fatalf("failed to add package: %v", err)
	}

	if id == 0 {
		t.Errorf("expected non-zero package ID")
	}

	// Test GetPackage
	retrieved, err := manager.GetPackage("test-pkg")
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if retrieved.Name != "test-pkg" {
		t.Errorf("expected name 'test-pkg', got %s", retrieved.Name)
	}

	if retrieved.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", retrieved.Version)
	}

	// Test AddFiles
	files := []string{
		"/home/user/.gpkg/bin/test-pkg",
		"/home/user/.gpkg/share/doc/test-pkg/README.md",
	}

	err = manager.AddFiles(id, files)
	if err != nil {
		t.Fatalf("failed to add files: %v", err)
	}

	// Test GetFiles
	retrievedFiles, err := manager.GetFiles(id)
	if err != nil {
		t.Fatalf("failed to get files: %v", err)
	}

	if len(retrievedFiles) != 2 {
		t.Errorf("expected 2 files, got %d", len(retrievedFiles))
	}

	// Test UpdatePackage
	pkg.Version = "1.1.0"
	err = manager.UpdatePackage(pkg)
	if err != nil {
		t.Fatalf("failed to update package: %v", err)
	}

	updated, err := manager.GetPackage("test-pkg")
	if err != nil {
		t.Fatalf("failed to get updated package: %v", err)
	}

	if updated.Version != "1.1.0" {
		t.Errorf("expected version '1.1.0', got %s", updated.Version)
	}

	// Test ListPackages
	packages, err := manager.ListPackages()
	if err != nil {
		t.Fatalf("failed to list packages: %v", err)
	}

	if len(packages) == 0 {
		t.Errorf("expected at least 1 package")
	}

	// Test AddVersion
	versionRecord := &pkgdb.VersionRecord{
		Version:   "1.0.0",
		Checksums: map[string]string{"sha256": "old123"},
		Files:     files,
	}

	versionID, err := manager.AddVersion(id, versionRecord)
	if err != nil {
		t.Fatalf("failed to add version: %v", err)
	}

	if versionID == 0 {
		t.Errorf("expected non-zero version ID")
	}

	// Test GetVersionHistory
	history, err := manager.GetVersionHistory(id)
	if err != nil {
		t.Fatalf("failed to get version history: %v", err)
	}

	if len(history) == 0 {
		t.Errorf("expected version history")
	}

	// Test DeletePackage
	err = manager.DeletePackage("test-pkg")
	if err != nil {
		t.Fatalf("failed to delete package: %v", err)
	}

	_, err = manager.GetPackage("test-pkg")
	if err == nil {
		t.Errorf("expected error when getting deleted package")
	}
}

func TestVersionHistory(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	manager, err := pkgdb.NewSQLiteManager(dbPath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	defer manager.Close()

	pkg := &pkgdb.PackageRecord{
		Name:    "versioned-pkg",
		Version: "2.0.0",
		Source:  "release",
		Prefix:  "/home/user/.gpkg",
	}

	id, _ := manager.AddPackage(pkg)

	// Add multiple versions
	for i := 0; i < 3; i++ {
		version := &pkgdb.VersionRecord{
			Version: "1.0." + string(rune('0'+i)),
		}
		manager.AddVersion(id, version)
		time.Sleep(10 * time.Millisecond)
	}

	history, err := manager.GetVersionHistory(id)
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}

	if len(history) != 3 {
		t.Errorf("expected 3 versions, got %d", len(history))
	}
}
