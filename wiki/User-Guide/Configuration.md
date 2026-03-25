# Configuration Guide

This guide explains how to configure gpkg for your needs.

## Configuration File Locations

gpkg uses a hierarchical configuration system with the following precedence (highest to lowest):

1. **CLI flags** - Command-line arguments (highest priority)
2. **Environment variables** - Shell environment variables
3. **User config** - `~/.gpkg/config.toml`
4. **System config** - `/etc/gpkg/config.toml` (lowest priority)

## Default Configuration

When you first run gpkg, it uses these defaults:

```toml
install_prefix = "~/.gpkg"
pkgdb_path = "~/.gpkg/pkgdb.sqlite"
sources_dir = "~/.gpkg/sources.d"
require_checksums = true
parallel_downloads = 4
```

## Creating a Configuration File

### User Configuration

Create a user configuration file at `~/.gpkg/config.toml`:

```bash
mkdir -p ~/.gpkg
touch ~/.gpkg/config.toml
```

### Example Configuration

```toml
# Installation directory
install_prefix = "~/.gpkg"

# Package database location
pkgdb_path = "~/.gpkg/pkgdb.sqlite"

# Package sources directory
sources_dir = "~/.gpkg/sources.d"

# Require checksum verification for downloads
require_checksums = true

# Number of parallel downloads
parallel_downloads = 4

# Default log level (error, warn, info, debug)
log_level = "info"

# Colorized output
no_color = false

# Offline mode by default
offline = false
```

## Configuration Options

### Core Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `install_prefix` | string | `~/.gpkg` | Directory where packages are installed |
| `pkgdb_path` | string | `~/.gpkg/pkgdb.sqlite` | SQLite database location |
| `sources_dir` | string | `~/.gpkg/sources.d` | Directory for package source indices |

### Security Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `require_checksums` | boolean | `true` | Require checksum verification for all downloads |

### Performance Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `parallel_downloads` | integer | `4` | Number of concurrent downloads |

### Output Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `log_level` | string | `info` | Logging verbosity (error, warn, info, debug) |
| `no_color` | boolean | `false` | Disable colorized output |
| `json` | boolean | `false` | Output in JSON format |

### Behavior Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `offline` | boolean | `false` | Disallow network operations |
| `assume_yes` | boolean | `false` | Assume "yes" for all prompts |

## Environment Variables

You can override configuration using environment variables:

```bash
# Override install prefix
export GPKG_INSTALL_PREFIX=/opt/gpkg

# Override database path
export GPKG_PKGDB_PATH=/var/lib/gpkg/pkgdb.sqlite

# Override config file location
export GPKG_CONFIG=/etc/gpkg/custom-config.toml

# Set log level
export LOG_LEVEL=debug
```

## Using the Config Command

### View Current Configuration

```bash
# Show merged configuration (all sources)
gpkg config show

# Get a specific value
gpkg config get install_prefix
```

### Modify Configuration

```bash
# Set a configuration value
gpkg config set install_prefix /opt/gpkg

# Set multiple values
gpkg config set log_level debug
gpkg config set parallel_downloads 8
```

## Custom Install Prefix

To use a custom installation directory:

### Via Configuration File

```toml
install_prefix = "/opt/gpkg"
```

### Via Environment Variable

```bash
export GPKG_INSTALL_PREFIX=/opt/gpkg
```

### Via Command-Line Flag

```bash
gpkg install package --prefix=/opt/gpkg
```

## Common Configuration Scenarios

### Minimal Network Usage

For systems with limited bandwidth:

```toml
parallel_downloads = 1
offline = false  # Set to true when working offline
```

### Development/Debug Mode

For troubleshooting:

```toml
log_level = "debug"
no_color = false
```

### CI/CD Environment

For automated environments:

```toml
assume_yes = true
json = true
no_color = true
log_level = "warn"
```

### System-Wide Installation

For shared/system installations:

```toml
install_prefix = "/usr/local"
pkgdb_path = "/var/lib/gpkg/pkgdb.sqlite"
sources_dir = "/etc/gpkg/sources.d"
```

## Troubleshooting Configuration

### Configuration Not Loading

1. Check file location: `~/.gpkg/config.toml`
2. Verify file permissions: `chmod 644 ~/.gpkg/config.toml`
3. Check TOML syntax: Use a TOML validator
4. Use `gpkg config show` to see merged configuration

### Permission Issues

If you get permission errors:

```bash
# Ensure proper ownership
chown -R $USER:$USER ~/.gpkg

# Ensure proper permissions
chmod 755 ~/.gpkg
chmod 644 ~/.gpkg/config.toml
```

### Reset to Defaults

To reset configuration:

```bash
# Remove user configuration
mv ~/.gpkg/config.toml ~/.gpkg/config.toml.bak

# gpkg will use defaults
gpkg config show
```

## See Also

- [Configuration Reference](../Reference/Configuration-Reference) - Complete list of all options
- [Basic Commands](Basic-Commands) - Using gpkg commands
- [Troubleshooting](../Troubleshooting) - Common issues and solutions

---

**Next:** Learn about [Package Management](Package-Management) workflows.
