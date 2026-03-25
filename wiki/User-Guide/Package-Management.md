# Package Management Guide

This guide covers the complete lifecycle of managing packages with gpkg.

## Table of Contents

- [Installing Packages](#installing-packages)
- [Upgrading Packages](#upgrading-packages)
- [Uninstalling Packages](#uninstalling-packages)
- [Rolling Back Packages](#rolling-back-packages)
- [Package Sources](#package-sources)
- [Package Information](#package-information)
- [Best Practices](#best-practices)

## Installing Packages

### From GitHub Release

The simplest way to install a package:

```bash
# Install latest release
gpkg install owner/repo

# Install specific version
gpkg install owner/repo@v1.2.0
```

### From Local Manifest

Install from a local manifest file:

```bash
# Install from manifest
gpkg install ./path/to/manifest.yaml

# Install with custom prefix
gpkg install ./manifest.yaml --prefix=/opt
```

### From Source

Build and install from source:

```bash
# Build from source instead of using release
gpkg install owner/repo --from-source
```

### Installation Options

```bash
# Custom installation prefix
gpkg install package --prefix=/custom/path

# Force reinstall (overwrite existing)
gpkg install package --force

# Skip dependency installation
gpkg install package --no-deps

# Dry run (preview without installing)
gpkg install package --dry-run

# Non-interactive mode
gpkg install package --yes
```

### Example Installation Workflow

```bash
# Search for a package
gpkg search cool-tool

# View package information
gpkg info owner/cool-tool

# Install the package
gpkg install owner/cool-tool

# Verify installation
gpkg list --installed | grep cool-tool
```

## Upgrading Packages

### Upgrade All Packages

```bash
# Upgrade all installed packages
gpkg upgrade

# Preview upgrades without applying
gpkg upgrade --dry-run

# Non-interactive upgrade
gpkg upgrade --yes
```

### Upgrade Specific Packages

```bash
# Upgrade one package
gpkg upgrade package-name

# Upgrade multiple packages
gpkg upgrade pkg1 pkg2 pkg3
```

### Upgrade Options

```bash
# Parallel upgrades (faster)
gpkg upgrade --concurrency 4

# Skip confirmation prompts
gpkg upgrade --yes

# Show what would be upgraded
gpkg upgrade --dry-run
```

## Uninstalling Packages

### Basic Uninstall

```bash
# Uninstall a package
gpkg uninstall package-name

# Skip confirmation
gpkg uninstall package-name --yes
```

### Purge (Complete Removal)

```bash
# Remove package and all data
gpkg uninstall package-name --purge
```

### Safety Features

- gpkg only removes files it installed (tracked in database)
- Files outside install_prefix are never removed
- Confirmation required unless `--yes` is used

## Rolling Back Packages

### Rollback to Previous Version

```bash
# List available versions
gpkg info package-name

# Rollback to specific version
gpkg rollback package-name v1.0.0

# Rollback to previous version
gpkg rollback package-name
```

## Package Sources

### Adding Sources

```bash
# Add HTTP/HTTPS source
gpkg add-source https://packages.example.com/index.json

# Add with custom name
gpkg add-source https://example.com/index.json --name my-source

# Add with priority
gpkg add-source https://example.com/index.json --priority 10
```

### Managing Sources

```bash
# List all sources
gpkg list-sources

# Update package indices
gpkg update

# Force update (ignore cache)
gpkg update --force

# Remove a source
gpkg remove-source source-id

# Remove with confirmation skip
gpkg remove-source source-id --yes
```

### Source Priority

Sources with higher priority are checked first for packages. Default priority is 0.

```bash
# Add high-priority source
gpkg add-source https://priority.example.com/index.json --priority 100

# Add low-priority source
gpkg add-source https://fallback.example.com/index.json --priority -10
```

## Package Information

### View Package Details

```bash
# Show package information
gpkg info package-name

# Show raw manifest
gpkg info package-name --raw

# Show parsed/structured view
gpkg info package-name --parsed

# Show dependency tree
gpkg info package-name --deps-tree
```

### Search for Packages

```bash
# Search all sources
gpkg search term

# Search specific source
gpkg search term --source my-source

# JSON output for scripting
gpkg search term --json
```

### List Packages

```bash
# List installed packages
gpkg list --installed

# List available packages
gpkg list --available

# Filter by name pattern
gpkg list --filter "tool"

# Sort by version
gpkg list --sort version

# JSON output
gpkg list --json
```

## Best Practices

### 1. Regular Updates

Keep package indices fresh:

```bash
# Update weekly or before installing
gpkg update
```

### 2. Review Before Installing

Always check package information:

```bash
# Review package details
gpkg info package-name

# Check what will be installed
gpkg install package-name --dry-run
```

### 3. Use Checksum Verification

Never disable checksum verification in production:

```toml
# In config file
require_checksums = true
```

### 4. Test Upgrades First

Use dry-run before upgrading:

```bash
# See what would change
gpkg upgrade --dry-run

# If satisfied, upgrade
gpkg upgrade
```

### 5. Keep Track of Installed Packages

Regularly review installed packages:

```bash
# List all installed
gpkg list --installed

# Remove unused packages
gpkg uninstall old-package
```

### 6. Use Appropriate Install Prefix

For user-specific tools:
```bash
# Default user installation
gpkg install package
```

For system-wide tools:
```bash
# System installation (requires sudo)
sudo gpkg install package --prefix=/usr/local
```

### 7. Backup Package Database

The database is critical for tracking installations:

```bash
# Backup database
cp ~/.gpkg/pkgdb.sqlite ~/.gpkg/pkgdb.sqlite.backup

# Restore if needed
cp ~/.gpkg/pkgdb.sqlite.backup ~/.gpkg/pkgdb.sqlite
```

## Advanced Workflows

### Batch Installation

Install multiple packages:

```bash
# Create package list
cat > packages.txt <<EOF
owner/repo1
owner/repo2
owner/repo3
EOF

# Install all
while read pkg; do gpkg install "$pkg"; done < packages.txt
```

### Scripted Installation for CI/CD

```bash
#!/bin/bash
set -e

# Non-interactive, JSON output
export GPKG_ASSUME_YES=true

# Install required tools
gpkg update
gpkg install --json owner/tool1
gpkg install --json owner/tool2

# Verify installations
gpkg list --installed --json
```

### Custom Installation Workflows

```bash
# Install to custom location with specific version
gpkg install owner/repo@v2.0.0 --prefix=/opt/mytools --yes

# Verify before production
gpkg info owner/repo --deps-tree

# Monitor installation
gpkg install owner/repo --verbose
```

## Troubleshooting

### Installation Fails

1. Check network connectivity
2. Verify checksum in manifest
3. Check disk space
4. Review logs with `--verbose` or `--log-level debug`

### Upgrade Issues

1. Try dry-run first: `gpkg upgrade --dry-run`
2. Upgrade packages individually
3. Check for conflicting dependencies

### Database Issues

```bash
# Check database integrity
sqlite3 ~/.gpkg/pkgdb.sqlite "PRAGMA integrity_check;"

# Backup and reset if corrupted
mv ~/.gpkg/pkgdb.sqlite ~/.gpkg/pkgdb.sqlite.bad
# Re-install packages
```

## See Also

- [Basic Commands](Basic-Commands) - Command reference
- [Configuration](Configuration) - Configure gpkg settings
- [Advanced Usage](Advanced-Usage) - Power user features
- [Troubleshooting](../Troubleshooting) - Common issues

---

**Next:** Explore [Advanced Usage](Advanced-Usage) for power user features.
