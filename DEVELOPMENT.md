# Development Guide

## Project Structure

```
gpkg/
├── cmd/gpkg/cmd/          # CLI commands
│   ├── root.go            # Root command and global flags
│   ├── source.go          # Source management commands
│   ├── install.go         # Install command
│   ├── package.go         # Package management commands
│   └── globals.go         # Global state
├── internal/
│   ├── config/            # Configuration management
│   │   ├── config.go      # Config structures
│   │   └── loader.go      # YAML config loader
│   ├── download/          # Download and installation
│   │   ├── download.go    # Interfaces
│   │   └── http_downloader.go # HTTP implementation
│   ├── manifest/          # Manifest parsing
│   │   ├── manifest.go    # Manifest structures
│   │   └── parser.go      # YAML parser
│   ├── package/           # Package management
│   │   ├── package.go     # Package structures
│   │   └── installer.go   # Installation logic
│   └── source/            # Package sources
│       ├── source.go      # Source interfaces
│       └── registry.go    # JSON registry implementation
├── pkg/                   # Public packages (future)
├── main.go               # Entry point
├── go.mod                # Go module definition
└── README.md             # User documentation
```

## Building

### Prerequisites
- Go 1.20 or higher
- Git (for source builds)

### Build Steps

```bash
# Clone the repository
git clone https://github.com/grave0x/gpkg.git
cd gpkg

# Fetch dependencies
go mod download

# Build the binary
go build -o bin/gpkg ./cmd/gpkg

# Run tests
go test ./...

# Build with version info (optional)
go build -ldflags "-X main.version=0.1.0" -o bin/gpkg ./cmd/gpkg
```

### Installation

```bash
# Install to $GOPATH/bin
go install ./cmd/gpkg

# Or copy binary manually
cp bin/gpkg /usr/local/bin/
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/manifest
go test -v ./internal/config
go test -v ./internal/source

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Code Organization

### Interfaces vs Implementations

The project uses interface-based design for flexibility:

- **Manifest Parser**: `Parser` interface with `YAMLParser` implementation
- **Config Loader**: `Loader` interface with `YAMLLoader` implementation
- **Source Registry**: `Registry` interface with `JSONRegistry` implementation
- **Downloader**: `Downloader` interface with `HTTPDownloader` implementation
- **Atomic Installer**: `AtomicInstaller` interface with `AtomicInstallerImpl` implementation

This allows easy substitution for testing and alternative implementations.

### Error Handling

All errors are wrapped with context using `fmt.Errorf(...%w...)` pattern:

```go
err := someFunction()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

## Adding New Commands

1. Create a new file in `cmd/gpkg/cmd/`
2. Define `var <name>Cmd = &cobra.Command{...}`
3. Add subcommand to root in `init()`:
   ```go
   rootCmd.AddCommand(<name>Cmd)
   ```

Example:
```go
var newCmd = &cobra.Command{
    Use:   "new-cmd <arg>",
    Short: "Description",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

## Configuration Precedence

Configuration is loaded in this order (later overrides earlier):

1. System config: `/etc/gpkg/config.yaml` or `/etc/gpkg.yaml`
2. User config: `~/.gpkg/config.yaml`, `~/.gpkg.yaml`, or `~/.config/gpkg/config.yaml`
3. Environment variables: `GPKG_PREFIX`, `GPKG_CACHE_DIR`, `GPKG_LOG_LEVEL`
4. CLI flags: `--prefix`, `--config`, etc.

## Manifest Format

Manifests describe how to install a package. See `CLI+SPEC.md` for detailed examples.

### Example: Release Binary
```yaml
package:
  name: my-tool
  version: 1.0.0
  author: Author Name
  url: https://example.com
  license: MIT
install:
  type: binary
  source: https://github.com/owner/repo/releases/download/v1.0.0/my-tool-linux-x64
  checksum:
    sha256: abc123...
  executable: my-tool
```

### Example: Source Build
```yaml
package:
  name: my-app
  version: 1.0.0
build_source:
  type: git
  source: https://github.com/owner/repo.git
  tag: v1.0.0
  commands:
    - make
    - make install PREFIX=$GPKG_PREFIX
  env:
    GPKG_PREFIX: /path/to/install
```

## Testing Patterns

### Unit Tests

Created test files use `*_test.go` naming convention:

```go
package manifest_test

import "testing"

func TestSomething(t *testing.T) {
    // Test implementation
}
```

### Test Utilities

- Use `t.TempDir()` for temporary file operations
- Use `filepath.Join()` for cross-platform paths
- Check for errors with `if err != nil { t.Fatalf(...) }`

## Logging

Currently using standard library logging. Future enhancement: structured logging with levels (error, warn, info, debug) controlled by `--log-level` flag.

## Performance Considerations

1. **Parallel Downloads**: Not yet implemented; could be added for multiple packages
2. **Caching**: Source metadata is cached; individual package cache in progress
3. **Atomic Operations**: Installations use temp files to prevent partial states

## Future Enhancements

- [ ] GitHub package index support
- [ ] GPG signature verification
- [ ] Parallel package downloads
- [ ] Post-install hooks
- [ ] Package dependency resolution
- [ ] Upgrade notifications
- [ ] Shell completions (bash, zsh, fish)
