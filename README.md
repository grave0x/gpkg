# gpkg

[![Tests](https://github.com/grave0x/gpkg/workflows/Test/badge.svg)](https://github.com/grave0x/gpkg/actions/workflows/test.yml)
[![Code Quality](https://github.com/grave0x/gpkg/workflows/Code%20Quality/badge.svg)](https://github.com/grave0x/gpkg/actions/workflows/code-quality.yml)
[![Security](https://github.com/grave0x/gpkg/workflows/Security/badge.svg)](https://github.com/grave0x/gpkg/actions/workflows/security.yml)
[![Release](https://github.com/grave0x/gpkg/workflows/Release/badge.svg)](https://github.com/grave0x/gpkg/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/go-1.20%2B-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Wiki](https://img.shields.io/badge/docs-wiki-blue.svg)](https://github.com/grave0x/gpkg/wiki)

**gpkg** is a simple, user-focused package manager that installs either release binaries or builds from source into a configurable prefix. Think "pacman for GitHub releases + source builds" — designed to be secure, script-friendly, and easy to use.

## ✨ Features

- 📦 **Install from GitHub releases** or build from source
- 🔒 **Secure downloads** with checksum verification (SHA256/SHA512)
- 🎯 **Isolated installations** in configurable prefix (default: `~/.gpkg`)
- 🔄 **Package management** - install, upgrade, rollback, uninstall
- 📋 **Source repositories** - add multiple package sources
- 🔍 **Package discovery** - search and view package information
- 🤖 **CI-friendly** - JSON output, non-interactive mode, dry-run support
- 📊 **SQLite database** for package tracking and version history

## 🚀 Quick Start

### Installation

#### From Release (Recommended)

```bash
# Download the latest release for your platform
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg-linux-amd64.tar.gz

# Extract and install
tar -xzf gpkg-linux-amd64.tar.gz
sudo mv gpkg-linux-amd64 /usr/local/bin/gpkg
chmod +x /usr/local/bin/gpkg

# Verify installation
gpkg --version
```

#### From Source

**Requirements:**
- Go 1.20 or higher
- SQLite3 (for CGO)
- Git

```bash
# Clone the repository
git clone https://github.com/grave0x/gpkg.git
cd gpkg

# Build
go build -o bin/gpkg .

# Install (optional)
sudo mv bin/gpkg /usr/local/bin/

# Verify
gpkg --version
```

### Basic Usage

```bash
# Add a package source
gpkg add-source https://packages.example.com/index.json

# Update package indices
gpkg update

# Search for packages
gpkg search tool

# View package information
gpkg info owner/repo

# Install a package from GitHub release
gpkg install owner/repo

# Install from local manifest
gpkg install ./examples/manifest.yaml

# List installed packages
gpkg list --installed

# Upgrade all packages
gpkg upgrade

# Uninstall a package
gpkg uninstall package-name
```

## 📚 Documentation

### Commands

#### Package Installation
```bash
gpkg install <package>              # Install from release
gpkg install <package> --from-source  # Build from source
gpkg install ./manifest.yaml        # Install from local manifest
gpkg install <package> --prefix=/opt  # Custom install location
gpkg install <package> --dry-run    # Preview without installing
```

#### Package Management
```bash
gpkg upgrade [package...]           # Upgrade packages (all if no args)
gpkg uninstall <package>            # Remove installed package
gpkg rollback <package> <version>   # Rollback to previous version
gpkg list --installed               # List installed packages
gpkg list --available               # List available packages
```

#### Source Management
```bash
gpkg add-source <uri>               # Add package source
gpkg remove-source <uri|id>         # Remove source
gpkg list-sources                   # List all sources
gpkg update                         # Refresh package indices
```

#### Package Information
```bash
gpkg info <package>                 # Show package details
gpkg search <term>                  # Search package indices
gpkg validate <manifest>            # Validate manifest file
```

#### Configuration
```bash
gpkg config show                    # Display current config
gpkg config get <key>               # Get config value
gpkg config set <key> <value>       # Set config value
```

### Global Flags

All commands support these global flags:

```bash
-c, --config <path>     # Custom config file
    --json              # JSON output (machine-readable)
-y, --yes               # Non-interactive mode (assume yes)
    --dry-run           # Preview actions without modifying system
    --log-level <level> # Logging level (error, warn, info, debug)
-v, --verbose           # Increase verbosity (stackable: -vv, -vvv)
    --quiet             # Minimal output
    --no-color          # Disable colored output
    --offline           # Disable network operations
```

### Configuration

**Config File Locations** (in order of precedence):
1. CLI flag: `--config /path/to/config.toml`
2. User config: `~/.gpkg/config.toml`
3. System config: `/etc/gpkg/config.toml`

**Default Configuration:**
```toml
install_prefix = "~/.gpkg"
pkgdb_path = "~/.gpkg/pkgdb.sqlite"
sources_dir = "~/.gpkg/sources.d"
require_checksums = true
parallel_downloads = 4
```

**Environment Variables:**
- `GPKG_INSTALL_PREFIX` - Override install directory
- `GPKG_PKGDB_PATH` - Override database location
- `GPKG_CONFIG` - Override config file location
- `LOG_LEVEL` - Set logging level

### Manifest Format

Packages are defined using YAML manifests. See `examples/manifest.yaml` for a complete example:

```yaml
name: cool-tool
version: 1.2.0
description: "A nifty CLI tool"
homepage: "https://github.com/owner/cool-tool"
license: "MIT"

# Release-based installation
release_assets:
  - url: "https://github.com/owner/cool-tool/releases/download/v1.2.0/cool-tool-linux-amd64.tar.gz"
    platform: "linux/amd64"
    checksum:
      sha256: "abc123..."

# Optional: Build from source
build:
  repo: "https://github.com/owner/cool-tool.git"
  build_steps:
    - "make"
    - "make install PREFIX={install_prefix}"
  build_env:
    - "CGO_ENABLED=0"

files_to_install:
  - "bin/cool-tool"
  - "share/man/man1/cool-tool.1"
```

## 🛠️ Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/grave0x/gpkg.git
cd gpkg

# Install dependencies
go mod download

# Build
go build -o bin/gpkg .

# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Build for multiple platforms
./scripts/build.sh 0.1.0
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with race detector
go test -race ./...

# Run specific package tests
go test ./internal/pkgdb

# Run integration tests
go test ./tests/integration
```

### Project Structure

```
gpkg/
├── cmd/gpkg/           # CLI commands (Cobra)
│   └── cmd/            # Command implementations
├── internal/           # Internal packages
│   ├── config/         # Configuration management
│   ├── download/       # Download utilities
│   ├── manifest/       # Manifest parsing
│   ├── package/        # Package installation logic
│   ├── pkgdb/          # SQLite database operations
│   ├── planner/        # Installation planning
│   ├── resolver/       # Dependency resolution
│   └── source/         # Source management
├── tests/              # Integration tests
├── examples/           # Example manifests
├── scripts/            # Build and release scripts
└── .github/workflows/  # CI/CD workflows
```

### Test Coverage

| Package | Coverage |
|---------|----------|
| pkgdb | 85.0% |
| planner | 70.3% |
| source | 66.1% |
| manifest | 62.5% |
| resolver | 52.8% |
| config | 35.8% |
| download | 30.5% |

## 🔒 Security

- **Checksum verification**: All downloads are verified with SHA256/SHA512 checksums
- **Atomic installations**: Packages are staged in temporary directories and moved atomically
- **Isolated prefix**: Packages install to `~/.gpkg` by default, isolated from system
- **File tracking**: All installed files are recorded in SQLite database
- **Safe uninstalls**: Only files within the install prefix can be removed

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Current priority areas:**
- Source manager improvements ([#2](https://github.com/grave0x/gpkg/issues/2))
- Repository viewer enhancements ([#6](https://github.com/grave0x/gpkg/issues/6))
- Test coverage improvements
- Documentation updates

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🗺️ Roadmap

### Sprint 1 (MVP) ✅
- [x] CLI framework with Cobra
- [x] Install from release binaries
- [x] Manifest parsing and validation
- [x] SQLite package database
- [x] Source management (add/remove/list)
- [x] Package information viewer
- [x] Update/upgrade workflows
- [x] Uninstall functionality

### Sprint 2 (In Progress)
- [ ] Enhanced checksum verification
- [ ] Resumable downloads
- [ ] Atomic installation improvements
- [ ] Rollback support
- [ ] Manifest linter
- [ ] Dependency resolution

### Sprint 3 (Future)
- [ ] Source builds with sandboxing
- [ ] Package signatures
- [ ] Background updater
- [ ] GUI/TUI interface
- [ ] Plugin system

## 📞 Support & Documentation

- **📚 Wiki**: [Complete Documentation](https://github.com/grave0x/gpkg/wiki)
  - [Getting Started](https://github.com/grave0x/gpkg/wiki/Getting-Started)
  - [User Guide](https://github.com/grave0x/gpkg/wiki/User-Guide)
  - [Developer Guide](https://github.com/grave0x/gpkg/wiki/Developer-Guide)
  - [API Reference](https://github.com/grave0x/gpkg/wiki/Reference)
- **🐛 Issues**: [GitHub Issues](https://github.com/grave0x/gpkg/issues)
- **🔧 Troubleshooting**: [Common Issues](https://github.com/grave0x/gpkg/wiki/Troubleshooting)
- **🏗️ Workflows**: [CI/CD Documentation](.github/workflows/README.md)

## 🙏 Acknowledgments

- Inspired by pacman, Homebrew, and other package managers
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses [go-sqlite3](https://github.com/mattn/go-sqlite3) for database operations

---

**Made with ❤️ by the gpkg team**
