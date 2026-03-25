# API Reference

This document provides comprehensive API documentation for gpkg's internal packages. It is intended for contributors and advanced users who want to understand or extend gpkg's functionality.

## Table of Contents

- [Overview](#overview)
- [Core Packages](#core-packages)
  - [config](#package-config)
  - [download](#package-download)
  - [manifest](#package-manifest)
  - [package](#package-package)
  - [pkgdb](#package-pkgdb)
  - [planner](#package-planner)
  - [resolver](#package-resolver)
  - [source](#package-source)
- [Data Models](#data-models)
- [Database Schema](#database-schema)
- [Error Handling](#error-handling)
- [Context Usage](#context-usage)
- [Environment Variables](#environment-variables)

---

## Overview

gpkg's internal API is organized into focused packages, each handling a specific aspect of package management:

```
internal/
├── config/     # Configuration loading and management
├── download/   # HTTP downloads, checksums, atomic operations
├── manifest/   # Package manifest parsing and validation
├── package/    # Package installation logic
├── pkgdb/      # SQLite database operations
├── planner/    # Installation planning and validation
├── resolver/   # Package and dependency resolution
└── source/     # Package source registry management
```

All async operations accept `context.Context` for cancellation and timeout support. Errors are wrapped with context using `fmt.Errorf("message: %w", err)` for traceability.

---

## Core Packages

### Package: config

**Import:** `github.com/grave0x/gpkg/internal/config`

**Purpose:** Manages application configuration with precedence-based loading (CLI flags > environment > user config > system config).

#### Types

##### Config

Main configuration structure with application settings.

```go
type Config struct {
    Prefix         string `yaml:"prefix" json:"prefix"`
    CacheDir       string `yaml:"cache_dir" json:"cache_dir"`
    SourcesFile    string `yaml:"sources_file" json:"sources_file"`
    LogLevel       string `yaml:"log_level" json:"log_level"`
    Color          bool   `yaml:"color" json:"color"`
    StrictChecksum bool   `yaml:"strict_checksum" json:"strict_checksum"`
    NetworkTimeout int    `yaml:"network_timeout" json:"network_timeout"`
}
```

**Fields:**
- `Prefix` - Installation directory (default: `~/.gpkg`)
- `CacheDir` - Download cache location
- `SourcesFile` - Path to sources registry JSON file
- `LogLevel` - Logging verbosity (`debug`, `info`, `warn`, `error`)
- `Color` - Enable/disable colored terminal output
- `StrictChecksum` - Require checksum validation for all downloads
- `NetworkTimeout` - HTTP timeout in seconds

##### Loader

Interface for configuration loading with precedence.

```go
type Loader interface {
    Load() (*Config, error)
    LoadFrom(path string) (*Config, error)
    MergeDefaults(cfg *Config) *Config
}
```

**Methods:**
- `Load()` - Loads config with precedence: CLI > env > user > system
- `LoadFrom(path)` - Loads from specific YAML file
- `MergeDefaults(cfg)` - Applies default values for unset fields

#### Implementations

##### YAMLLoader

YAML-based configuration loader.

```go
type YAMLLoader struct { /* private fields */ }

func NewYAMLLoader(configPath string) *YAMLLoader
```

**Configuration Precedence:**
1. CLI flags (handled by caller)
2. Environment variables (`GPKG_PREFIX`, `GPKG_CACHE_DIR`, `GPKG_LOG_LEVEL`)
3. User config files (in order):
   - `~/.gpkg/config.yaml`
   - `~/.gpkg.yaml`
   - `~/.config/gpkg/config.yaml`
4. System config files:
   - `/etc/gpkg/config.yaml`
   - `/etc/gpkg.yaml`
5. Built-in defaults

#### Functions

```go
func DefaultConfig() *Config
```

Returns a Config with sensible defaults:
- Prefix: `~/.gpkg`
- CacheDir: `~/.gpkg/cache`
- SourcesFile: `~/.gpkg/sources.json`
- LogLevel: `info`
- Color: `true`
- StrictChecksum: `true`
- NetworkTimeout: `30` seconds

#### Usage Example

```go
import "github.com/grave0x/gpkg/internal/config"

// Load configuration with precedence
loader := config.NewYAMLLoader("")
cfg, err := loader.Load()
if err != nil {
    log.Fatal(err)
}

// Load from specific file
cfg, err := loader.LoadFrom("/etc/gpkg/config.yaml")

// Apply defaults to partial config
cfg = loader.MergeDefaults(cfg)

fmt.Printf("Installation prefix: %s\n", cfg.Prefix)
```

#### Error Handling

- Returns wrapped errors with context: `"failed to read config file: %w"`
- Returns nil error if optional config files don't exist
- Validates YAML syntax and structure

---

### Package: download

**Import:** `github.com/grave0x/gpkg/internal/download`

**Purpose:** Handles secure HTTP downloads with checksum validation and atomic filesystem operations.

#### Types

##### Downloader

Interface for downloading files with optional checksum validation.

```go
type Downloader interface {
    Download(ctx context.Context, url, dest string) error
    DownloadWithChecksum(ctx context.Context, url, dest, expectedChecksum, algorithm string) error
    ValidateChecksum(filePath, expectedChecksum, algorithm string) (bool, error)
}
```

**Methods:**
- `Download(ctx, url, dest)` - Downloads file to destination
- `DownloadWithChecksum(ctx, url, dest, checksum, algo)` - Downloads and validates checksum
- `ValidateChecksum(path, checksum, algo)` - Verifies file integrity

##### AtomicInstaller

Interface for atomic filesystem operations with rollback support.

```go
type AtomicInstaller interface {
    Install(ctx context.Context, sourceFile, destPath string) error
    Uninstall(ctx context.Context, pkgName string) error
    Rollback(ctx context.Context, installID string) error
}
```

**Methods:**
- `Install(ctx, src, dest)` - Atomically moves file from temp to final location
- `Uninstall(ctx, pkg)` - Safely removes installed package
- `Rollback(ctx, id)` - Reverts failed installation using backup

##### ChecksumAlgorithm

Supported checksum types.

```go
type ChecksumAlgorithm string

const (
    SHA256 ChecksumAlgorithm = "sha256"  // Recommended
    SHA1   ChecksumAlgorithm = "sha1"
    MD5    ChecksumAlgorithm = "md5"
)
```

#### Implementations

##### HTTPDownloader

HTTP-based downloader with checksum validation.

```go
type HTTPDownloader struct { /* private fields */ }

func NewHTTPDownloader(timeout time.Duration, allowOffline bool) *HTTPDownloader
```

**Parameters:**
- `timeout` - HTTP request timeout (e.g., `30 * time.Second`)
- `allowOffline` - If true, fails fast when network unavailable

**Features:**
- Atomic downloads (temp file → final destination)
- Context-aware cancellation
- Automatic checksum validation
- Offline mode support

##### AtomicInstallerImpl

Filesystem-based atomic installer with backup support.

```go
type AtomicInstallerImpl struct { /* private fields */ }

func NewAtomicInstaller(backupDir string) *AtomicInstallerImpl
```

**Parameters:**
- `backupDir` - Directory for backup files (e.g., `~/.gpkg/backups`)

**Features:**
- Automatic backups before replacing files
- Rollback support on failure
- Proper file permissions (0755 for executables)

#### Usage Example

```go
import (
    "context"
    "time"
    "github.com/grave0x/gpkg/internal/download"
)

// Create downloader
dl := download.NewHTTPDownloader(30*time.Second, false)

// Simple download
ctx := context.Background()
err := dl.Download(ctx, "https://example.com/file.tar.gz", "/tmp/file.tar.gz")

// Download with checksum validation
err = dl.DownloadWithChecksum(
    ctx,
    "https://example.com/binary",
    "/tmp/binary",
    "abc123...",
    "sha256",
)

// Atomic installation
installer := download.NewAtomicInstaller("/home/user/.gpkg/backups")
err = installer.Install(ctx, "/tmp/binary", "/home/user/.gpkg/bin/tool")

// Verify checksum of existing file
valid, err := dl.ValidateChecksum("/tmp/file", "abc123...", "sha256")
if !valid {
    log.Fatal("Checksum mismatch!")
}
```

#### Error Handling

- Network errors include original error wrapped with context
- Checksum mismatches automatically clean up downloaded file
- Download failures leave no partial files (atomic temp file usage)
- Context cancellation is respected throughout

---

### Package: manifest

**Import:** `github.com/grave0x/gpkg/internal/manifest`

**Purpose:** Parses and validates package manifest files (YAML format).

#### Types

##### Manifest

Top-level manifest structure defining package installation specs.

```go
type Manifest struct {
    Package      PackageSpec      `yaml:"package" json:"package"`
    Install      InstallSpec      `yaml:"install" json:"install"`
    BuildSource  *BuildSourceSpec `yaml:"build_source,omitempty" json:"build_source,omitempty"`
    Dependencies []string         `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
}
```

**Fields:**
- `Package` - Package metadata (name, version, author, etc.)
- `Install` - Release binary installation specification
- `BuildSource` - Optional source build specification
- `Dependencies` - List of package dependencies (package names or identifiers)

##### PackageSpec

Package metadata and description.

```go
type PackageSpec struct {
    Name     string `yaml:"name" json:"name"`
    Version  string `yaml:"version" json:"version"`
    Author   string `yaml:"author,omitempty" json:"author,omitempty"`
    URL      string `yaml:"url,omitempty" json:"url,omitempty"`
    License  string `yaml:"license,omitempty" json:"license,omitempty"`
}
```

**Required Fields:** `Name`, `Version`

##### InstallSpec

Describes how to install a release binary.

```go
type InstallSpec struct {
    Type        string            `yaml:"type" json:"type"`
    Source      string            `yaml:"source" json:"source"`
    Pattern     string            `yaml:"pattern,omitempty" json:"pattern,omitempty"`
    Checksum    map[string]string `yaml:"checksum,omitempty" json:"checksum,omitempty"`
    ExtractPath string            `yaml:"extract_path,omitempty" json:"extract_path,omitempty"`
    Executable  string            `yaml:"executable,omitempty" json:"executable,omitempty"`
    PostInstall string            `yaml:"post_install,omitempty" json:"post_install,omitempty"`
}
```

**Fields:**
- `Type` - Installation type: `"release"`, `"binary"`, or `"archive"`
- `Source` - Download URL or release pattern
- `Pattern` - Optional file pattern for release selection
- `Checksum` - Map of algorithm → hash (e.g., `{"sha256": "abc123..."}`)
- `ExtractPath` - Path within archive to extract
- `Executable` - Name of executable within archive
- `PostInstall` - Shell command to run after installation

##### BuildSourceSpec

Describes how to build from source.

```go
type BuildSourceSpec struct {
    Type     string            `yaml:"type" json:"type"`
    Source   string            `yaml:"source" json:"source"`
    Tag      string            `yaml:"tag,omitempty" json:"tag,omitempty"`
    Branch   string            `yaml:"branch,omitempty" json:"branch,omitempty"`
    Commands []string          `yaml:"commands" json:"commands"`
    Env      map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
}
```

**Fields:**
- `Type` - Source type: `"git"` or `"tarball"`
- `Source` - Repository URL or tarball URL
- `Tag` - Git tag to checkout (takes precedence over branch)
- `Branch` - Git branch to clone
- `Commands` - Build commands to execute (e.g., `["make", "make install"]`)
- `Env` - Environment variables for build process

##### Parser

Interface for manifest parsing and validation.

```go
type Parser interface {
    Parse(path string) (*Manifest, error)
    ParseBytes(data []byte) (*Manifest, error)
    Validate(m *Manifest) error
}
```

**Methods:**
- `Parse(path)` - Reads and parses manifest from file
- `ParseBytes(data)` - Parses manifest from bytes
- `Validate(m)` - Validates manifest structure and required fields

#### Implementations

##### YAMLParser

YAML-based manifest parser using gopkg.in/yaml.v3.

```go
type YAMLParser struct{}

func NewYAMLParser() *YAMLParser
```

**Validation Rules:**
- `package.name` and `package.version` are required
- Either `install` or `build_source` must be specified
- If `install` is present, `install.source` is required
- If `build_source` is present, `build_source.source` and at least one command are required

#### Usage Example

```go
import "github.com/grave0x/gpkg/internal/manifest"

// Parse manifest file
parser := manifest.NewYAMLParser()
mf, err := parser.Parse("/path/to/manifest.yaml")
if err != nil {
    log.Fatal(err)
}

// Parse from bytes
data := []byte(`
package:
  name: example
  version: 1.0.0
install:
  type: release
  source: https://github.com/user/repo/releases/latest
  checksum:
    sha256: abc123...
`)
mf, err = parser.ParseBytes(data)

// Validate manifest
if err := parser.Validate(mf); err != nil {
    log.Fatal("Invalid manifest:", err)
}

fmt.Printf("Package: %s v%s\n", mf.Package.Name, mf.Package.Version)
```

#### Error Handling

- File read errors: `"failed to read manifest file: %w"`
- YAML parse errors: `"failed to parse manifest YAML: %w"`
- Validation errors: `"manifest validation failed: %w"`
- Specific validation errors: `"package name is required"`, etc.

---

### Package: package

**Import:** `github.com/grave0x/gpkg/internal/package`

**Purpose:** Handles package installation from release binaries or source builds.

#### Types

##### Package

Represents a software package.

```go
type Package struct {
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Version     string            `json:"version"`
    Author      string            `json:"author,omitempty"`
    URL         string            `json:"url,omitempty"`
    License     string            `json:"license,omitempty"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}
```

##### InstalledPackage

Package with installation metadata.

```go
type InstalledPackage struct {
    Package      *Package
    InstalledAt  time.Time         `json:"installed_at"`
    Source       string            `json:"source"`          // "release" or "source"
    Prefix       string            `json:"prefix"`
    Checksums    map[string]string `json:"checksums,omitempty"`
    Dependencies []string          `json:"dependencies,omitempty"`
}
```

**Fields:**
- `Package` - Package metadata
- `InstalledAt` - Installation timestamp
- `Source` - Installation source: `"release"` (binary) or `"source"` (built)
- `Prefix` - Installation directory prefix
- `Checksums` - File checksums for verification
- `Dependencies` - List of installed dependencies

##### VersionInfo

Version information for a package.

```go
type VersionInfo struct {
    Latest      string    `json:"latest"`
    Current     string    `json:"current,omitempty"`
    Available   []string  `json:"available,omitempty"`
    ReleaseDate time.Time `json:"release_date,omitempty"`
}
```

##### Installer

Handles package installation operations.

```go
type Installer struct { /* private fields */ }

func NewInstaller(
    dl download.Downloader,
    atomic download.AtomicInstaller,
    prefix string,
) *Installer
```

**Parameters:**
- `dl` - Downloader for fetching release binaries
- `atomic` - AtomicInstaller for safe filesystem operations
- `prefix` - Installation prefix (e.g., `~/.gpkg`)

**Methods:**

```go
func (i *Installer) InstallFromRelease(
    ctx context.Context,
    mf *manifest.Manifest,
) (*InstalledPackage, error)
```

Installs a package from a release binary.

**Process:**
1. Downloads release from `manifest.Install.Source`
2. Validates checksum if provided
3. Atomically installs to `{prefix}/bin/{package-name}`
4. Returns `InstalledPackage` record

```go
func (i *Installer) InstallFromSource(
    ctx context.Context,
    mf *manifest.Manifest,
) (*InstalledPackage, error)
```

Installs a package by building from source.

**Process:**
1. Clones source repository to `{prefix}/src/{package-name}`
2. Checks out specified tag/branch
3. Executes build commands with environment variables
4. Returns `InstalledPackage` record

#### Usage Example

```go
import (
    "context"
    "github.com/grave0x/gpkg/internal/package"
    "github.com/grave0x/gpkg/internal/download"
    "github.com/grave0x/gpkg/internal/manifest"
)

// Create installer
dl := download.NewHTTPDownloader(30*time.Second, false)
atomic := download.NewAtomicInstaller("/home/user/.gpkg/backups")
installer := pkg.NewInstaller(dl, atomic, "/home/user/.gpkg")

// Parse manifest
parser := manifest.NewYAMLParser()
mf, _ := parser.Parse("manifest.yaml")

// Install from release
ctx := context.Background()
installed, err := installer.InstallFromRelease(ctx, mf)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Installed %s v%s at %s\n",
    installed.Package.Name,
    installed.Package.Version,
    installed.InstalledAt)

// Install from source
installed, err = installer.InstallFromSource(ctx, mf)
```

#### Error Handling

- Missing specification: `"manifest does not have install specification"`
- Download failures: `"failed to download with checksum validation: %w"`
- Installation failures: `"failed to install package: %w"`
- Build failures: `"build command failed: %w"` (includes command output)
- Git clone failures: `"git clone failed: <output>"`

---

### Package: pkgdb

**Import:** `github.com/grave0x/gpkg/internal/pkgdb`

**Purpose:** Manages SQLite database operations for tracking installed packages, versions, and files.

#### Types

##### PackageRecord

Database record for an installed package.

```go
type PackageRecord struct {
    ID            int64
    Name          string
    Version       string
    InstalledAt   time.Time
    UpdatedAt     time.Time
    Source        string                // "release" or "source"
    Prefix        string
    Author        string
    URL           string
    License       string
    Checksums     map[string]string
    Dependencies  []string
    Files         []string
    BuildMetadata map[string]string
}
```

**Fields:**
- `ID` - Auto-incremented primary key
- `Name` - Package name (unique)
- `Version` - Currently installed version
- `InstalledAt` - Initial installation timestamp
- `UpdatedAt` - Last update timestamp
- `Source` - `"release"` or `"source"`
- `Prefix` - Installation directory
- `Checksums` - File checksums (JSON-encoded in DB)
- `Dependencies` - List of dependencies (JSON-encoded)
- `Files` - List of installed file paths
- `BuildMetadata` - Additional metadata (JSON-encoded)

##### VersionRecord

Historical version record for rollback support.

```go
type VersionRecord struct {
    ID          int64
    PackageID   int64
    Version     string
    InstalledAt time.Time
    Checksums   map[string]string
    Files       []string
}
```

##### FileRecord

Individual file installed by a package.

```go
type FileRecord struct {
    ID        int64
    PackageID int64
    FilePath  string
    Checksum  string
    FileSize  int64
}
```

##### Manager

Interface for database operations.

```go
type Manager interface {
    AddPackage(p *PackageRecord) (int64, error)
    GetPackage(name string) (*PackageRecord, error)
    UpdatePackage(p *PackageRecord) error
    DeletePackage(name string) error
    ListPackages() ([]*PackageRecord, error)
    
    AddFiles(packageID int64, files []string) error
    GetFiles(packageID int64) ([]string, error)
    
    AddVersion(packageID int64, v *VersionRecord) (int64, error)
    GetVersionHistory(packageID int64) ([]*VersionRecord, error)
    
    Close() error
}
```

**Methods:**
- `AddPackage(p)` - Inserts new package, returns ID
- `GetPackage(name)` - Retrieves package by name
- `UpdatePackage(p)` - Updates package metadata
- `DeletePackage(name)` - Removes package and cascades to files/versions
- `ListPackages()` - Returns all installed packages
- `AddFiles(id, files)` - Records installed files
- `GetFiles(id)` - Retrieves file list for package
- `AddVersion(id, v)` - Records version in history
- `GetVersionHistory(id)` - Returns version history ordered by date
- `Close()` - Closes database connection

#### Implementations

##### SQLiteManager

SQLite-based package database.

```go
type SQLiteManager struct { /* private fields */ }

func NewSQLiteManager(dbPath string) (*SQLiteManager, error)
```

**Parameters:**
- `dbPath` - Path to SQLite database file (e.g., `~/.gpkg/packages.db`)

**Features:**
- Auto-creates schema on first use
- JSON encoding for complex fields (checksums, dependencies, metadata)
- Foreign key constraints with CASCADE delete
- Indexed queries for performance

#### Database Schema

See [Database Schema](#database-schema) section below.

#### Usage Example

```go
import "github.com/grave0x/gpkg/internal/pkgdb"

// Open database
db, err := pkgdb.NewSQLiteManager("/home/user/.gpkg/packages.db")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Add package
pkg := &pkgdb.PackageRecord{
    Name:    "example",
    Version: "1.0.0",
    Source:  "release",
    Prefix:  "/home/user/.gpkg",
    Checksums: map[string]string{
        "sha256": "abc123...",
    },
}
id, err := db.AddPackage(pkg)

// Record installed files
files := []string{
    "/home/user/.gpkg/bin/example",
    "/home/user/.gpkg/share/example/config.yaml",
}
db.AddFiles(id, files)

// List all packages
packages, err := db.ListPackages()
for _, p := range packages {
    fmt.Printf("%s v%s (%d files)\n", p.Name, p.Version, len(p.Files))
}

// Get specific package
pkg, err = db.GetPackage("example")
if err != nil {
    log.Fatal(err)
}

// Update package
pkg.Version = "1.1.0"
db.UpdatePackage(pkg)

// Delete package (cascades to files and versions)
db.DeletePackage("example")
```

#### Error Handling

- Package not found: `"package not found: <name>"`
- Database errors: `"failed to add package: %w"`
- Foreign key violations are prevented by schema constraints
- Duplicate package names return error: unique constraint on `name`

---

### Package: planner

**Import:** `github.com/grave0x/gpkg/internal/planner`

**Purpose:** Generates and validates installation plans before execution.

#### Types

##### Action

Single step in an installation plan.

```go
type Action struct {
    Type        string `json:"type"`
    Description string `json:"description"`
    Size        int64  `json:"size,omitempty"`
    Duration    int64  `json:"duration,omitempty"`
    Required    bool   `json:"required"`
    Reversible  bool   `json:"reversible"`
}
```

**Fields:**
- `Type` - Action type: `"download"`, `"extract"`, `"build"`, `"write"`, `"cleanup"`
- `Description` - Human-readable description
- `Size` - Estimated size in bytes
- `Duration` - Estimated time in seconds
- `Required` - Whether action is required for installation
- `Reversible` - Whether action can be rolled back

##### Plan

Complete installation plan.

```go
type Plan struct {
    Package         string    `json:"package"`
    Version         string    `json:"version"`
    Actions         []*Action `json:"actions"`
    TotalSize       int64     `json:"total_size"`
    EstimatedTime   int64     `json:"estimated_time_seconds"`
    Dependencies    []string  `json:"dependencies,omitempty"`
    Conflicts       []string  `json:"conflicts,omitempty"`
    Warnings        []string  `json:"warnings,omitempty"`
    WillReplace     bool      `json:"will_replace"`
    PreviousVersion string    `json:"previous_version,omitempty"`
}
```

**Fields:**
- `Package` - Package name
- `Version` - Target version
- `Actions` - Ordered list of actions to perform
- `TotalSize` - Total download/disk size
- `EstimatedTime` - Total estimated time in seconds
- `Dependencies` - List of dependency package names
- `Conflicts` - List of conflicting packages
- `Warnings` - User warnings (missing dependencies, etc.)
- `WillReplace` - Whether upgrading existing installation
- `PreviousVersion` - Version being replaced (if upgrading)

##### Planner

Interface for generating installation plans.

```go
type Planner interface {
    PlanInstallation(ctx context.Context, mf *manifest.Manifest, fromRelease bool) (*Plan, error)
    PlanUpgrade(ctx context.Context, current, target *manifest.Manifest) (*Plan, error)
    PlanUninstall(ctx context.Context, pkg string) (*Plan, error)
    ValidatePlan(p *Plan) error
}
```

**Methods:**
- `PlanInstallation(ctx, mf, fromRelease)` - Creates installation plan
- `PlanUpgrade(ctx, current, target)` - Creates upgrade plan with backup
- `PlanUninstall(ctx, pkg)` - Creates uninstallation plan
- `ValidatePlan(p)` - Validates plan is executable

#### Implementations

##### DefaultPlanner

Default planner implementation.

```go
type DefaultPlanner struct { /* private fields */ }

func NewDefaultPlanner(offline, dryRun bool) *DefaultPlanner
```

**Parameters:**
- `offline` - Offline mode (affects plan validation)
- `dryRun` - Dry-run mode (for simulation)

#### Usage Example

```go
import (
    "context"
    "github.com/grave0x/gpkg/internal/planner"
    "github.com/grave0x/gpkg/internal/manifest"
)

// Create planner
pl := planner.NewDefaultPlanner(false, false)

// Parse manifest
parser := manifest.NewYAMLParser()
mf, _ := parser.Parse("manifest.yaml")

// Plan installation
ctx := context.Background()
plan, err := pl.PlanInstallation(ctx, mf, true)
if err != nil {
    log.Fatal(err)
}

// Display plan
fmt.Printf("Installing %s v%s\n", plan.Package, plan.Version)
fmt.Printf("Total size: %d MB\n", plan.TotalSize/1024/1024)
fmt.Printf("Estimated time: %d seconds\n", plan.EstimatedTime)

for i, action := range plan.Actions {
    fmt.Printf("%d. [%s] %s\n", i+1, action.Type, action.Description)
}

if len(plan.Warnings) > 0 {
    fmt.Println("\nWarnings:")
    for _, w := range plan.Warnings {
        fmt.Printf("  - %s\n", w)
    }
}

// Validate plan
if err := pl.ValidatePlan(plan); err != nil {
    log.Fatal("Invalid plan:", err)
}

// Plan upgrade
currentMf, _ := parser.Parse("current-manifest.yaml")
upgradePlan, err := pl.PlanUpgrade(ctx, currentMf, mf)
fmt.Printf("Upgrading from %s to %s\n",
    upgradePlan.PreviousVersion,
    upgradePlan.Version)
```

#### Error Handling

- Invalid plans: `"invalid plan: package name is empty"`
- Missing actions: `"invalid plan: no actions defined"`
- Invalid action types: `"invalid plan: action has no type"`

---

### Package: resolver

**Import:** `github.com/grave0x/gpkg/internal/resolver`

**Purpose:** Resolves package identifiers to manifests and handles dependency resolution.

#### Types

##### PackageResolver

Interface for resolving package identifiers.

```go
type PackageResolver interface {
    ResolvePackage(ctx context.Context, identifier string) (*manifest.Manifest, error)
    ResolveGitHub(ctx context.Context, owner, repo string) (*manifest.Manifest, error)
    ParseIdentifier(identifier string) (Type, Owner, Repo string, err error)
}
```

**Methods:**
- `ResolvePackage(ctx, id)` - Resolves identifier to manifest
- `ResolveGitHub(ctx, owner, repo)` - Resolves GitHub repository
- `ParseIdentifier(id)` - Parses identifier format

**Supported Identifier Formats:**
- `owner/repo` - GitHub shorthand
- `github:owner/repo` - Explicit GitHub
- `https://github.com/owner/repo` - Full GitHub URL

##### GitHubResolver

GitHub-based package resolver.

```go
type GitHubResolver struct { /* private fields */ }

func NewGitHubResolver() *GitHubResolver
```

**Features:**
- Fetches manifest from repository (manifest.yaml or gpkg.yaml)
- Retrieves latest release information
- Constructs manifest from release metadata

##### DependencyResolver

Resolves package dependencies recursively.

```go
type DependencyResolver struct { /* private fields */ }

func NewDependencyResolver(pr PackageResolver) *DependencyResolver
```

**Methods:**

```go
func (dr *DependencyResolver) ResolveDependencies(
    ctx context.Context,
    mf *manifest.Manifest,
) ([]string, error)
```

Recursively resolves all dependencies, returns installation order.

```go
func (dr *DependencyResolver) CheckConflicts(
    ctx context.Context,
    mf *manifest.Manifest,
) ([]string, error)
```

Checks for version or file conflicts with installed packages.

#### Usage Example

```go
import (
    "context"
    "github.com/grave0x/gpkg/internal/resolver"
)

// Create resolver
res := resolver.NewGitHubResolver()

// Resolve package identifier
ctx := context.Background()
mf, err := res.ResolvePackage(ctx, "golang/go")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Resolved: %s v%s\n", mf.Package.Name, mf.Package.Version)

// Parse identifier
pkgType, owner, repo, err := res.ParseIdentifier("github:user/repo")
fmt.Printf("Type: %s, Owner: %s, Repo: %s\n", pkgType, owner, repo)

// Resolve dependencies
depResolver := resolver.NewDependencyResolver(res)
order, err := depResolver.ResolveDependencies(ctx, mf)
fmt.Printf("Install order: %v\n", order)

// Check conflicts
conflicts, err := depResolver.CheckConflicts(ctx, mf)
if len(conflicts) > 0 {
    fmt.Println("Conflicts detected:", conflicts)
}
```

#### Error Handling

- Invalid identifiers: `"invalid package identifier: <id>"`
- Unsupported formats: `"unsupported identifier format: <id>"`
- Invalid GitHub URLs: `"invalid GitHub URL: <url>"`
- Unsupported package types: `"unsupported package type: <type>"`

---

### Package: source

**Import:** `github.com/grave0x/gpkg/internal/source`

**Purpose:** Manages package source registries (indices where packages can be found).

#### Types

##### Source

Represents a package source/registry.

```go
type Source struct {
    ID          string `json:"id"`
    URI         string `json:"uri"`
    Name        string `json:"name,omitempty"`
    Description string `json:"description,omitempty"`
    LastUpdated int64  `json:"last_updated,omitempty"`
    Enabled     bool   `json:"enabled"`
    SourceType  string `json:"type,omitempty"`  // "http", "local", "github"
}
```

**Fields:**
- `ID` - Unique identifier (slug)
- `URI` - Source URL or path
- `Name` - Display name
- `Description` - Human-readable description
- `LastUpdated` - Unix timestamp of last update
- `Enabled` - Whether source is active
- `SourceType` - Source type: `"http"`, `"local"`, or `"github"`

##### Registry

Interface for managing sources.

```go
type Registry interface {
    AddSource(ctx context.Context, src *Source) error
    RemoveSource(ctx context.Context, idOrURI string) error
    ListSources(ctx context.Context) ([]*Source, error)
    GetSource(ctx context.Context, id string) (*Source, error)
    UpdateSource(ctx context.Context, src *Source) error
}
```

**Methods:**
- `AddSource(ctx, src)` - Adds new source to registry
- `RemoveSource(ctx, idOrURI)` - Removes by ID or URI
- `ListSources(ctx)` - Returns all sources
- `GetSource(ctx, id)` - Retrieves specific source
- `UpdateSource(ctx, src)` - Updates source metadata

##### Fetcher

Interface for retrieving package information from sources.

```go
type Fetcher interface {
    Fetch(ctx context.Context, src *Source) (interface{}, error)
    FetchPackage(ctx context.Context, src *Source, pkgName string) (interface{}, error)
}
```

**Methods:**
- `Fetch(ctx, src)` - Fetches all package metadata from source
- `FetchPackage(ctx, src, name)` - Fetches specific package info

#### Implementations

##### JSONRegistry

JSON file-based source registry.

```go
type JSONRegistry struct { /* private fields */ }

func NewJSONRegistry(filePath string) *JSONRegistry
```

**Parameters:**
- `filePath` - Path to sources JSON file (e.g., `~/.gpkg/sources.json`)

**Additional Methods:**

```go
func (r *JSONRegistry) Load() error
func (r *JSONRegistry) Save() error
```

- `Load()` - Reads registry from disk
- `Save()` - Writes registry to disk (called automatically by mutating methods)

#### Usage Example

```go
import (
    "context"
    "github.com/grave0x/gpkg/internal/source"
)

// Create registry
registry := source.NewJSONRegistry("/home/user/.gpkg/sources.json")
if err := registry.Load(); err != nil {
    log.Fatal(err)
}

// Add source
ctx := context.Background()
src := &source.Source{
    ID:          "official",
    URI:         "https://packages.gpkg.dev/index.json",
    Name:        "Official gpkg Registry",
    Description: "Official package registry",
    SourceType:  "http",
}
err := registry.AddSource(ctx, src)

// List sources
sources, err := registry.ListSources(ctx)
for _, s := range sources {
    fmt.Printf("[%s] %s - %s\n", s.ID, s.Name, s.URI)
}

// Get specific source
src, err = registry.GetSource(ctx, "official")
if err != nil {
    log.Fatal(err)
}

// Update source
src.Description = "Updated description"
registry.UpdateSource(ctx, src)

// Remove source
registry.RemoveSource(ctx, "official")
```

#### Error Handling

- Empty ID: `"source ID cannot be empty"`
- Duplicate source: `"source with ID <id> already exists"`
- Source not found: `"source not found: <id>"`
- File I/O errors: `"failed to read registry file: %w"`
- JSON errors: `"failed to parse registry JSON: %w"`

---

## Data Models

### Core Structures Summary

#### Package Information
```go
// Runtime package representation
Package {
    Name, Description, Version
    Author, URL, License
    Metadata map[string]string
}

// Database record with installation info
PackageRecord {
    ID, Name, Version
    InstalledAt, UpdatedAt
    Source, Prefix, Author, URL, License
    Checksums, Dependencies, Files
    BuildMetadata
}
```

#### Manifest Specifications
```go
// Release installation
InstallSpec {
    Type: "release" | "binary" | "archive"
    Source: URL
    Checksum: {algo: hash}
    ExtractPath, Executable, PostInstall
}

// Source build
BuildSourceSpec {
    Type: "git" | "tarball"
    Source: URL
    Tag, Branch
    Commands: []string
    Env: {key: value}
}
```

#### Installation Planning
```go
// Installation plan
Plan {
    Package, Version
    Actions: []{Type, Description, Size, Duration}
    TotalSize, EstimatedTime
    Dependencies, Conflicts, Warnings
    WillReplace, PreviousVersion
}

// Single action
Action {
    Type: "download" | "extract" | "build" | "write" | "cleanup"
    Description: string
    Required, Reversible: bool
}
```

#### Source Registry
```go
Source {
    ID: "unique-id"
    URI: "https://..." | "/path/to/local"
    Name, Description
    SourceType: "http" | "local" | "github"
    Enabled: bool
}
```

---

## Database Schema

### SQLite Database Structure

**Location:** `~/.gpkg/packages.db` (configurable via `Config.Prefix`)

#### Table: packages

Primary package records.

```sql
CREATE TABLE packages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    version TEXT NOT NULL,
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    source TEXT,                    -- "release" or "source"
    prefix TEXT,                    -- Installation prefix
    author TEXT,
    url TEXT,
    license TEXT,
    checksums TEXT,                 -- JSON: {"sha256": "abc..."}
    dependencies TEXT,              -- JSON: ["pkg1", "pkg2"]
    build_metadata TEXT             -- JSON: additional metadata
);

CREATE INDEX idx_packages_name ON packages(name);
```

**Fields:**
- `id` - Auto-increment primary key
- `name` - Package name (unique constraint)
- `version` - Current installed version
- `installed_at` - Initial installation timestamp
- `updated_at` - Last update/upgrade timestamp
- `source` - `"release"` or `"source"`
- `checksums` - JSON-encoded map: `{"sha256": "abc123..."}`
- `dependencies` - JSON-encoded array: `["dep1", "dep2"]`
- `build_metadata` - JSON-encoded metadata map

#### Table: versions

Version history for rollback support.

```sql
CREATE TABLE versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    package_id INTEGER NOT NULL,
    version TEXT NOT NULL,
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    checksums TEXT,                 -- JSON
    files TEXT,                     -- JSON: list of files
    FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE
);

CREATE INDEX idx_versions_package_id ON versions(package_id);
```

**Fields:**
- `package_id` - Foreign key to packages.id
- `version` - Historical version string
- `checksums` - JSON-encoded checksums for this version
- `files` - JSON-encoded file list for this version

#### Table: files

Individual files installed by packages.

```sql
CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    package_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    checksum TEXT,
    file_size INTEGER,
    FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE,
    UNIQUE(package_id, file_path)
);

CREATE INDEX idx_files_package_id ON files(package_id);
```

**Fields:**
- `package_id` - Foreign key to packages.id
- `file_path` - Absolute file path
- `checksum` - SHA256 checksum of file
- `file_size` - File size in bytes

#### Relationships

```
packages (1) ──< (many) versions
packages (1) ──< (many) files
```

- **CASCADE DELETE:** Deleting a package removes all associated versions and files
- **UNIQUE CONSTRAINT:** One package name, one file per package
- **INDEXES:** Fast lookups by package name and package_id

#### Query Examples

```sql
-- Get package with files
SELECT p.*, GROUP_CONCAT(f.file_path) as files
FROM packages p
LEFT JOIN files f ON p.id = f.package_id
WHERE p.name = 'example'
GROUP BY p.id;

-- Get version history
SELECT version, installed_at
FROM versions
WHERE package_id = (SELECT id FROM packages WHERE name = 'example')
ORDER BY installed_at DESC;

-- Find packages by source type
SELECT name, version, source
FROM packages
WHERE source = 'release'
ORDER BY installed_at DESC;

-- Check for file conflicts
SELECT p.name, f.file_path
FROM files f
JOIN packages p ON f.package_id = p.id
WHERE f.file_path = '/usr/local/bin/tool';
```

---

## Error Handling

### Error Conventions

All packages follow consistent error handling patterns:

#### Wrapping Errors

Errors are wrapped with context using `fmt.Errorf` with `%w`:

```go
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

This allows:
- Error unwrapping with `errors.Unwrap()`
- Error checking with `errors.Is()` and `errors.As()`
- Contextual error messages

#### Common Error Patterns

```go
// Not found errors
fmt.Errorf("package not found: %s", name)
fmt.Errorf("source not found: %q", id)

// Validation errors
fmt.Errorf("package name is required")
fmt.Errorf("invalid plan: %s", reason)

// I/O errors
fmt.Errorf("failed to read config file: %w", err)
fmt.Errorf("failed to create directory: %w", err)

// Network errors
fmt.Errorf("download failed with status %d", statusCode)
fmt.Errorf("failed to download file: %w", err)

// Database errors
fmt.Errorf("failed to add package: %w", err)
fmt.Errorf("failed to query packages: %w", err)
```

### Error Checking

```go
// Check specific error type
var notFoundErr *NotFoundError
if errors.As(err, &notFoundErr) {
    // Handle not found
}

// Check wrapped error
if errors.Is(err, sql.ErrNoRows) {
    // Handle no rows
}

// Unwrap error
cause := errors.Unwrap(err)
```

### Context Cancellation

All context-aware functions respect cancellation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := downloader.Download(ctx, url, dest)
if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("Download timed out")
}
if errors.Is(err, context.Canceled) {
    fmt.Println("Download cancelled")
}
```

---

## Context Usage

### Context Patterns

All async operations accept `context.Context` for:
- Timeout control
- Cancellation signals
- Request-scoped values (in future)

#### Basic Usage

```go
// Background context (no timeout)
ctx := context.Background()
err := installer.InstallFromRelease(ctx, manifest)

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
err := downloader.Download(ctx, url, dest)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    // Cancel on signal
    <-signalChan
    cancel()
}()
err := installer.InstallFromSource(ctx, manifest)
```

#### Context Propagation

Always propagate context through call chains:

```go
func InstallPackage(ctx context.Context, name string) error {
    mf, err := resolver.ResolvePackage(ctx, name)  // Pass ctx
    if err != nil {
        return err
    }
    
    plan, err := planner.PlanInstallation(ctx, mf, true)  // Pass ctx
    if err != nil {
        return err
    }
    
    return installer.InstallFromRelease(ctx, mf)  // Pass ctx
}
```

---

## Environment Variables

### Configuration Overrides

Environment variables override config file settings:

| Variable | Description | Default |
|----------|-------------|---------|
| `GPKG_PREFIX` | Installation prefix | `~/.gpkg` |
| `GPKG_CACHE_DIR` | Cache directory | `~/.gpkg/cache` |
| `GPKG_LOG_LEVEL` | Logging level | `info` |

### Usage

```bash
# Custom installation prefix
export GPKG_PREFIX=/opt/gpkg
gpkg install example

# Debug logging
export GPKG_LOG_LEVEL=debug
gpkg install example

# Custom cache directory
export GPKG_CACHE_DIR=/var/cache/gpkg
gpkg update
```

### Precedence Order

Configuration is loaded with the following precedence (highest to lowest):

1. **CLI flags** (e.g., `--prefix /custom`)
2. **Environment variables** (`GPKG_PREFIX`, etc.)
3. **User config files**:
   - `~/.gpkg/config.yaml`
   - `~/.gpkg.yaml`
   - `~/.config/gpkg/config.yaml`
4. **System config files**:
   - `/etc/gpkg/config.yaml`
   - `/etc/gpkg.yaml`
5. **Built-in defaults**

---

## Appendix

### Supported Checksum Algorithms

| Algorithm | Type | Recommended |
|-----------|------|-------------|
| `sha256` | SHA-256 | ✅ Yes |
| `sha1` | SHA-1 | ⚠️ Legacy only |
| `md5` | MD5 | ⚠️ Legacy only |

**Recommendation:** Always use `sha256` for new manifests.

### Default File Locations

| Purpose | Path | Configurable |
|---------|------|--------------|
| Config file | `~/.gpkg/config.yaml` | ✅ `--config` |
| Sources registry | `~/.gpkg/sources.json` | ✅ `SourcesFile` |
| Package database | `~/.gpkg/packages.db` | ✅ `Prefix` |
| Binaries | `~/.gpkg/bin/` | ✅ `Prefix` |
| Cache | `~/.gpkg/cache/` | ✅ `CacheDir` |
| Source builds | `~/.gpkg/src/` | ✅ `Prefix` |
| Backups | `~/.gpkg/backups/` | ✅ `Prefix` |

### Manifest File Formats

**Primary:** `manifest.yaml` or `gpkg.yaml` in repository root

**Example:**

```yaml
package:
  name: example-tool
  version: 1.2.3
  author: John Doe
  url: https://github.com/user/example
  license: MIT

install:
  type: release
  source: https://github.com/user/example/releases/download/v1.2.3/binary
  checksum:
    sha256: abc123def456...
  executable: example-tool

dependencies:
  - dependency-pkg

build_source:
  type: git
  source: https://github.com/user/example.git
  tag: v1.2.3
  commands:
    - make build
    - make install PREFIX=$GPKG_PREFIX
  env:
    CGO_ENABLED: "0"
```

---

## See Also

- [CLI Specification](../CLI+SPEC.md) - Command-line interface documentation
- [Implementation Summary](../IMPLEMENTATION_SUMMARY.md) - Implementation overview
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute to gpkg
- [Development Guide](../DEVELOPMENT.md) - Setting up development environment

---

**Last Updated:** 2025-01-XX  
**Version:** 1.0.0
