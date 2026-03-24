# API Reference

This document describes the internal APIs and interfaces used by gpkg.

## Package: internal/config

### Config Structure

```go
type Config struct {
    Prefix string         // Installation prefix (default: ~/.gpkg)
    CacheDir string       // Cache directory
    SourcesFile string    // Sources registry file
    LogLevel string       // Logging level
    Color bool            // Enable color output
    StrictChecksum bool   // Checksum verification strict mode
    NetworkTimeout int    // Network timeout in seconds
}
```

### Loader Interface

```go
type Loader interface {
    Load() (*Config, error)                    // Load with precedence
    LoadFrom(path string) (*Config, error)     // Load from specific file
    MergeDefaults(cfg *Config) *Config         // Apply defaults
}
```

### Functions

- `DefaultConfig() *Config` - Returns config with sensible defaults
- `NewYAMLLoader(path string) *Loader` - Creates YAML-based loader

## Package: internal/manifest

### Manifest Structure

```go
type Manifest struct {
    Package PackageSpec           // Package metadata
    Install *InstallSpec          // Release installation spec
    BuildSource *BuildSourceSpec  // Source build spec
    Dependencies []string         // Package dependencies
}

type PackageSpec struct {
    Name string    // Package name
    Version string // Package version
    Author string  // Author name
    URL string     // Package URL
    License string // License identifier
}

type InstallSpec struct {
    Type string               // "release", "binary", "archive"
    Source string             // Download URL
    Pattern string            // File pattern
    Checksum map[string]string // Algorithm -> hash
    ExtractPath string        // Path within archive
    Executable string         // Executable name
    PostInstall string        // Post-install command
}

type BuildSourceSpec struct {
    Type string                // "git", "tarball"
    Source string              // Repository/source URL
    Tag string                 // Git tag
    Branch string              // Git branch
    Commands []string          // Build commands
    Env map[string]string     // Environment variables
}
```

### Parser Interface

```go
type Parser interface {
    Parse(path string) (*Manifest, error)
    ParseBytes(data []byte) (*Manifest, error)
    Validate(m *Manifest) error
}
```

### Functions

- `NewYAMLParser() *Parser` - Creates YAML manifest parser

## Package: internal/package

### Package Structure

```go
type Package struct {
    Name string
    Description string
    Version string
    Author string
    URL string
    License string
    Metadata map[string]string
}

type InstalledPackage struct {
    Package *Package
    InstalledAt time.Time
    Source string             // "release" or "source"
    Prefix string             // Installation prefix
    Checksums map[string]string
    Dependencies []string
}

type VersionInfo struct {
    Latest string
    Current string
    Available []string
    ReleaseDate time.Time
}
```

### Installer

```go
type Installer struct {
    // Private fields
}

// Constructor
NewInstaller(dl Downloader, atomic AtomicInstaller, prefix string) *Installer

// Methods
InstallFromRelease(ctx context.Context, mf *Manifest) (*InstalledPackage, error)
InstallFromSource(ctx context.Context, mf *Manifest) (*InstalledPackage, error)
```

## Package: internal/source

### Source Structure

```go
type Source struct {
    ID string           // Unique identifier
    URI string          // Source URL
    Name string         // Display name
    Description string  // Description
    LastUpdated int64   // Unix timestamp
    Enabled bool        // Is enabled
    SourceType string   // "http", "local", "github"
}
```

### Registry Interface

```go
type Registry interface {
    AddSource(ctx context.Context, src *Source) error
    RemoveSource(ctx context.Context, idOrURI string) error
    ListSources(ctx context.Context) ([]*Source, error)
    GetSource(ctx context.Context, id string) (*Source, error)
    UpdateSource(ctx context.Context, src *Source) error
}
```

### JSONRegistry

```go
type JSONRegistry struct {
    // Private fields
}

// Constructor
NewJSONRegistry(filePath string) *JSONRegistry

// Methods
Load() error      // Load from disk
Save() error      // Save to disk
```

## Package: internal/download

### Downloader Interface

```go
type Downloader interface {
    Download(ctx context.Context, url, dest string) error
    DownloadWithChecksum(ctx context.Context, url, dest, hash, algo string) error
    ValidateChecksum(filePath, hash, algo string) (bool, error)
}
```

### HTTPDownloader

```go
type HTTPDownloader struct {
    // Private fields
}

// Constructor
NewHTTPDownloader(timeout time.Duration, offline bool) *HTTPDownloader
```

### AtomicInstaller Interface

```go
type AtomicInstaller interface {
    Install(ctx context.Context, src, dest string) error
    Uninstall(ctx context.Context, pkgName string) error
    Rollback(ctx context.Context, installID string) error
}
```

### Supported Checksum Algorithms

- `sha256` (recommended)
- `sha1`
- `md5`

## Global Configuration

### Environment Variables

- `GPKG_PREFIX` - Override installation prefix
- `GPKG_CACHE_DIR` - Override cache directory
- `GPKG_LOG_LEVEL` - Override log level

### Default Locations

- Config: `~/.gpkg/config.yaml`, `~/.config/gpkg/config.yaml`
- Sources: `~/.gpkg/sources.json`
- Install prefix: `~/.gpkg`
- Cache: `~/.gpkg/cache`

## CLI Global Flags

- `-c, --config <path>` - Custom config file
- `--json` - JSON output
- `-y, --yes` - Non-interactive mode
- `--dry-run` - Simulate without performing actions
- `--log-level <level>` - Set log level
- `-v, --verbose` - Increase verbosity (stackable)
- `--quiet` - Minimal output
- `--no-color` - Disable colors
- `--offline` - Disallow network access
- `-V, --version` - Show version
- `-h, --help` - Show help

## Context Usage

All async operations accept `context.Context` for cancellation and timeout support:

```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := installer.InstallFromRelease(ctx, manifest)
```

## Error Handling

All functions return errors wrapped with context:

```go
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

Errors can be unwrapped using `errors.Unwrap()` for checking specific causes.
