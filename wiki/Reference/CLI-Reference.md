# gpkg CLI Reference

Complete command-line reference for gpkg, the user-focused package manager.

---

## Table of Contents

1. [Overview](#overview)
2. [Global Options](#global-options)
3. [Commands](#commands)
   - [Package Installation & Management](#package-installation--management)
   - [Source Management](#source-management)
   - [Package Discovery](#package-discovery)
   - [Configuration & Utilities](#configuration--utilities)
4. [Environment Variables](#environment-variables)
5. [Configuration Precedence](#configuration-precedence)
6. [Exit Codes](#exit-codes)
7. [JSON Output](#json-output)

---

## Overview

### CLI Design Philosophy

gpkg follows a git/pacman-like subcommand design with these principles:

- **First-class subcommands**: Each command is a discrete operation with its own help and flags
- **POSIX-style flags**: Both long (`--flag`) and short (`-f`) forms where practical
- **Sensible defaults**: Optimized for interactive use with explicit flags for automation/CI
- **Machine-readable output**: `--json` flag for structured output on supported commands
- **Safety first**: Checksums required by default, confirmations for destructive operations
- **Dual installation modes**: Install from release binaries or build from source

### Usage Pattern

```
gpkg [global-options] <command> [command-options] [arguments]
```

All commands support `--help` for detailed usage information.

---

## Global Options

These flags apply to **all** gpkg commands:

### Help and Version

| Flag | Short | Description |
|------|-------|-------------|
| `--help` | `-h` | Display help for gpkg or any command |
| `--version` | `-V` | Print gpkg version and exit |

### Configuration

| Flag | Short | Description |
|------|-------|-------------|
| `--config <path>` | `-c` | Path to config file (overrides default locations) |

### Output Control

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | | Machine-readable JSON output (where supported) |
| `--verbose` | `-v` | Increase verbosity (stackable: `-vvv`) |
| `--quiet` | | Minimal output (conflicts with `--verbose`) |
| `--no-color` | | Disable colorized output for CI/scripting |
| `--log-level <level>` | | Override logging level: `error`, `warn`, `info`, `debug` (default: `info`) |

### Execution Mode

| Flag | Short | Description |
|------|-------|-------------|
| `--yes` | `-y` | Assume "yes" for all prompts (non-interactive mode) |
| `--dry-run` | | Plan actions without modifying filesystem/pkgdb (also `--simulate`) |
| `--offline` | | Disallow network fetches; fail when remote data required |

**Note**: Global flags can be placed before or after the command name.

---

## Commands

### Package Installation & Management

#### install

Install a package from a release binary or build from source.

**Synopsis**

```
gpkg install <pkg|manifest> [--from-release|--from-source] [flags]
```

**Description**

Installs a package into the configured prefix (default: `~/.gpkg`). Can install from:
- Package name (requires configured sources)
- GitHub repository (`owner/repo` format)
- Local manifest file path

By default, prefers release binaries when available. Use `--from-source` to force building from source.

**Arguments**

- `<pkg|manifest>`: Package identifier, GitHub repo (owner/repo), or path to manifest file

**Flags**

| Flag | Description |
|------|-------------|
| `--from-release` | Install from release binary (default when releases available) |
| `--from-source` | Install by building from source |
| `--prefix <path>` | Installation prefix (overrides config) |

**Examples**

```bash
# Install from GitHub repository
gpkg install owner/cool-tool

# Force source build
gpkg install owner/cool-tool --from-source

# Install from local manifest
gpkg install ./examples/cool-tool.yaml

# Install to custom location
gpkg install my-package --prefix=/opt/gpkg

# Plan installation without executing
gpkg install owner/repo --dry-run

# Non-interactive install
gpkg install owner/repo --yes
```

**Exit Codes**

- `0`: Package installed successfully
- `3`: Network error (remote unavailable)
- `4`: Checksum verification failed
- `5`: Build/install failed
- `6`: Manifest validation error
- `7`: Package not found
- `8`: Package database error

**Related Commands**

- `gpkg uninstall` - Remove installed package
- `gpkg upgrade` - Upgrade to newer version
- `gpkg info` - View package information before installing

**JSON Output**

When `--json` is specified:

```json
{
  "package": "cool-tool",
  "version": "1.2.0",
  "status": "installed",
  "files": [
    "bin/cool-tool",
    "share/man/cool-tool.1"
  ],
  "prefix": "/home/user/.gpkg"
}
```

---

#### uninstall

Remove an installed package from the system.

**Synopsis**

```
gpkg uninstall <pkg> [flags]
```

**Description**

Uninstalls a package by removing all tracked files and database records. Files are only removed from within the configured install prefix for safety. Prompts for confirmation unless `--yes` is specified.

**Arguments**

- `<pkg>`: Name of installed package to remove

**Flags**

| Flag | Description |
|------|-------------|
| `--prefix <path>` | Installation prefix (must match install location) |

**Examples**

```bash
# Uninstall a package (prompts for confirmation)
gpkg uninstall my-package

# Non-interactive uninstall
gpkg uninstall my-package --yes

# Uninstall from custom prefix
gpkg uninstall my-package --prefix=/opt/gpkg

# Plan uninstall without executing
gpkg uninstall my-package --dry-run
```

**Exit Codes**

- `0`: Package uninstalled successfully
- `1`: General failure
- `7`: Package not found
- `8`: Package database error

**Related Commands**

- `gpkg install` - Install packages
- `gpkg list` - View installed packages

---

#### upgrade

Upgrade installed packages to their latest available versions.

**Synopsis**

```
gpkg upgrade [pkg...] [flags]
```

**Description**

Upgrades packages to the latest versions available in configured sources:
- Without arguments: upgrades **all** installed packages with newer versions available
- With arguments: upgrades only the specified packages

**Arguments**

- `[pkg...]`: Optional list of packages to upgrade (defaults to all)

**Flags**

No command-specific flags.

**Examples**

```bash
# Upgrade all installed packages
gpkg upgrade

# Upgrade specific package
gpkg upgrade my-tool

# Upgrade multiple packages
gpkg upgrade pkg1 pkg2 pkg3

# See what would be upgraded (dry-run)
gpkg upgrade --dry-run

# Non-interactive upgrade
gpkg upgrade --yes

# Verbose upgrade output
gpkg upgrade -vv
```

**Exit Codes**

- `0`: Upgrade completed successfully
- `3`: Network error
- `4`: Checksum verification failed
- `5`: Upgrade failed
- `8`: Package database error

**Related Commands**

- `gpkg update` - Refresh package source metadata first
- `gpkg list` - View installed packages and versions

**JSON Output**

When `--json` is specified:

```json
{
  "upgraded": [
    {
      "package": "tool1",
      "from_version": "1.0.0",
      "to_version": "1.2.0",
      "status": "success"
    }
  ],
  "skipped": [
    {
      "package": "tool2",
      "version": "2.0.0",
      "reason": "already latest"
    }
  ]
}
```

---

#### rollback

Rollback an installed package to a previous version.

**Synopsis**

```
gpkg rollback <pkg> --to-version <version> [flags]
```

**Description**

Reverts a package to a previously installed version. Requires that the target version's installation data is still available in the package database history.

**Arguments**

- `<pkg>`: Name of installed package to rollback

**Flags**

| Flag | Description |
|------|-------------|
| `--to-version <ver>` | **Required**. Version to rollback to |
| `--prefix <path>` | Installation prefix |

**Examples**

```bash
# Rollback to specific version
gpkg rollback my-tool --to-version 1.0.0

# Rollback with custom prefix
gpkg rollback my-tool --to-version 1.0.0 --prefix=/opt/gpkg

# Plan rollback (dry-run)
gpkg rollback my-tool --to-version 1.0.0 --dry-run
```

**Exit Codes**

- `0`: Rollback completed successfully
- `1`: General failure
- `7`: Package or version not found
- `8`: Package database error

**Related Commands**

- `gpkg upgrade` - Upgrade to newer versions
- `gpkg info` - View version history

---

### Source Management

#### add-source

Add a new package source.

**Synopsis**

```
gpkg add-source <uri> [flags]
```

**Description**

Registers a package source (index) where packages can be discovered and installed from. Sources can be:
- HTTP/HTTPS JSON index URLs
- Local directories containing manifest files
- GitHub shorthand: `github:owner/repo`

**Arguments**

- `<uri>`: Source URI (URL, file path, or github: shorthand)

**Flags**

Currently no command-specific flags (spec mentions `--name <id>`, `--priority <n>`, `--force` - implementation may vary).

**Examples**

```bash
# Add HTTP source
gpkg add-source https://packages.example.com/index.json

# Add GitHub source
gpkg add-source github:myorg/packages

# Add local directory source
gpkg add-source /opt/local-packages

# Verify what would be added (dry-run)
gpkg add-source https://example.com/index.json --dry-run
```

**Exit Codes**

- `0`: Source added successfully
- `1`: General failure
- `3`: Network error (when fetching remote source)

**Related Commands**

- `gpkg remove-source` - Remove sources
- `gpkg list-sources` - View configured sources
- `gpkg update` - Refresh source metadata

**JSON Output**

When `--json` is specified:

```json
{
  "added": "source-id",
  "uri": "https://packages.example.com/index.json"
}
```

---

#### remove-source

Remove a package source.

**Synopsis**

```
gpkg remove-source <id|uri> [flags]
```

**Description**

Removes a configured package source by its ID or URI. Prompts for confirmation unless `--yes` is specified.

**Arguments**

- `<id|uri>`: Source identifier or URI to remove

**Flags**

No command-specific flags.

**Examples**

```bash
# Remove by source ID
gpkg remove-source source-1

# Remove by URI
gpkg remove-source https://packages.example.com/index.json

# Non-interactive removal
gpkg remove-source source-1 --yes
```

**Exit Codes**

- `0`: Source removed successfully
- `1`: General failure
- `7`: Source not found

**Related Commands**

- `gpkg add-source` - Add sources
- `gpkg list-sources` - View configured sources

---

#### list-sources

List all registered package sources.

**Synopsis**

```
gpkg list-sources [flags]
```

**Description**

Displays all configured package sources with their IDs, URIs, priorities, and last update times.

**Arguments**

None.

**Flags**

No command-specific flags.

**Examples**

```bash
# List all sources
gpkg list-sources

# JSON output
gpkg list-sources --json
```

**Exit Codes**

- `0`: Success

**Related Commands**

- `gpkg add-source` - Add sources
- `gpkg remove-source` - Remove sources
- `gpkg update` - Refresh source metadata

**JSON Output**

When `--json` is specified:

```json
{
  "sources": [
    {
      "id": "default",
      "uri": "https://packages.gpkg.io/index.json",
      "priority": 100,
      "last_updated": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

#### update

Refresh package source metadata from all configured sources.

**Synopsis**

```
gpkg update [flags]
```

**Description**

Fetches the latest package lists and version information from all enabled sources. This updates the local cache but does **not** install or upgrade any packages.

Similar to `apt update` or `pacman -Sy`.

**Arguments**

None.

**Flags**

No command-specific flags (spec mentions `--concurrency N`, `--force` - implementation may vary).

**Examples**

```bash
# Update all sources
gpkg update

# Update with verbose output
gpkg update -v

# JSON output
gpkg update --json
```

**Exit Codes**

- `0`: Update completed successfully
- `3`: Network error

**Related Commands**

- `gpkg upgrade` - Upgrade packages after updating
- `gpkg list-sources` - View configured sources

**JSON Output**

When `--json` is specified:

```json
{
  "updated": [
    {
      "source": "default",
      "packages": 127,
      "status": "success"
    }
  ],
  "failed": []
}
```

---

### Package Discovery

#### search

Search across configured package sources.

**Synopsis**

```
gpkg search <term> [flags]
```

**Description**

Searches package names and descriptions for the specified term across all configured sources.

**Arguments**

- `<term>`: Search term or pattern

**Flags**

| Flag | Description |
|------|-------------|
| `--source <id>` | Search only in specific source |

**Examples**

```bash
# Search all sources
gpkg search golang

# Search specific source
gpkg search --source github golang

# JSON output
gpkg search --json curl
```

**Exit Codes**

- `0`: Search completed (even with zero results)
- `3`: Network error

**Related Commands**

- `gpkg list` - List all available packages
- `gpkg info` - Get detailed package information

**JSON Output**

When `--json` is specified:

```json
{
  "results": [
    {
      "name": "go-tool",
      "version": "1.2.0",
      "description": "A tool written in Go",
      "source": "default"
    }
  ],
  "total": 1
}
```

---

#### info

Display detailed information about a package.

**Synopsis**

```
gpkg info <repo|pkg|manifest> [flags]
```

**Description**

Shows comprehensive package information including:
- Package metadata (name, version, description)
- Latest available version
- Installed version (if applicable)
- Release assets and checksums
- Build steps (for source builds)
- Installation methods
- Dependencies (with `--deps-tree`)

**Arguments**

- `<repo|pkg|manifest>`: GitHub repo (owner/repo), package name, or local manifest path

**Flags**

| Flag | Description |
|------|-------------|
| `--deps-tree` | Show dependency tree |
| `--raw` | Show raw manifest content |
| `--parsed` | Show parsed/structured manifest |

**Examples**

```bash
# Info from GitHub repo
gpkg info owner/cool-tool

# Info with dependency tree
gpkg info my-package --deps-tree

# Show raw manifest
gpkg info ./manifest.yaml --raw

# JSON output
gpkg info owner/repo --json
```

**Exit Codes**

- `0`: Success
- `3`: Network error
- `6`: Manifest validation error
- `7`: Package not found

**Related Commands**

- `gpkg search` - Find packages
- `gpkg install` - Install after reviewing info

**JSON Output**

When `--json` is specified:

```json
{
  "name": "cool-tool",
  "version": "1.2.0",
  "description": "A cool tool",
  "installed": {
    "version": "1.1.0",
    "prefix": "/home/user/.gpkg"
  },
  "releases": [
    {
      "version": "1.2.0",
      "assets": [
        {
          "platform": "linux",
          "arch": "amd64",
          "url": "https://...",
          "checksum": "sha256:..."
        }
      ]
    }
  ],
  "manifest": {
    "name": "cool-tool",
    "version": "1.2.0"
  }
}
```

---

#### list

List installed or available packages.

**Synopsis**

```
gpkg list [--installed|--available] [flags]
```

**Description**

Lists packages with filtering and sorting options:
- `--installed` (default): Shows packages currently installed
- `--available`: Shows packages available from configured sources

**Arguments**

None.

**Flags**

| Flag | Description |
|------|-------------|
| `--installed` | Show installed packages (default) |
| `--available` | Show available packages from sources |
| `--filter <pattern>` | Filter by name pattern |
| `--sort <field>` | Sort by: `name`, `version`, `date` (default: `name`) |

**Examples**

```bash
# List installed packages
gpkg list

# List available packages
gpkg list --available

# Filter packages
gpkg list --filter "tool"

# Sort by version
gpkg list --sort version

# JSON output
gpkg list --json
```

**Exit Codes**

- `0`: Success

**Related Commands**

- `gpkg search` - Search for specific packages
- `gpkg info` - Get detailed package information

**JSON Output**

When `--json` is specified:

```json
{
  "packages": [
    {
      "name": "tool1",
      "version": "1.0.0",
      "installed_at": "2024-01-15T10:30:00Z",
      "prefix": "/home/user/.gpkg"
    }
  ],
  "total": 1
}
```

---

### Configuration & Utilities

#### config

Manage gpkg configuration.

**Synopsis**

```
gpkg config [command]
```

**Description**

Configuration management with three subcommands:
- `get <key>`: Retrieve configuration value
- `set <key> <value>`: Set configuration value
- `show`: Display merged configuration from all sources

**Subcommands**

##### config get

Retrieve a configuration value.

**Synopsis**

```
gpkg config get <key> [flags]
```

**Examples**

```bash
# Get install prefix
gpkg config get install_prefix

# Get with JSON output
gpkg config get install_prefix --json
```

##### config set

Set a configuration value.

**Synopsis**

```
gpkg config set <key> <value> [flags]
```

**Examples**

```bash
# Set install prefix
gpkg config set install_prefix /opt/gpkg

# Set log level
gpkg config set log_level debug
```

##### config show

Display merged configuration.

**Synopsis**

```
gpkg config show [flags]
```

**Examples**

```bash
# Show all configuration
gpkg config show

# Show as JSON
gpkg config show --json
```

**Configuration Keys**

Common configuration keys:

| Key | Description | Default |
|-----|-------------|---------|
| `install_prefix` | Installation directory | `~/.gpkg` |
| `pkgdb_path` | Package database location | `~/.gpkg/pkgdb.sqlite` |
| `sources_dir` | Sources directory | `~/.gpkg/sources.d` |
| `require_checksums` | Require checksums for installs | `true` |
| `log_level` | Logging level | `info` |
| `allow_unverified_source_builds` | Allow builds without provenance | `false` |

**Exit Codes**

- `0`: Success
- `1`: General failure
- `2`: Invalid key or value

**Related Commands**

None.

**JSON Output**

When `--json` is specified on `config show`:

```json
{
  "install_prefix": "/home/user/.gpkg",
  "pkgdb_path": "/home/user/.gpkg/pkgdb.sqlite",
  "sources_dir": "/home/user/.gpkg/sources.d",
  "require_checksums": true,
  "log_level": "info"
}
```

---

#### validate

Validate a manifest file against the schema.

**Synopsis**

```
gpkg validate <manifest-path> [flags]
```

**Description**

Validates a gpkg manifest file, checking for:
- Required fields (name, version)
- Valid install or build_source specification
- Proper checksum format
- Valid build commands
- Schema compliance

**Arguments**

- `<manifest-path>`: Path to manifest file to validate

**Flags**

| Flag | Description |
|------|-------------|
| `--fix` | Attempt to automatically fix trivial issues |

**Examples**

```bash
# Validate manifest
gpkg validate ./my-tool.yaml

# Validate and auto-fix
gpkg validate ./manifest.yaml --fix

# JSON output
gpkg validate ./manifest.yaml --json
```

**Exit Codes**

- `0`: Manifest is valid
- `1`: General failure
- `6`: Validation errors found

**Related Commands**

- `gpkg info` - View manifest information
- `gpkg install` - Install from validated manifest

**JSON Output**

When `--json` is specified:

```json
{
  "valid": false,
  "errors": [
    {
      "field": "version",
      "message": "version field is required"
    }
  ],
  "warnings": [
    {
      "field": "description",
      "message": "description is recommended but not required"
    }
  ]
}
```

---

#### completion

Generate shell completion scripts.

**Synopsis**

```
gpkg completion [bash|zsh|fish] [flags]
```

**Description**

Generates shell completion scripts for bash, zsh, or fish shells. The output should be placed in your shell's completion directory.

**Arguments**

- `[shell]`: Shell type: `bash`, `zsh`, or `fish`

**Flags**

No command-specific flags.

**Examples**

```bash
# Generate bash completion
gpkg completion bash | sudo tee /usr/share/bash-completion/completions/gpkg

# Generate zsh completion
gpkg completion zsh | sudo tee /usr/share/zsh/site-functions/_gpkg

# Generate fish completion
gpkg completion fish | sudo tee /usr/share/fish/vendor_completions.d/gpkg.fish

# Or install to user directory (bash)
gpkg completion bash > ~/.local/share/bash-completion/completions/gpkg
```

**Exit Codes**

- `0`: Success
- `2`: Invalid shell specified

**Related Commands**

None.

---

## Environment Variables

gpkg respects the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `GPKG_CONFIG` | Path to config file | See config locations |
| `GPKG_INSTALL_PREFIX` | Installation prefix | `~/.gpkg` |
| `GPKG_PKGDB_PATH` | Package database path | `~/.gpkg/pkgdb.sqlite` |
| `GPKG_SOURCES_DIR` | Sources directory | `~/.gpkg/sources.d` |
| `GPKG_LOG_LEVEL` | Logging level | `info` |
| `GPKG_NO_COLOR` | Disable colors (any value) | - |
| `GPKG_OFFLINE` | Force offline mode (any value) | - |
| `LOG_LEVEL` | Alternative log level variable | `info` |
| `NO_COLOR` | Standard no-color variable | - |

**Example Usage**

```bash
# Install with custom prefix via environment
GPKG_INSTALL_PREFIX=/opt/custom gpkg install my-tool

# Run in offline mode
GPKG_OFFLINE=1 gpkg list

# Enable debug logging
GPKG_LOG_LEVEL=debug gpkg install owner/repo
```

---

## Configuration Precedence

Configuration values are resolved in the following order (highest to lowest priority):

1. **CLI flags** (e.g., `--config`, `--prefix`, `--log-level`)
2. **Environment variables** (e.g., `GPKG_INSTALL_PREFIX`)
3. **User config file** (`~/.gpkg/config.toml`)
4. **System config file** (`/etc/gpkg/config.toml`)
5. **Built-in defaults**

### Configuration File Locations

gpkg searches for configuration files in this order:

1. Path specified by `--config` flag
2. Path specified by `GPKG_CONFIG` environment variable
3. `~/.gpkg/config.toml` (user config)
4. `/etc/gpkg/config.toml` (system config)

### Merging Behavior

- CLI flags override everything
- Environment variables override config files
- User config overrides system config
- Values not specified at higher levels fall through to lower levels

**Example**

```bash
# System config: install_prefix = /usr/local
# User config: install_prefix = ~/.gpkg
# Environment: GPKG_INSTALL_PREFIX=/opt/custom
# CLI: --prefix=/tmp/test

# Resolved prefix: /tmp/test (CLI wins)
```

---

## Exit Codes

gpkg uses specific exit codes to indicate different failure conditions:

| Code | Meaning | Description |
|------|---------|-------------|
| `0` | Success | Operation completed successfully |
| `1` | General failure | Unspecified error occurred |
| `2` | Usage error | Invalid arguments or command usage |
| `3` | Network error | Remote resource unavailable, connection failed |
| `4` | Checksum failed | Checksum verification failed |
| `5` | Install/build failed | Installation or build process failed |
| `6` | Validation error | Manifest validation failed |
| `7` | Not found | Package, source, or resource not found |
| `8` | Database error | Package database error |

**Exit Code Usage in Scripts**

```bash
#!/bin/bash

gpkg install my-tool
case $? in
  0) echo "Installed successfully" ;;
  3) echo "Network error - check connectivity" ;;
  4) echo "Checksum failed - possible corruption" ;;
  7) echo "Package not found" ;;
  *) echo "Installation failed" ;;
esac
```

---

## JSON Output

Many gpkg commands support `--json` for machine-readable output, making them suitable for automation and scripting.

### Commands Supporting JSON

- `gpkg install`
- `gpkg upgrade`
- `gpkg info`
- `gpkg list`
- `gpkg list-sources`
- `gpkg search`
- `gpkg update`
- `gpkg add-source`
- `gpkg config show`
- `gpkg config get`
- `gpkg validate`

### JSON Output Features

- **Structured data**: Consistent, parseable output
- **Error handling**: Errors output as JSON to stderr when `--json` is used
- **Exit codes**: Same exit codes as regular mode
- **Logging**: Log messages suppressed; only JSON output on stdout

### Error JSON Schema

```json
{
  "error": {
    "code": 7,
    "message": "package not found",
    "detail": "no package named 'nonexistent-pkg' found in configured sources"
  }
}
```

### Success JSON Examples

#### install

```json
{
  "package": "cool-tool",
  "version": "1.2.0",
  "status": "installed",
  "files": [
    "bin/cool-tool",
    "share/man/cool-tool.1",
    "share/doc/cool-tool/README.md"
  ],
  "prefix": "/home/user/.gpkg"
}
```

#### info

```json
{
  "name": "cool-tool",
  "version": "1.2.0",
  "description": "A cool command-line tool",
  "homepage": "https://github.com/owner/cool-tool",
  "installed": {
    "version": "1.1.0",
    "prefix": "/home/user/.gpkg",
    "installed_at": "2024-01-10T08:00:00Z"
  },
  "latest_version": "1.2.0",
  "releases": [
    {
      "version": "1.2.0",
      "created_at": "2024-01-15T12:00:00Z",
      "assets": [
        {
          "platform": "linux",
          "arch": "amd64",
          "url": "https://github.com/owner/cool-tool/releases/download/v1.2.0/cool-tool-linux-amd64.tar.gz",
          "checksum": "sha256:abcdef123456..."
        }
      ]
    }
  ],
  "dependencies": [],
  "manifest": {
    "name": "cool-tool",
    "version": "1.2.0"
  }
}
```

#### list

```json
{
  "packages": [
    {
      "name": "tool1",
      "version": "1.0.0",
      "installed_at": "2024-01-15T10:30:00Z",
      "prefix": "/home/user/.gpkg"
    },
    {
      "name": "tool2",
      "version": "2.5.1",
      "installed_at": "2024-01-14T15:20:00Z",
      "prefix": "/home/user/.gpkg"
    }
  ],
  "total": 2
}
```

#### list-sources

```json
{
  "sources": [
    {
      "id": "default",
      "uri": "https://packages.gpkg.io/index.json",
      "priority": 100,
      "last_updated": "2024-01-15T10:30:00Z",
      "package_count": 127
    },
    {
      "id": "custom",
      "uri": "/opt/local-packages",
      "priority": 50,
      "last_updated": "2024-01-15T09:00:00Z",
      "package_count": 5
    }
  ]
}
```

#### search

```json
{
  "results": [
    {
      "name": "go-tool",
      "version": "1.2.0",
      "description": "A tool written in Go",
      "source": "default",
      "homepage": "https://github.com/owner/go-tool"
    },
    {
      "name": "golang-linter",
      "version": "0.5.0",
      "description": "Linter for Go code",
      "source": "default",
      "homepage": "https://github.com/owner/golang-linter"
    }
  ],
  "total": 2,
  "query": "golang"
}
```

#### upgrade

```json
{
  "upgraded": [
    {
      "package": "tool1",
      "from_version": "1.0.0",
      "to_version": "1.2.0",
      "status": "success"
    },
    {
      "package": "tool2",
      "from_version": "2.0.0",
      "to_version": "2.1.0",
      "status": "success"
    }
  ],
  "failed": [],
  "skipped": [
    {
      "package": "tool3",
      "version": "3.0.0",
      "reason": "already latest version"
    }
  ],
  "total_upgraded": 2
}
```

#### validate

```json
{
  "valid": false,
  "file": "./manifest.yaml",
  "errors": [
    {
      "field": "version",
      "message": "version field is required",
      "line": 0
    }
  ],
  "warnings": [
    {
      "field": "description",
      "message": "description is recommended but not required",
      "line": 0
    }
  ]
}
```

### Using JSON Output in Scripts

**Example: Check if package is installed**

```bash
#!/bin/bash

if gpkg list --json | jq -e '.packages[] | select(.name == "my-tool")' > /dev/null; then
  echo "my-tool is installed"
else
  echo "my-tool is not installed"
fi
```

**Example: Install and verify**

```bash
#!/bin/bash

result=$(gpkg install owner/repo --json 2>&1)
if [ $? -eq 0 ]; then
  version=$(echo "$result" | jq -r '.version')
  echo "Successfully installed version $version"
else
  error=$(echo "$result" | jq -r '.error.message')
  echo "Installation failed: $error"
fi
```

**Example: Upgrade all packages and report**

```bash
#!/bin/bash

gpkg upgrade --json --yes | jq -r '
  "Upgraded \(.total_upgraded) packages:",
  (.upgraded[] | "  - \(.package): \(.from_version) → \(.to_version)")
'
```

---

## See Also

- [Manifest Format Reference](./Manifest-Reference.md) - gpkg manifest file specification
- [Configuration Guide](../Guides/Configuration.md) - Detailed configuration documentation
- [Installation Guide](../Guides/Installation.md) - Getting started with gpkg

---

**Last Updated**: 2024-01-15  
**gpkg Version**: 0.1.0
