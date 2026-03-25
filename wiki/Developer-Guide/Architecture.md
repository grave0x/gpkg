# Developer Architecture Guide

This guide explains the architecture of gpkg, a Git-based package manager written in Go. It covers the project structure, package organization, data flow, key components, database schema, design patterns, and dependencies.

## Table of Contents

- [Project Structure](#project-structure)
- [Package Organization](#package-organization)
- [Data Flow](#data-flow)
- [Key Components](#key-components)
- [Database Schema](#database-schema)
- [Design Patterns](#design-patterns)
- [Dependencies](#dependencies)
- [Configuration System](#configuration-system)
- [Architectural Strengths](#architectural-strengths)

---

## Project Structure

gpkg follows Go best practices with a clean separation between public APIs, internal implementation, and command-line interface:

```
/home/grave/Documents/git-repos/gpkg/
├── cmd/gpkg/cmd/           # CLI Command handlers (Cobra commands)
├── internal/                # Core packages (non-exported)
│   ├── config/             # Configuration management
│   ├── download/           # Download & atomic operations
│   ├── manifest/           # Manifest parsing
│   ├── package/            # Package installation logic
│   ├── pkgdb/              # SQLite database layer
│   ├── planner/            # Installation planning
│   ├── resolver/           # Package resolution
│   └── source/             # Package source registry
├── pkg/                    # Public API packages (reserved for future use)
├── main.go                 # Entry point
├── go.mod                  # Dependencies
├── tests/                  # Test files
├── bin/                    # Build output
├── examples/               # Example manifests
└── wiki/                   # Documentation
```

### Design Principle

The directory structure enforces **separation of concerns**:

- **`cmd/`** - User interface layer (CLI commands and flags)
- **`internal/`** - Business logic (non-exported, encapsulated)
- **`pkg/`** - Public APIs for external consumers (reserved for future)
- **`main.go`** - Minimal entry point that delegates to cmd layer

This structure prevents accidental coupling and makes testing easier by isolating the CLI from core logic.

---

## Package Organization

Each package in `internal/` has a focused responsibility. Here's a detailed breakdown:

### 1. `internal/config` - Configuration Management

**Purpose**: Load and merge configuration from multiple sources with precedence rules.

**Main Types**:
```go
type Config struct {
    Prefix            string // Installation prefix (~/.gpkg)
    CacheDir          string // Cache location
    SourcesFile       string // Sources registry path
    LogLevel          string // Logging verbosity
    Color             bool   // Colorized output
    StrictChecksum    bool   // Strict checksum validation
    NetworkTimeout    int    // Timeout in seconds
}

type Loader interface {
    Load() (*Config, error)
    LoadFrom(path string) (*Config, error)
    MergeDefaults(cfg *Config) *Config
}
```

**Implementation**: `YAMLLoader`

**Configuration Precedence** (highest to lowest):
1. CLI flags (`--prefix`, `--config`, etc.)
2. Environment variables (`GPKG_PREFIX`, `GPKG_CACHE_DIR`, `GPKG_LOG_LEVEL`)
3. User config (`~/.gpkg/config.yaml`, `~/.config/gpkg/config.yaml`)
4. System config (`/etc/gpkg/config.yaml`)
5. Hardcoded defaults

**Files**: `config.go` (56 lines), `loader.go` (124 lines)

---

### 2. `internal/manifest` - Manifest Parsing

**Purpose**: Parse and validate package manifests (YAML format).

**Main Types**:
```go
type Manifest struct {
    Package      PackageSpec         // Package metadata
    Install      InstallSpec         // Binary release installation
    BuildSource  *BuildSourceSpec    // Source build (optional)
    Dependencies []string            // Package dependencies
}

type PackageSpec struct {
    Name     string  // Package name (required)
    Version  string  // Version (required)
    Author   string
    URL      string  // Homepage/repository
    License  string
}

type InstallSpec struct {
    Type        string            // "release", "binary", "archive"
    Source      string            // Download URL
    Pattern     string            // File pattern (optional)
    Checksum    map[string]string // Algorithm -> Hash
    ExtractPath string            // Path in archive
    Executable  string            // Executable name
    PostInstall string            // Post-install script
}

type BuildSourceSpec struct {
    Type     string            // "git", "tarball"
    Source   string            // Repository URL
    Tag      string
    Branch   string
    Commands []string          // Build commands
    Env      map[string]string // Environment variables
}

type Parser interface {
    Parse(path string) (*Manifest, error)
    ParseBytes(data []byte) (*Manifest, error)
    Validate(m *Manifest) error
}
```

**Validation Rules**:
- Package name and version are required
- Either `install` or `build_source` must be specified
- Install requires `type` and `source` URL
- BuildSource requires at least one build command

**Files**: `manifest.go` (51 lines), `parser.go` (68 lines)

---

### 3. `internal/package` - Installation Logic

**Purpose**: Install packages from binary releases or source builds.

**Main Types**:
```go
type Package struct {
    Name        string
    Description string
    Version     string
    Author      string
    URL         string
    License     string
    Metadata    map[string]string
}

type InstalledPackage struct {
    Package      *Package
    InstalledAt  time.Time
    Source       string            // "release" or "source"
    Prefix       string            // Installation prefix
    Checksums    map[string]string
    Dependencies []string
}

type Installer struct {
    downloader download.Downloader
    atomicOp   download.AtomicInstaller
    prefix     string
}
```

**Key Methods**:
- `NewInstaller(downloader, atomicOp, prefix)` - Create installer with dependencies
- `InstallFromRelease(ctx, manifest)` - Install binary release
  - Downloads from URL
  - Validates checksum (SHA256/SHA1/MD5)
  - Performs atomic installation to `~/.gpkg/bin/{name}`
- `InstallFromSource(ctx, manifest)` - Build from source
  - Clones git repository or downloads tarball
  - Runs build commands
  - Stores in `~/.gpkg/src/{name}`

**Files**: `package.go` (32 lines), `installer.go` (183 lines)

---

### 4. `internal/pkgdb` - SQLite Database Layer

**Purpose**: Persist package metadata, track versions, and maintain file inventories.

**Database Schema**:

```sql
-- Installed packages (current state)
CREATE TABLE packages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    version TEXT NOT NULL,
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    source TEXT,              -- "release" or "source"
    prefix TEXT,
    author TEXT,
    url TEXT,
    license TEXT,
    checksums TEXT,           -- JSON: {algo: hash}
    dependencies TEXT,        -- JSON: [dep1, dep2]
    build_metadata TEXT       -- JSON: custom metadata
);

-- Version history for rollback
CREATE TABLE versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    package_id INTEGER NOT NULL,
    version TEXT NOT NULL,
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    checksums TEXT,
    files TEXT,               -- JSON: [file1, file2]
    FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE
);

-- File tracking for uninstall
CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    package_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    checksum TEXT,
    file_size INTEGER,
    FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE,
    UNIQUE(package_id, file_path)
);

CREATE INDEX idx_packages_name ON packages(name);
CREATE INDEX idx_versions_package_id ON versions(package_id);
CREATE INDEX idx_files_package_id ON files(package_id);
```

**Types**:
```go
type PackageRecord struct {
    ID            int64
    Name          string
    Version       string
    InstalledAt   time.Time
    UpdatedAt     time.Time
    Source        string            // "release" or "source"
    Prefix        string
    Author        string
    URL           string
    License       string
    Checksums     map[string]string
    Dependencies  []string
    Files         []string
    BuildMetadata map[string]string
}

type VersionRecord struct {
    ID          int64
    PackageID   int64
    Version     string
    InstalledAt time.Time
    Checksums   map[string]string
    Files       []string
}

type FileRecord struct {
    ID        int64
    PackageID int64
    FilePath  string
    Checksum  string
    FileSize  int64
}

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

**Implementation**: `SQLiteManager`
- Database location: `~/.gpkg/pkgdb.sqlite`
- JSON serialization for complex fields (checksums, dependencies, metadata)
- Cascading deletes maintain referential integrity
- Indexes optimize common queries (lookup by name, version history, file lists)

**Files**: `types.go` (73 lines), `sqlite.go` (292 lines)

---

### 5. `internal/download` - Download & Atomic Operations

**Purpose**: Download files from URLs and perform atomic installations with backups.

**Interfaces**:
```go
type Downloader interface {
    Download(ctx context.Context, url, dest string) error
    DownloadWithChecksum(ctx, url, dest, expectedChecksum, algorithm string) error
    ValidateChecksum(filePath, expectedChecksum, algorithm string) (bool, error)
}

type AtomicInstaller interface {
    Install(ctx context.Context, sourceFile, destPath string) error
    Uninstall(ctx context.Context, pkgName string) error
    Rollback(ctx context.Context, installID string) error
}
```

**Implementation**: `HTTPDownloader`
- HTTP downloads with context support (cancellation, timeouts)
- Atomic temp file operations (download to temp, then rename)
- Checksum validation (SHA256, SHA1, MD5)
- Offline mode support (prevents network calls)
- Progress reporting (size tracking)

**Implementation**: `AtomicInstallerImpl`
- Creates timestamped backups before installation
- Atomic file operations (write to temp, then rename)
- Sets executable permissions (0755)
- Backup directory: `~/.gpkg/backups/{timestamp}`

**Files**: `download.go` (36 lines), `http_downloader.go` (192 lines)

---

### 6. `internal/resolver` - Package Resolution

**Purpose**: Resolve package identifiers to manifests and handle dependency graphs.

**Interfaces**:
```go
type PackageResolver interface {
    ResolvePackage(ctx context.Context, identifier string) (*manifest.Manifest, error)
    ResolveGitHub(ctx context.Context, owner, repo string) (*manifest.Manifest, error)
    ParseIdentifier(identifier string) (Type, Owner, Repo string, err error)
}

type DependencyResolver struct {
    resolver PackageResolver
    cache    map[string]*manifest.Manifest
}
```

**Identifier Formats Supported**:
- `owner/repo` (GitHub shorthand)
- `github:owner/repo` (explicit GitHub)
- `https://github.com/owner/repo` (full URL)

**Features**:
- Recursive dependency traversal
- In-memory caching to avoid duplicate resolution
- Cycle detection (prevents infinite loops)
- Conflict detection (version/file conflicts)

**Files**: `resolver.go` (155 lines)

---

### 7. `internal/planner` - Installation Planning

**Purpose**: Create execution plans for install/upgrade/uninstall operations.

**Types**:
```go
type Action struct {
    Type        string // "download", "extract", "build", "write", "cleanup"
    Description string
    Size        int64  // bytes
    Duration    int64  // estimated seconds
    Required    bool   // Must succeed
    Reversible  bool   // Can be rolled back
}

type Plan struct {
    Package         string
    Version         string
    Actions         []*Action  // Ordered action sequence
    TotalSize       int64      // Total bytes
    EstimatedTime   int64      // Seconds
    Dependencies    []string
    Conflicts       []string
    Warnings        []string
    WillReplace     bool       // Overwriting existing?
    PreviousVersion string
}

type Planner interface {
    PlanInstallation(ctx context.Context, mf *manifest.Manifest, fromRelease bool) (*Plan, error)
    PlanUpgrade(ctx context.Context, current, target *manifest.Manifest) (*Plan, error)
    PlanUninstall(ctx context.Context, pkg string) (*Plan, error)
    ValidatePlan(p *Plan) error
}
```

**Implementation**: `DefaultPlanner`

**Installation Plan** (binary release):
1. Download release (5s estimated, reversible)
2. Extract archive (2s, reversible)
3. Write to destination (1s, reversible)

**Installation Plan** (from source):
1. Download source (10s estimated, reversible)
2. Build commands (15s each, **non-reversible**)
3. Write to destination (1s, reversible)

**Upgrade Plan**:
1. Backup current version
2. Execute installation plan

**Uninstall Plan**:
1. Cleanup action (remove files)

**Files**: `planner.go` (193 lines)

---

### 8. `internal/source` - Package Source Registry

**Purpose**: Manage package sources (repositories of package metadata).

**Types**:
```go
type Source struct {
    ID          string // Unique identifier
    URI         string // HTTP/file URL
    Name        string
    Description string
    LastUpdated int64  // Unix timestamp
    Enabled     bool
    SourceType  string // "http", "local", "github"
}

type Registry interface {
    AddSource(ctx context.Context, src *Source) error
    RemoveSource(ctx context.Context, idOrURI string) error
    ListSources(ctx context.Context) ([]*Source, error)
    GetSource(ctx context.Context, id string) (*Source, error)
    UpdateSource(ctx context.Context, src *Source) error
}

type Fetcher interface {
    Fetch(ctx context.Context, src *Source) (interface{}, error)
    FetchPackage(ctx context.Context, src *Source, pkgName string) (interface{}, error)
}
```

**Implementation**: `JSONRegistry`
- Storage: `~/.gpkg/sources.json`
- Persists on every add/remove/update
- Lookup by ID or URI
- Tracks `LastUpdated` timestamp

**Files**: `source.go` (41 lines), `registry.go` (136 lines)

---

## Data Flow

### Install Operation (Binary Release)

```
User Command: gpkg install owner/repo --from-release
              ↓
┌─────────────────────────────────────────────────────────────┐
│                     root.go (Cobra)                          │
│  Parse flags: --from-release, --prefix, --offline, etc.     │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│                  install.go (installCmd.RunE)                │
│  1. Load config (precedence: CLI > ENV > user > system)     │
│  2. Parse manifest (local file or resolve from identifier)  │
│  3. Create dependencies:                                     │
│     • HTTPDownloader (timeout, offline mode)                │
│     • AtomicInstaller (backup directory)                    │
│     • Installer (downloader + atomic + prefix)              │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│          installer.go (InstallFromRelease)                   │
│                                                              │
│  1. Download Phase:                                          │
│     HTTPDownloader.DownloadWithChecksum()                   │
│     • HTTP GET request → temp file                          │
│     • Validate checksum (SHA256/SHA1/MD5)                   │
│     • Atomic rename: temp → destination                     │
│                                                              │
│  2. Install Phase:                                           │
│     AtomicInstaller.Install()                               │
│     • Check if binary exists                                │
│     • Backup existing → ~/.gpkg/backups/{timestamp}         │
│     • Copy new binary to ~/.gpkg/bin/{name}                 │
│     • Set permissions (chmod 0755)                          │
│     • Atomic rename to final location                       │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│                  pkgdb/sqlite.go (AddPackage)                │
│  1. Create PackageRecord from manifest                       │
│  2. Serialize JSON fields (checksums, dependencies)          │
│  3. INSERT INTO packages                                     │
│  4. INSERT INTO files (installed file paths)                 │
│  5. INSERT INTO versions (version history)                   │
└─────────────────────────────────────────────────────────────┘
              ↓
       Installation Complete
```

### Install Operation (From Source)

```
User Command: gpkg install owner/repo --from-source
              ↓
┌─────────────────────────────────────────────────────────────┐
│          installer.go (InstallFromSource)                    │
│                                                              │
│  1. Clone/Download Phase:                                    │
│     • git clone {source} → ~/.gpkg/src/{name}               │
│     • OR download tarball → extract                         │
│                                                              │
│  2. Build Phase:                                             │
│     • cd ~/.gpkg/src/{name}                                 │
│     • Execute build commands in sequence                    │
│     • Set environment variables from manifest               │
│                                                              │
│  3. Install Phase:                                           │
│     • Find built executable                                 │
│     • Copy to ~/.gpkg/bin/{name}                            │
│     • Set permissions (chmod 0755)                          │
└─────────────────────────────────────────────────────────────┘
              ↓
       (same database operations as above)
```

### Upgrade Operation

```
User Command: gpkg upgrade {package-name}
              ↓
┌─────────────────────────────────────────────────────────────┐
│              package.go (upgradeCmd.RunE)                    │
│  1. Query database for current version                       │
│  2. Resolve latest version from source                       │
│  3. Compare versions (skip if already latest)                │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│          planner.go (PlanUpgrade)                            │
│  1. Create backup action                                     │
│  2. Create download action (new version)                     │
│  3. Create install action                                    │
│  4. Calculate total size and time                            │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│          installer.go (InstallFromRelease/Source)            │
│  (same as install operation)                                 │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│          pkgdb/sqlite.go (UpdatePackage)                     │
│  1. INSERT INTO versions (previous version)                  │
│  2. UPDATE packages SET version, updated_at, checksums       │
│  3. UPDATE files (new file list)                             │
└─────────────────────────────────────────────────────────────┘
```

### Uninstall Operation

```
User Command: gpkg uninstall {package-name}
              ↓
┌─────────────────────────────────────────────────────────────┐
│            package.go (uninstallCmd.RunE)                    │
│  1. Query database for package metadata                      │
│  2. Get list of installed files                              │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│          planner.go (PlanUninstall)                          │
│  1. Create cleanup action (delete files)                     │
│  2. Calculate total size to free                             │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│              Filesystem Operations                           │
│  1. Delete ~/.gpkg/bin/{name}                               │
│  2. Delete ~/.gpkg/src/{name} (if from source)              │
│  3. Delete associated files                                  │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│          pkgdb/sqlite.go (DeletePackage)                     │
│  1. DELETE FROM packages WHERE name = ?                      │
│  2. CASCADE DELETE from versions and files tables            │
└─────────────────────────────────────────────────────────────┘
```

### Rollback Operation

```
User Command: gpkg rollback {package-name} --to-version {version}
              ↓
┌─────────────────────────────────────────────────────────────┐
│              list.go (rollbackCmd.RunE)                      │
│  1. Query database for package                               │
│  2. Get version history (versions table)                     │
│  3. Find requested version                                   │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│          Backup Restoration                                  │
│  1. Find backup in ~/.gpkg/backups/{timestamp}              │
│  2. Verify backup integrity (checksums)                      │
│  3. Restore files from backup                                │
└─────────────────────────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────────────────────────┐
│          pkgdb/sqlite.go (UpdatePackage)                     │
│  1. UPDATE packages SET version, updated_at                  │
│  2. UPDATE checksums and files                               │
└─────────────────────────────────────────────────────────────┘
```

---

## Key Components

### Component Interaction Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Layer (cmd/)                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ install  │  │ upgrade  │  │ uninstall│  │ rollback │    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘    │
└───────┼────────────┼─────────────┼──────────────┼──────────┘
        │            │             │              │
        └────────────┴─────────────┴──────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│              Business Logic Layer (internal/)                │
│                                                              │
│  Config                    Manifest                          │
│  ├─ Loader (interface)     ├─ Parser (interface)            │
│  └─ YAMLLoader             └─ YAMLParser                    │
│                                                              │
│  Source Registry           Resolver                          │
│  ├─ Registry (interface)   ├─ PackageResolver (interface)   │
│  └─ JSONRegistry           ├─ GitHubResolver                │
│                            └─ DependencyResolver            │
│                                                              │
│  Planner                   Download                          │
│  ├─ Planner (interface)    ├─ Downloader (interface)        │
│  └─ DefaultPlanner         ├─ HTTPDownloader                │
│                            ├─ AtomicInstaller (interface)   │
│                            └─ AtomicInstallerImpl            │
│                                                              │
│  Installer                 PkgDB                             │
│  └─ Installer              ├─ Manager (interface)           │
│     ├─ downloader          └─ SQLiteManager                 │
│     ├─ atomicOp            ├─ PackageRecord                 │
│     └─ prefix              ├─ VersionRecord                 │
│                            └─ FileRecord                    │
└─────────────────────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│              Data Layer (SQLite + Filesystem)                │
│  ~/.gpkg/                                                    │
│  ├── pkgdb.sqlite         (package database)                │
│  ├── bin/                 (installed binaries)              │
│  ├── src/                 (source builds)                   │
│  ├── cache/               (download cache)                  │
│  ├── backups/             (version backups)                 │
│  ├── sources.json         (source registry)                 │
│  └── config.yaml          (user configuration)              │
└─────────────────────────────────────────────────────────────┘
```

### Interface Dependency Graph

```
Installer
   ├── requires: Downloader interface
   │   └── implemented by: HTTPDownloader
   │
   └── requires: AtomicInstaller interface
       └── implemented by: AtomicInstallerImpl

Config
   └── requires: Loader interface
       └── implemented by: YAMLLoader

Manifest
   └── requires: Parser interface
       └── implemented by: YAMLParser

PkgDB
   └── requires: Manager interface
       └── implemented by: SQLiteManager

Resolver
   └── requires: PackageResolver interface
       └── implemented by: GitHubResolver, DependencyResolver

Planner
   └── requires: Planner interface
       └── implemented by: DefaultPlanner

Source
   └── requires: Registry interface
       └── implemented by: JSONRegistry
```

---

## Database Schema

### Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         packages                             │
│  ─────────────────────────────────────────────────────────  │
│  id (PK)               INTEGER                               │
│  name (UNIQUE)         TEXT                                  │
│  version               TEXT                                  │
│  installed_at          TIMESTAMP                             │
│  updated_at            TIMESTAMP                             │
│  source                TEXT ("release" | "source")           │
│  prefix                TEXT                                  │
│  author                TEXT                                  │
│  url                   TEXT                                  │
│  license               TEXT                                  │
│  checksums             TEXT (JSON: {algo: hash})             │
│  dependencies          TEXT (JSON: [pkg1, pkg2])             │
│  build_metadata        TEXT (JSON: {key: value})             │
└────────────┬────────────────────────────────────────────────┘
             │
             │ 1:N
             │
    ┌────────┴─────────┐
    │                  │
    ↓                  ↓
┌──────────────┐  ┌──────────────┐
│  versions    │  │    files     │
│  ──────────  │  │  ──────────  │
│  id (PK)     │  │  id (PK)     │
│  package_id  │  │  package_id  │
│  version     │  │  file_path   │
│  installed_at│  │  checksum    │
│  checksums   │  │  file_size   │
│  files       │  └──────────────┘
└──────────────┘
    (FK)              (FK)
     ↓                 ↓
  CASCADE           CASCADE
  DELETE            DELETE
```

### Table Purposes

**`packages` table** (current installed packages):
- Stores metadata for currently installed packages
- One row per package (name is UNIQUE)
- JSON fields for complex data (checksums, dependencies, metadata)
- Primary source of truth for "what is installed"

**`versions` table** (version history):
- Tracks all installed versions of a package
- Enables rollback to previous versions
- Stores checksums and file lists for each version
- Foreign key to packages (CASCADE DELETE)

**`files` table** (file tracking):
- Tracks every file installed by a package
- Enables complete uninstallation
- Stores checksums for file integrity
- Foreign key to packages (CASCADE DELETE)

### Indexes

```sql
CREATE INDEX idx_packages_name ON packages(name);
-- Fast lookup by package name (most common query)

CREATE INDEX idx_versions_package_id ON versions(package_id);
-- Fast version history retrieval

CREATE INDEX idx_files_package_id ON files(package_id);
-- Fast file list retrieval for uninstall
```

---

## Design Patterns

### 1. Repository Pattern

**Where**: `internal/pkgdb/sqlite.go`

The `SQLiteManager` implements the `Manager` interface, encapsulating all database operations. Consumers interact with the interface, not the concrete SQLite implementation.

```go
type Manager interface {
    AddPackage(p *PackageRecord) (int64, error)
    GetPackage(name string) (*PackageRecord, error)
    UpdatePackage(p *PackageRecord) error
    DeletePackage(name string) error
    ListPackages() ([]*PackageRecord, error)
    // ... more methods
}

// Usage in commands:
db, _ := pkgdb.NewSQLiteManager(dbPath)
pkg, _ := db.GetPackage("example-pkg")  // Abstract database access
```

**Benefits**:
- Hides SQL complexity from consumers
- Easy to swap implementations (PostgreSQL, in-memory, etc.)
- Simplifies testing with mock implementations

---

### 2. Factory Pattern

**Where**: Throughout `internal/` packages

Each package provides factory functions that create fully configured instances:

```go
// internal/download/http_downloader.go
func NewHTTPDownloader(timeout time.Duration, offline bool) *HTTPDownloader

// internal/package/installer.go
func NewInstaller(downloader Downloader, atomicOp AtomicInstaller, prefix string) *Installer

// internal/planner/planner.go
func NewDefaultPlanner(db Manager) *DefaultPlanner

// internal/config/loader.go
func NewYAMLLoader() *YAMLLoader

// internal/manifest/parser.go
func NewYAMLParser() *YAMLParser
```

**Benefits**:
- Encapsulates construction logic
- Ensures proper initialization
- Simplifies object creation

---

### 3. Dependency Injection

**Where**: `cmd/gpkg/cmd/install.go`, command handlers

Dependencies are injected via constructors, not created internally:

```go
// install.go
downloader := download.NewHTTPDownloader(timeout, offline)
atomicInstaller := download.NewAtomicInstaller(backupDir)
installer := pkgmod.NewInstaller(downloader, atomicInstaller, prefix)
```

**Benefits**:
- Loose coupling between components
- Easy to test with mocks/stubs
- Flexible configuration

---

### 4. Strategy Pattern

**Where**: `internal/planner/planner.go`

The `Planner` interface defines different strategies for different operations:

```go
type Planner interface {
    PlanInstallation(ctx, manifest, fromRelease) (*Plan, error)  // Install strategy
    PlanUpgrade(ctx, current, target) (*Plan, error)             // Upgrade strategy
    PlanUninstall(ctx, pkg) (*Plan, error)                       // Uninstall strategy
}
```

Each method implements a different planning strategy with different actions and estimations.

---

### 5. Decorator Pattern

**Where**: `internal/download/http_downloader.go`

`DownloadWithChecksum()` decorates `Download()` with checksum validation:

```go
func (h *HTTPDownloader) DownloadWithChecksum(ctx, url, dest, checksum, algo string) error {
    // Base operation
    if err := h.Download(ctx, url, dest); err != nil {
        return err
    }
    
    // Decoration: checksum validation
    valid, err := h.ValidateChecksum(dest, checksum, algo)
    if !valid {
        return fmt.Errorf("checksum mismatch")
    }
    
    return nil
}
```

---

### 6. Interface-Based Design (Dependency Inversion Principle)

**Where**: Throughout the codebase

All major components depend on interfaces, not concrete implementations:

- **`Downloader`** interface → `HTTPDownloader` implementation
- **`AtomicInstaller`** interface → `AtomicInstallerImpl` implementation
- **`Parser`** interface → `YAMLParser` implementation
- **`Loader`** interface → `YAMLLoader` implementation
- **`Manager`** interface → `SQLiteManager` implementation
- **`Planner`** interface → `DefaultPlanner` implementation
- **`Registry`** interface → `JSONRegistry` implementation

**Benefits**:
- Easy to add new implementations
- Testability through mocking
- Loose coupling
- Flexibility to swap implementations

---

### 7. Command Pattern (via Cobra)

**Where**: `cmd/gpkg/cmd/`

Each CLI command is encapsulated as a Cobra command:

```go
var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install a package",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Command execution logic
    },
}
```

**Benefits**:
- Encapsulates request as an object
- Supports undo (rollback command)
- Easy to add new commands
- Standardized CLI interface

---

### 8. Atomic Operations Pattern

**Where**: `internal/download/http_downloader.go`, `AtomicInstallerImpl`

All file operations use the "write to temp, then rename" pattern:

```go
// Download to temp file
tempFile := filepath.Join(os.TempDir(), "gpkg-download-*")
// ... download to tempFile ...

// Atomic rename
os.Rename(tempFile, destPath)  // Atomic on same filesystem
```

**Benefits**:
- Prevents partial writes
- Ensures consistency (all-or-nothing)
- Safe concurrent access

---

## Dependencies

From `go.mod`:

| Package | Version | Purpose | Where Used |
|---------|---------|---------|------------|
| `github.com/mattn/go-sqlite3` | v1.14.18 | SQLite3 driver for Go (CGO-based) | `internal/pkgdb/sqlite.go` - Database operations |
| `github.com/spf13/cobra` | v1.7.0 | CLI framework (commands, subcommands, flags, help) | `cmd/gpkg/cmd/` - All CLI commands |
| `gopkg.in/yaml.v3` | v3.0.1 | YAML parsing and serialization | `internal/config/loader.go`, `internal/manifest/parser.go` |
| `github.com/inconshreveable/mousetrap` | v1.1.0 | (indirect) Windows Cobra helper | Windows support |
| `github.com/spf13/pflag` | v1.0.5 | (indirect) POSIX-style flags for Cobra | Flag parsing |

### Dependency Rationale

**Why go-sqlite3?**
- Embedded database (no external server)
- Strong Go support
- ACID transactions
- Sufficient for package metadata storage

**Why Cobra?**
- De facto standard for Go CLIs
- Automatic help generation
- Subcommand support
- Flag parsing
- Aliases and suggestions

**Why yaml.v3?**
- Manifest format is YAML
- Pure Go implementation (v3)
- Supports anchors, aliases, and complex structures

---

## Configuration System

### Configuration Sources

gpkg loads configuration from multiple sources with the following **precedence** (highest to lowest):

1. **CLI Flags** (highest priority)
   ```bash
   gpkg install example --prefix=/custom/path --offline
   ```

2. **Environment Variables**
   ```bash
   export GPKG_PREFIX=/custom/path
   export GPKG_CACHE_DIR=/var/cache/gpkg
   export GPKG_LOG_LEVEL=debug
   ```

3. **User Configuration Files**
   - `~/.gpkg/config.yaml`
   - `~/.config/gpkg/config.yaml`
   - `~/.gpkg.yaml`

4. **System Configuration Files**
   - `/etc/gpkg/config.yaml`
   - `/etc/gpkg.yaml`

5. **Hardcoded Defaults** (lowest priority)
   ```yaml
   prefix: ~/.gpkg
   cache_dir: ~/.gpkg/cache
   sources_file: ~/.gpkg/sources.json
   log_level: info
   color: true
   strict_checksum: true
   network_timeout: 30
   ```

### Configuration File Format

```yaml
# ~/.gpkg/config.yaml
prefix: /opt/gpkg
cache_dir: /var/cache/gpkg
sources_file: /etc/gpkg/sources.json
log_level: debug
color: false
strict_checksum: true
network_timeout: 60
```

### Environment Variables

| Variable | Type | Description | Default |
|----------|------|-------------|---------|
| `GPKG_PREFIX` | string | Installation prefix | `~/.gpkg` |
| `GPKG_CACHE_DIR` | string | Download cache directory | `~/.gpkg/cache` |
| `GPKG_LOG_LEVEL` | string | Logging level (debug, info, warn, error) | `info` |
| `GPKG_COLOR` | bool | Enable colored output | `true` |
| `GPKG_OFFLINE` | bool | Offline mode (no network) | `false` |

---

## Architectural Strengths

### 1. Separation of Concerns

Each package has a single, well-defined responsibility:
- **CLI layer** (`cmd/`) handles user interaction
- **Business logic** (`internal/`) handles core operations
- **Data layer** (`pkgdb/`) handles persistence

This makes the codebase easy to navigate and understand.

---

### 2. Testability

Interface-based design and dependency injection make testing straightforward:

```go
// Mock downloader for testing
type MockDownloader struct {
    DownloadFunc func(ctx, url, dest) error
}

func (m *MockDownloader) Download(ctx, url, dest string) error {
    return m.DownloadFunc(ctx, url, dest)
}

// Test installer without network calls
mockDownloader := &MockDownloader{
    DownloadFunc: func(ctx, url, dest string) error {
        // Simulate download
        return nil
    },
}
installer := NewInstaller(mockDownloader, atomicOp, prefix)
```

---

### 3. Atomic Operations

All file operations are atomic:
- Downloads write to temp files, then rename
- Installations create backups before overwriting
- Database transactions ensure consistency

This prevents corruption from crashes or interruptions.

---

### 4. Extensibility

The interface-based design makes it easy to add new features:
- **New manifest formats**: Implement `Parser` interface (JSON, TOML)
- **New download sources**: Implement `Downloader` interface (S3, FTP)
- **New databases**: Implement `Manager` interface (PostgreSQL, MySQL)
- **New package sources**: Implement `Registry` interface (HTTP API)

---

### 5. Configuration Flexibility

Multi-source configuration with precedence allows:
- **System-wide defaults** in `/etc/gpkg/config.yaml`
- **User overrides** in `~/.gpkg/config.yaml`
- **Environment customization** with `GPKG_*` variables
- **Per-command overrides** with CLI flags

---

### 6. Database Normalization

The 3-table schema properly normalizes data:
- **`packages`** table avoids duplication (one row per package)
- **`versions`** table enables version history without duplicating package metadata
- **`files`** table tracks files without duplicating package data
- **Cascading deletes** maintain referential integrity

---

### 7. Error Handling

Context-aware operations support:
- **Cancellation**: `ctx.Done()` checks
- **Timeouts**: Network operations respect timeouts
- **Descriptive errors**: Errors include context (file paths, URLs, etc.)

---

## Areas for Future Enhancement

1. **Dependency Resolution**
   - Current resolver is a stub
   - Need full dependency graph analysis
   - Conflict detection and resolution

2. **Parallel Operations**
   - Parallel downloads (skeleton exists)
   - Concurrent package installations
   - Batch operations

3. **Database Migrations**
   - No migration system currently
   - Schema versioning needed
   - Upgrade path management

4. **Testing Coverage**
   - Sparse unit tests
   - Missing integration tests
   - No end-to-end test suite

5. **Package Source Fetching**
   - `Fetcher` interface not implemented
   - Need HTTP source fetching
   - Package search functionality

6. **Build System**
   - Limited build source support
   - Need better build isolation
   - Container-based builds

---

## File Reference

| File | Lines | Purpose |
|------|-------|---------|
| `main.go` | 10 | Entry point → calls `cmd.Execute()` |
| `cmd/gpkg/cmd/root.go` | 70 | Cobra root command, global flags |
| `cmd/gpkg/cmd/install.go` | 131 | Install command (release/source) |
| `cmd/gpkg/cmd/package.go` | 191 | upgrade, uninstall, info, update |
| `cmd/gpkg/cmd/list.go` | 224 | rollback, list commands |
| `internal/config/config.go` | 56 | Config struct & Loader interface |
| `internal/config/loader.go` | 124 | YAMLLoader implementation |
| `internal/manifest/manifest.go` | 51 | Manifest types |
| `internal/manifest/parser.go` | 68 | YAMLParser with validation |
| `internal/package/package.go` | 32 | Package types |
| `internal/package/installer.go` | 183 | Installer (release/source) |
| `internal/pkgdb/types.go` | 73 | Database types + Manager interface |
| `internal/pkgdb/sqlite.go` | 292 | SQLiteManager + schema |
| `internal/download/download.go` | 36 | Downloader/AtomicInstaller interfaces |
| `internal/download/http_downloader.go` | 192 | HTTPDownloader + AtomicInstallerImpl |
| `internal/resolver/resolver.go` | 155 | GitHubResolver, DependencyResolver |
| `internal/planner/planner.go` | 193 | DefaultPlanner (install/upgrade/uninstall) |
| `internal/source/source.go` | 41 | Source/Registry/Fetcher interfaces |
| `internal/source/registry.go` | 136 | JSONRegistry with persistence |

**Total**: ~2,594 lines of Go code (excluding tests)

---

## Contributing

When contributing to gpkg, follow these architectural guidelines:

1. **Maintain separation of concerns**: Keep CLI logic in `cmd/`, business logic in `internal/`
2. **Use interfaces**: Define interfaces for new components to enable testing and extensibility
3. **Follow dependency injection**: Pass dependencies via constructors
4. **Preserve atomic operations**: Use temp files and rename for all file operations
5. **Update database schema carefully**: Add migrations for schema changes
6. **Write tests**: Add unit tests for new functionality
7. **Document public APIs**: Add godoc comments for exported types and functions

---

## Additional Resources

- **Building gpkg**: See `BUILD.md` for build instructions
- **CLI Specification**: See `CLI+SPEC.md` for command documentation
- **API Documentation**: See `API.md` for public API reference
- **Contributing Guide**: See `CONTRIBUTING.md` for contribution guidelines

---

*Last Updated: 2025*
