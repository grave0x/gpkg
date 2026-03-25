# Configuration Reference

Complete reference documentation for all gpkg configuration options, environment variables, command-line flags, and configuration management.

---

## Table of Contents

1. [Overview](#overview)
2. [Configuration Files](#configuration-files)
3. [Configuration Options](#configuration-options)
4. [Environment Variables](#environment-variables)
5. [Command-Line Flags](#command-line-flags)
6. [Precedence Rules](#precedence-rules)
7. [Validation Rules](#validation-rules)
8. [Configuration Management](#configuration-management)
9. [Common Configuration Scenarios](#common-configuration-scenarios)
10. [Troubleshooting](#troubleshooting)

---

## Overview

gpkg uses a layered configuration system that allows fine-grained control over package management behavior. Configuration can be set through:

- **Configuration files** (YAML format)
- **Environment variables** (GPKG_* prefix)
- **Command-line flags** (highest priority)

The configuration system follows a clear precedence order, ensuring predictable behavior in different environments (development, production, CI/CD).

### Configuration Architecture

```
┌─────────────────────────┐
│   CLI Flags             │ ← Highest priority
├─────────────────────────┤
│   Environment Variables │
├─────────────────────────┤
│   User Config File      │
├─────────────────────────┤
│   System Config File    │
├─────────────────────────┤
│   Built-in Defaults     │ ← Lowest priority
└─────────────────────────┘
```

---

## Configuration Files

### File Locations

gpkg searches for configuration files in the following locations (first found wins):

#### User Configuration
1. `~/.gpkg/config.yaml` (primary, recommended)
2. `~/.gpkg.yaml` (alternative)
3. `~/.config/gpkg/config.yaml` (XDG-compliant)

#### System Configuration
1. `/etc/gpkg/config.yaml` (primary)
2. `/etc/gpkg.yaml` (alternative)

#### Custom Configuration
- Specified via `--config` flag or `GPKG_CONFIG` environment variable

### File Format

gpkg uses **YAML** syntax for configuration files. The file must have a `.yaml` extension.

### Example Configuration File

**Location**: `~/.gpkg/config.yaml`

```yaml
# Installation prefix - where packages are installed
prefix: ~/.gpkg

# Cache directory for downloaded packages and build artifacts
cache_dir: ~/.gpkg/cache

# Sources registry file - package source definitions
sources_file: ~/.gpkg/sources.json

# Logging level: error, warn, info, debug
log_level: info

# Enable colored terminal output
color: true

# Require strict checksum verification for downloads
strict_checksum: true

# Network operation timeout in seconds
network_timeout: 30
```

### Minimal Configuration

If no configuration file exists, gpkg uses built-in defaults. A minimal config file might only override specific values:

```yaml
# Minimal config - only override what you need
prefix: /opt/gpkg
log_level: debug
```

---

## Configuration Options

All configuration options that can be set in configuration files:

| Option | Type | Default | Description | Valid Values |
|--------|------|---------|-------------|--------------|
| `prefix` | string | `~/.gpkg` | Installation prefix directory where packages are installed | Any valid directory path |
| `cache_dir` | string | `~/.gpkg/cache` | Cache directory for downloads and build artifacts | Any valid directory path |
| `sources_file` | string | `~/.gpkg/sources.json` | Path to package sources registry file | Any valid file path |
| `log_level` | string | `info` | Logging verbosity level | `error`, `warn`, `info`, `debug` |
| `color` | boolean | `true` | Enable colored terminal output | `true`, `false` |
| `strict_checksum` | boolean | `true` | Require checksum verification for all downloads | `true`, `false` |
| `network_timeout` | integer | `30` | Network operation timeout in seconds | Positive integer (recommended: 15-300) |

### Option Details

#### prefix
- **Purpose**: Root directory where all packages are installed
- **Path Expansion**: Tilde (`~`) is expanded to user home directory
- **Subdirectories Created**:
  - `bin/` - Executable binaries
  - `lib/` - Shared libraries
  - `include/` - Header files
  - `share/` - Shared data files
  - `pkgdb.sqlite` - Package database
- **Example**: 
  ```yaml
  prefix: /opt/gpkg
  ```

#### cache_dir
- **Purpose**: Temporary storage for downloads before installation
- **Contents**:
  - Downloaded release archives
  - Build artifacts
  - Temporary files during installation
- **Cleanup**: Can be safely deleted; will be recreated as needed
- **Example**:
  ```yaml
  cache_dir: /var/cache/gpkg
  ```

#### sources_file
- **Purpose**: JSON file containing package source definitions
- **Format**: JSON array of source objects
- **Created Automatically**: If it doesn't exist, gpkg creates an empty sources file
- **Example**:
  ```yaml
  sources_file: /etc/gpkg/sources.json
  ```

#### log_level
- **Purpose**: Controls verbosity of log output
- **Levels** (least to most verbose):
  - `error` - Only critical errors
  - `warn` - Warnings and errors
  - `info` - General information, warnings, and errors (default)
  - `debug` - Detailed debugging information
- **Case Insensitive**: `INFO`, `info`, `Info` all work
- **Override**: Can be overridden with `--log-level` flag
- **Example**:
  ```yaml
  log_level: debug
  ```

#### color
- **Purpose**: Enable/disable colored terminal output
- **Auto-Detection**: Automatically disabled when output is not a TTY
- **CI/CD**: Set to `false` for cleaner logs in automated environments
- **Override**: Use `--no-color` flag to disable
- **Example**:
  ```yaml
  color: false
  ```

#### strict_checksum
- **Purpose**: Controls checksum verification behavior
- **When `true`** (default, recommended):
  - All downloads must have checksums in manifest
  - Installation fails if checksum doesn't match
  - Ensures package integrity and security
- **When `false`** (use with caution):
  - Downloads without checksums are allowed
  - Security warning is displayed
  - Use only for development/testing
- **Example**:
  ```yaml
  strict_checksum: true
  ```

#### network_timeout
- **Purpose**: Maximum time to wait for network operations
- **Applies To**:
  - Package downloads
  - Source index fetching
  - Git clone operations
- **Units**: Seconds
- **Recommended Range**: 30-120 seconds
- **Example**:
  ```yaml
  network_timeout: 60
  ```

---

## Environment Variables

Environment variables override configuration file settings. All gpkg environment variables use the `GPKG_` prefix.

### Standard Environment Variables

| Variable | Type | Maps To | Description | Example |
|----------|------|---------|-------------|---------|
| `GPKG_PREFIX` | string | `prefix` | Override installation prefix | `export GPKG_PREFIX=/usr/local/gpkg` |
| `GPKG_CACHE_DIR` | string | `cache_dir` | Override cache directory | `export GPKG_CACHE_DIR=/tmp/gpkg-cache` |
| `GPKG_LOG_LEVEL` | string | `log_level` | Override log level | `export GPKG_LOG_LEVEL=debug` |
| `GPKG_CONFIG` | string | - | Specify custom config file location | `export GPKG_CONFIG=/etc/custom/gpkg.yaml` |

### Additional Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `LOG_LEVEL` | Alternative to GPKG_LOG_LEVEL (lower priority) | `export LOG_LEVEL=warn` |
| `NO_COLOR` | Disable colored output (any value disables) | `export NO_COLOR=1` |
| `HOME` | User home directory (used for ~ expansion) | Automatically set by system |

### Using Environment Variables

#### Temporary Override (Single Command)
```bash
GPKG_LOG_LEVEL=debug gpkg install package-name
```

#### Session Override
```bash
export GPKG_PREFIX=/opt/packages
gpkg install package1
gpkg install package2
```

#### Shell Profile (Permanent)
Add to `~/.bashrc` or `~/.zshrc`:
```bash
export GPKG_PREFIX="$HOME/mypackages"
export GPKG_LOG_LEVEL=info
export GPKG_CACHE_DIR="/var/cache/gpkg"
```

---

## Command-Line Flags

Command-line flags have the highest priority and override all other configuration sources.

### Global Flags

These flags apply to **all commands**:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config <path>` | `-c` | string | - | Path to custom config file |
| `--json` | - | boolean | `false` | Output in JSON format (machine-readable) |
| `--yes` | `-y` | boolean | `false` | Assume "yes" to all prompts (non-interactive) |
| `--dry-run` | - | boolean | `false` | Preview actions without making changes |
| `--log-level <level>` | - | string | `info` | Set logging level |
| `--verbose` | `-v` | counter | `0` | Increase verbosity (stackable: `-v`, `-vv`, `-vvv`) |
| `--quiet` | - | boolean | `false` | Minimal output (conflicts with `--verbose`) |
| `--no-color` | - | boolean | `false` | Disable colored output |
| `--offline` | - | boolean | `false` | Disable all network operations |
| `--help` | `-h` | - | - | Show help for command |
| `--version` | `-V` | - | - | Show gpkg version |

### Command-Specific Flags

#### install Command

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-release` | boolean | `false` | Install from release binary (default behavior) |
| `--from-source` | boolean | `false` | Build and install from source |
| `--prefix <path>` | string | - | Override installation prefix for this package |
| `--force` | boolean | `false` | Overwrite existing installation |
| `--no-deps` | boolean | `false` | Skip dependency installation (dangerous) |

**Examples**:
```bash
# Install with custom prefix
gpkg install tool --prefix /opt/custom

# Force reinstall
gpkg install tool --force

# Build from source
gpkg install tool --from-source

# Non-interactive installation
gpkg install tool --yes --json
```

#### list Command

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--installed` | boolean | `true` | Show installed packages |
| `--available` | boolean | `false` | Show available packages from sources |
| `--filter <pattern>` | string | - | Filter by name pattern |
| `--sort <field>` | string | `name` | Sort by: `name`, `version`, `date` |

**Examples**:
```bash
# List all installed packages
gpkg list --installed

# List available packages
gpkg list --available

# Filter packages
gpkg list --filter "dev-*"
```

#### info Command

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--raw` | boolean | `false` | Show raw manifest file |
| `--parsed` | boolean | `false` | Show parsed manifest structure |
| `--deps-tree` | boolean | `false` | Show dependency tree |

**Examples**:
```bash
# Basic package info
gpkg info owner/repo

# Show dependency tree
gpkg info owner/repo --deps-tree

# Raw manifest
gpkg info owner/repo --raw --json
```

#### search Command

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--source <id>` | string | - | Search in specific source |

**Examples**:
```bash
# Search all sources
gpkg search keyword

# Search specific source
gpkg search keyword --source github-releases
```

#### upgrade Command

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | boolean | `false` | Upgrade all packages |
| `--dry-run` | boolean | `false` | Preview upgrades without installing |

**Examples**:
```bash
# Upgrade specific packages
gpkg upgrade package1 package2

# Upgrade all packages
gpkg upgrade --all

# Preview upgrades
gpkg upgrade --all --dry-run --json
```

#### uninstall Command

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--prefix <path>` | string | - | Installation prefix override |
| `--force` | boolean | `false` | Skip confirmation prompts |

**Examples**:
```bash
# Uninstall package
gpkg uninstall package-name

# Force uninstall without confirmation
gpkg uninstall package-name --force --yes
```

#### rollback Command

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--to-version <ver>` | string | - | Version to rollback to (required) |
| `--prefix <path>` | string | - | Installation prefix override |

**Examples**:
```bash
# Rollback to specific version
gpkg rollback package-name --to-version 1.2.0
```

#### validate Command

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--fix` | boolean | `false` | Attempt to auto-fix trivial issues |

**Examples**:
```bash
# Validate manifest
gpkg validate manifest.yaml

# Validate and auto-fix
gpkg validate manifest.yaml --fix
```

---

## Precedence Rules

Configuration values are determined by the following precedence order (highest to lowest):

### 1. Command-Line Flags (Highest Priority)
```bash
gpkg install package --log-level debug
# log_level = debug (from CLI flag)
```

### 2. Environment Variables
```bash
export GPKG_LOG_LEVEL=warn
gpkg install package
# log_level = warn (from environment)
```

### 3. User Configuration File
```yaml
# ~/.gpkg/config.yaml
log_level: info
```

### 4. System Configuration File
```yaml
# /etc/gpkg/config.yaml
log_level: error
```

### 5. Built-in Defaults (Lowest Priority)
```go
// Hardcoded in application
LogLevel: "info"
```

### Precedence Examples

#### Example 1: Multiple Sources
```bash
# System config
# /etc/gpkg/config.yaml
prefix: /usr/local/gpkg
log_level: warn

# User config  
# ~/.gpkg/config.yaml
prefix: ~/.gpkg
# log_level not specified, inherits 'warn' from system config

# Environment
export GPKG_LOG_LEVEL=debug

# Command
gpkg install package --prefix /opt/packages

# Effective configuration:
# prefix: /opt/packages (CLI flag wins)
# log_level: debug (environment variable wins)
```

#### Example 2: Partial Override
```bash
# Config file
# ~/.gpkg/config.yaml
prefix: ~/.gpkg
cache_dir: ~/.gpkg/cache
log_level: info
network_timeout: 60

# Environment
export GPKG_LOG_LEVEL=debug

# Effective configuration:
# prefix: ~/.gpkg (from config)
# cache_dir: ~/.gpkg/cache (from config)
# log_level: debug (from environment - overrides config)
# network_timeout: 60 (from config)
```

### Config File Search Order

When multiple config file locations exist:

#### User Config Search (First Found Wins)
1. `~/.gpkg/config.yaml`
2. `~/.gpkg.yaml`
3. `~/.config/gpkg/config.yaml`

#### System Config Search (First Found Wins)
1. `/etc/gpkg/config.yaml`
2. `/etc/gpkg.yaml`

**Important**: System config is loaded first, then user config *merges* over it, preserving unset values from system config.

---

## Validation Rules

### Configuration File Validation

#### YAML Syntax
- **Must be valid YAML**: Syntax errors cause config load to fail
- **Case Sensitive**: Keys must match exactly (`prefix`, not `Prefix`)
- **Type Checking**: Values must match expected types

**Invalid**:
```yaml
prefix: ~/.gpkg
log_level: 123  # Error: must be string
color: "yes"    # Error: must be boolean (true/false)
```

**Valid**:
```yaml
prefix: ~/.gpkg
log_level: debug
color: true
network_timeout: 60
```

#### Field Validation

##### log_level
- **Valid Values**: `error`, `warn`, `info`, `debug` (case-insensitive)
- **Invalid Value Behavior**: Falls back to `info` with warning
- **Examples**:
  ```yaml
  log_level: DEBUG  # OK (normalized to 'debug')
  log_level: trace  # Invalid → defaults to 'info' with warning
  ```

##### color and strict_checksum
- **Valid Values**: `true`, `false` (boolean)
- **Invalid Values**: Strings like `"yes"`, `"no"`, `1`, `0` are rejected
- **Example**:
  ```yaml
  color: true       # OK
  color: false      # OK
  color: "true"     # Error: must be boolean, not string
  ```

##### network_timeout
- **Valid Values**: Positive integers
- **Invalid Values**: 
  - Zero or negative numbers → rejected
  - Floats → truncated to integer
  - Strings → rejected
- **Recommended Range**: 15-300 seconds
- **Example**:
  ```yaml
  network_timeout: 30    # OK
  network_timeout: 0     # Error: must be positive
  network_timeout: -10   # Error: must be positive
  network_timeout: "30"  # Error: must be integer, not string
  ```

##### Paths (prefix, cache_dir, sources_file)
- **Valid Values**: Any string representing a file system path
- **Tilde Expansion**: `~` is expanded to user home directory
- **Relative Paths**: Allowed (resolved from current working directory)
- **Validation**: Path doesn't need to exist (created on first use)
- **Examples**:
  ```yaml
  prefix: ~/.gpkg              # OK (expands to /home/user/.gpkg)
  prefix: /opt/gpkg            # OK (absolute path)
  prefix: ../packages          # OK (relative path)
  cache_dir: $HOME/cache/gpkg  # OK (shell variable - must be expanded by shell)
  ```

### Manifest Validation

When validating package manifests (`gpkg validate manifest.yaml`):

#### Required Fields
- `name` - Package name (non-empty string)
- `version` - Semantic version (e.g., `1.2.3`)

#### Installation Method (One Required)
- **Either** `release_assets` (for binary installation)
- **Or** `build` (for source build)
- **Not both** in most cases

#### Release Assets Validation
```yaml
release_assets:
  - url: "https://..."      # Required
    platform: "linux/amd64" # Required
    checksum:               # Required if strict_checksum=true
      sha256: "abc123..."
```

#### Build Configuration Validation
```yaml
build:
  repo: "https://..."       # Required (git repository URL)
  build_steps:              # Required (at least one command)
    - "make"
    - "make install"
```

#### Common Validation Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `missing required field: name` | `name` not specified | Add `name: package-name` |
| `missing required field: version` | `version` not specified | Add `version: 1.0.0` |
| `no installation method` | Neither `release_assets` nor `build` | Add one installation method |
| `invalid version format` | Version not semver | Use format `X.Y.Z` |
| `checksum required` | Missing checksum with `strict_checksum=true` | Add SHA256/SHA512 checksum |
| `invalid URL` | Malformed URL in `url` or `repo` | Fix URL syntax |

### Runtime Validation

#### Directory Creation
- Directories are created automatically if they don't exist
- Parent directories must be writable
- Fails with permission error if cannot create

#### Network Validation
- URLs must be valid HTTP/HTTPS
- Timeout enforced per `network_timeout` setting
- Certificate validation performed (unless insecure flag used)

#### Checksum Validation
- Downloaded files must match manifest checksums
- Supported: SHA256, SHA512
- Fails installation if mismatch detected
- Can be disabled with `strict_checksum: false` (not recommended)

---

## Configuration Management

gpkg provides commands to view and modify configuration.

### View Configuration

#### Show All Configuration
```bash
# Human-readable format
gpkg config show

# JSON format (machine-readable)
gpkg config show --json
```

**Output Example**:
```
Configuration:
  Prefix:          /home/user/.gpkg
  Cache Dir:       /home/user/.gpkg/cache
  Sources File:    /home/user/.gpkg/sources.json
  Log Level:       info
  Color:           true
  Strict Checksum: true
  Network Timeout: 30s

Configuration Sources:
  - User config: /home/user/.gpkg/config.yaml
```

#### Get Specific Value
```bash
gpkg config get prefix
# Output: /home/user/.gpkg

gpkg config get log_level
# Output: info

gpkg config get network_timeout
# Output: 30
```

### Modify Configuration

#### Set Configuration Value
```bash
# Set installation prefix
gpkg config set prefix /opt/gpkg

# Set log level
gpkg config set log_level debug

# Set network timeout
gpkg config set network_timeout 60

# Disable color output
gpkg config set color false

# Enable strict checksums
gpkg config set strict_checksum true
```

**Note**: `config set` writes to user config file (`~/.gpkg/config.yaml`), creating it if it doesn't exist.

#### Settable Keys
- `prefix`
- `cache_dir`
- `sources_file`
- `log_level`
- `color`
- `strict_checksum`
- `network_timeout`

### Reset Configuration

#### Remove User Configuration
```bash
# Backup existing config
cp ~/.gpkg/config.yaml ~/.gpkg/config.yaml.backup

# Remove user config (falls back to system config or defaults)
rm ~/.gpkg/config.yaml

# Verify defaults
gpkg config show
```

#### Create Default Configuration
```bash
# Generate a default config file
cat > ~/.gpkg/config.yaml <<EOF
prefix: ~/.gpkg
cache_dir: ~/.gpkg/cache
sources_file: ~/.gpkg/sources.json
log_level: info
color: true
strict_checksum: true
network_timeout: 30
EOF
```

---

## Common Configuration Scenarios

### Development Environment

**Goal**: Verbose logging, local prefix, relaxed checksums

**Config**: `~/.gpkg/config.yaml`
```yaml
prefix: ~/dev/packages
cache_dir: ~/dev/cache
log_level: debug
strict_checksum: false
network_timeout: 60
```

**Usage**:
```bash
gpkg install tool --from-source -vv
```

---

### Production Environment

**Goal**: Secure, system-wide, strict verification

**Config**: `/etc/gpkg/config.yaml`
```yaml
prefix: /opt/packages
cache_dir: /var/cache/gpkg
sources_file: /etc/gpkg/sources.json
log_level: warn
strict_checksum: true
network_timeout: 30
color: false
```

**Usage**:
```bash
sudo gpkg install tool --yes --json
```

---

### CI/CD Pipeline

**Goal**: Non-interactive, JSON output, no colors

**Environment** (`.gitlab-ci.yml` or similar):
```yaml
variables:
  GPKG_PREFIX: /ci-build/packages
  GPKG_LOG_LEVEL: warn
  NO_COLOR: 1
```

**Usage**:
```bash
gpkg install tool --yes --json --dry-run
gpkg install tool --yes --json
```

---

### Offline Environment

**Goal**: No network access, use cached/local packages

**Config**: `~/.gpkg/config.yaml`
```yaml
prefix: ~/.gpkg
cache_dir: /mnt/offline-cache
log_level: info
```

**Usage**:
```bash
# All operations use --offline flag
gpkg install ./local-manifest.yaml --offline
gpkg list --installed --offline
```

---

### Multi-User System

**Goal**: Shared packages, per-user caching

**System Config**: `/etc/gpkg/config.yaml`
```yaml
prefix: /usr/local/packages
sources_file: /etc/gpkg/sources.json
strict_checksum: true
log_level: info
```

**User Config**: `~/.gpkg/config.yaml`
```yaml
cache_dir: ~/.cache/gpkg
log_level: debug  # User preference
```

**Result**:
- Packages install to `/usr/local/packages` (requires sudo)
- Each user has their own cache in `~/.cache/gpkg`
- Users can set their own log level

---

### Docker Container

**Goal**: Minimal, reproducible, fast

**Dockerfile**:
```dockerfile
FROM ubuntu:22.04

# Install gpkg
RUN curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg-linux-amd64.tar.gz \
    && tar -xzf gpkg-linux-amd64.tar.gz \
    && mv gpkg-linux-amd64 /usr/local/bin/gpkg \
    && chmod +x /usr/local/bin/gpkg

# Configure gpkg
ENV GPKG_PREFIX=/usr/local/packages \
    GPKG_CACHE_DIR=/var/cache/gpkg \
    GPKG_LOG_LEVEL=warn

# Install packages
RUN gpkg install tool --yes --json

# Add to PATH
ENV PATH="/usr/local/packages/bin:$PATH"
```

---

## Troubleshooting

### Configuration Not Loading

**Symptom**: Configuration file changes are ignored

**Diagnosis**:
```bash
# Check which config file is being used
gpkg config show

# Verify file exists and is readable
ls -la ~/.gpkg/config.yaml

# Check YAML syntax
cat ~/.gpkg/config.yaml
```

**Common Causes**:
1. **Wrong file location**: Config in wrong directory
   - **Fix**: Move to `~/.gpkg/config.yaml`

2. **YAML syntax error**: Invalid YAML
   - **Fix**: Validate YAML syntax, check indentation

3. **File permissions**: Not readable
   - **Fix**: `chmod 644 ~/.gpkg/config.yaml`

4. **Environment override**: Environment variable overriding config
   - **Fix**: Check `env | grep GPKG_`

5. **CLI flag override**: Command-line flag taking precedence
   - **Fix**: Remove CLI flags to test config

---

### Invalid Log Level

**Symptom**: Warning about invalid log level

**Error**:
```
WARN: Invalid log level 'trace', using default 'info'
```

**Fix**: Use valid log level (`error`, `warn`, `info`, `debug`)
```yaml
log_level: debug  # Not 'trace'
```

---

### Permission Denied

**Symptom**: Cannot create directories or write files

**Error**:
```
Error: failed to create directory /opt/gpkg: permission denied
```

**Diagnosis**:
```bash
# Check prefix ownership
ls -ld /opt/gpkg

# Check parent directory permissions
ls -ld /opt
```

**Solutions**:

1. **Use sudo** (system-wide installation):
   ```bash
   sudo gpkg install package
   ```

2. **Change prefix** (user installation):
   ```bash
   gpkg install package --prefix ~/.local/packages
   ```

3. **Fix permissions**:
   ```bash
   sudo chown -R $USER:$USER /opt/gpkg
   ```

---

### Network Timeout

**Symptom**: Downloads timing out

**Error**:
```
Error: download failed: context deadline exceeded
```

**Fix**: Increase network timeout
```yaml
network_timeout: 120  # Increase to 120 seconds
```

Or via CLI:
```bash
# Set temporarily via environment
GPKG_NETWORK_TIMEOUT=120 gpkg install package
```

---

### Checksum Mismatch

**Symptom**: Installation fails due to checksum verification

**Error**:
```
Error: checksum mismatch: expected abc123..., got def456...
```

**Diagnosis**:
1. Network issue causing corrupted download
2. Manifest has wrong checksum
3. Package was updated without updating manifest

**Solutions**:

1. **Retry download** (clears cache):
   ```bash
   rm -rf ~/.gpkg/cache/*
   gpkg install package
   ```

2. **Verify manifest checksum** is correct

3. **Disable strict checksum** (development only):
   ```yaml
   strict_checksum: false
   ```
   ```bash
   gpkg install package
   ```

---

### Config File Not Created

**Symptom**: `gpkg config set` doesn't create file

**Diagnosis**:
```bash
# Check if .gpkg directory exists
ls -ld ~/.gpkg

# Check parent permissions
ls -ld ~
```

**Fix**: Create directory manually
```bash
mkdir -p ~/.gpkg
gpkg config set prefix ~/.gpkg
```

---

### Environment Variables Not Working

**Symptom**: Environment variables are ignored

**Diagnosis**:
```bash
# Check if variable is set
echo $GPKG_PREFIX

# Check if exported
export | grep GPKG
```

**Fix**: Ensure variable is exported
```bash
export GPKG_PREFIX=/opt/packages  # Not just 'GPKG_PREFIX=/opt/packages'
gpkg install package
```

---

### Colors Not Showing

**Symptom**: Output not colored despite `color: true`

**Diagnosis**:
```bash
# Check if stdout is a TTY
test -t 1 && echo "TTY" || echo "Not a TTY"

# Check NO_COLOR environment
echo $NO_COLOR
```

**Common Causes**:
1. Output piped or redirected (not a TTY)
2. `NO_COLOR` environment variable set
3. `--no-color` flag used
4. Terminal doesn't support colors

**Fix**:
```bash
# Force color (if supported)
unset NO_COLOR
gpkg install package

# Check terminal type
echo $TERM
```

---

### Can't Find Package Database

**Symptom**: Error about missing pkgdb.sqlite

**Error**:
```
Error: failed to open package database: no such file or directory
```

**Diagnosis**:
```bash
# Check if database exists
ls -la ~/.gpkg/pkgdb.sqlite

# Check prefix configuration
gpkg config get prefix
```

**Fix**: Database is created automatically on first install
```bash
# Install a package to initialize database
gpkg install package
```

Or create manually:
```bash
touch ~/.gpkg/pkgdb.sqlite
```

---

## Summary

### Quick Reference Card

| Aspect | Details |
|--------|---------|
| **Config Format** | YAML (`.yaml` extension) |
| **Primary Location** | `~/.gpkg/config.yaml` |
| **Total Options** | 7 core configuration options |
| **Environment Vars** | 3 primary (`GPKG_PREFIX`, `GPKG_CACHE_DIR`, `GPKG_LOG_LEVEL`) |
| **Global Flags** | 11 global flags applying to all commands |
| **Precedence** | CLI flags > Env vars > User config > System config > Defaults |
| **Validation** | YAML syntax, type checking, value ranges |
| **Management** | `gpkg config show/get/set` commands |

### Configuration Checklist

- [ ] Config file in correct location (`~/.gpkg/config.yaml`)
- [ ] Valid YAML syntax (check with validator)
- [ ] Correct data types (strings quoted, booleans true/false)
- [ ] Valid log level (error/warn/info/debug)
- [ ] Positive integer for network_timeout
- [ ] Paths exist or are writable
- [ ] Environment variables exported (if used)
- [ ] No conflicting CLI flags

### Default Values (Reference)

```yaml
prefix: ~/.gpkg
cache_dir: ~/.gpkg/cache
sources_file: ~/.gpkg/sources.json
log_level: info
color: true
strict_checksum: true
network_timeout: 30
```

---

## See Also

- [Installation Guide](../Installation-Guide.md) - Setting up gpkg
- [CLI Reference](../CLI-Reference.md) - Complete command documentation
- [Manifest Format](../Manifest-Format.md) - Package manifest specification
- [Security Best Practices](../Security-Best-Practices.md) - Secure configuration guidelines
- [Troubleshooting Guide](../Troubleshooting-Guide.md) - Common issues and solutions

---

**Last Updated**: 2024
**gpkg Version**: 0.1.0+
