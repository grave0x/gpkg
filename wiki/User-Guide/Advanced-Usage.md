# Advanced Usage Guide

This guide covers advanced features and power-user workflows for gpkg.

## Table of Contents

- [Advanced Installation Techniques](#advanced-installation-techniques)
- [Scripting and Automation](#scripting-and-automation)
- [Custom Package Sources](#custom-package-sources)
- [Working with Manifests](#working-with-manifests)
- [Dependency Management](#dependency-management)
- [Performance Optimization](#performance-optimization)
- [Security Best Practices](#security-best-practices)
- [Advanced Troubleshooting](#advanced-troubleshooting)

## Advanced Installation Techniques

### Parallel Installations

Install multiple packages concurrently:

```bash
# Enable parallel downloads in config
cat >> ~/.gpkg/config.toml <<EOF
parallel_downloads = 8
EOF

# Install multiple packages in background
for pkg in owner/pkg1 owner/pkg2 owner/pkg3; do
    gpkg install "$pkg" &
done
wait
```

### Conditional Installation

Install based on system detection:

```bash
#!/bin/bash

# Detect platform
case "$(uname -s)" in
    Linux*)  gpkg install linux-specific-tool ;;
    Darwin*) gpkg install mac-specific-tool ;;
    *)       echo "Unsupported platform" ;;
esac
```

### Version Pinning

Lock packages to specific versions:

```bash
# Install specific version
gpkg install owner/repo@v1.2.3

# Create version lock file
cat > package-lock.txt <<EOF
owner/repo1@v1.0.0
owner/repo2@v2.3.1
EOF

# Install from lock file
while IFS= read -r pkg; do
    gpkg install "$pkg"
done < package-lock.txt
```

### Custom Build from Source

Override build steps for source installations:

```yaml
# custom-manifest.yaml
name: my-tool
version: 1.0.0
build:
  repo: https://github.com/owner/tool.git
  build_steps:
    - "export CGO_ENABLED=0"
    - "go build -ldflags='-s -w' -o bin/tool ."
  build_env:
    - "GOOS=linux"
    - "GOARCH=amd64"
```

## Scripting and Automation

### JSON Output for Parsing

All commands support `--json` for machine-readable output:

```bash
# Get installed packages as JSON
gpkg list --installed --json | jq '.packages[].name'

# Check if package is installed
if gpkg list --installed --json | jq -e '.packages[] | select(.name=="my-tool")' > /dev/null; then
    echo "Package installed"
fi

# Get package info programmatically
VERSION=$(gpkg info owner/repo --json | jq -r '.version')
```

### Automated Upgrades

Set up automatic package upgrades:

```bash
#!/bin/bash
# auto-upgrade.sh

# Update indices
gpkg update --yes

# Upgrade all packages
gpkg upgrade --yes --dry-run > /tmp/upgrades.txt

# Review and apply
if [ -s /tmp/upgrades.txt ]; then
    gpkg upgrade --yes
    echo "Upgraded: $(cat /tmp/upgrades.txt)"
fi
```

### Batch Operations

Process multiple packages efficiently:

```bash
# Batch install with error handling
install_packages() {
    local failed=()
    for pkg in "$@"; do
        if ! gpkg install "$pkg" --yes; then
            failed+=("$pkg")
        fi
    done
    
    if [ ${#failed[@]} -gt 0 ]; then
        echo "Failed to install: ${failed[*]}"
        return 1
    fi
}

install_packages owner/pkg1 owner/pkg2 owner/pkg3
```

### CI/CD Integration

Integrate gpkg into your CI/CD pipeline:

```yaml
# .github/workflows/setup-tools.yml
name: Setup Tools
on: [push]

jobs:
  setup:
    runs-on: ubuntu-latest
    steps:
      - name: Install gpkg
        run: |
          curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg-linux-amd64.tar.gz
          tar -xzf gpkg-linux-amd64.tar.gz
          sudo mv gpkg-linux-amd64 /usr/local/bin/gpkg
          
      - name: Install tools
        run: |
          gpkg update --yes
          gpkg install --yes owner/tool1
          gpkg install --yes owner/tool2
          
      - name: Verify installations
        run: gpkg list --installed
```

## Custom Package Sources

### Create Your Own Package Index

Create a custom package source:

```json
{
  "name": "my-packages",
  "version": "1.0",
  "packages": [
    {
      "name": "my-tool",
      "repo": "owner/my-tool",
      "description": "My custom tool",
      "latest_version": "1.0.0",
      "manifest_url": "https://example.com/manifests/my-tool.yaml"
    }
  ]
}
```

Host it and add to gpkg:

```bash
gpkg add-source https://example.com/index.json --name my-packages
```

### Mirror Package Sources

Create a local mirror for offline use:

```bash
# Download package indices
gpkg update

# Mirror to local directory
mkdir -p ~/gpkg-mirror
cp -r ~/.gpkg/sources.d/* ~/gpkg-mirror/

# Use local mirror
gpkg add-source file://~/gpkg-mirror/index.json --priority 100
```

### Private Package Sources

Set up authentication for private sources:

```bash
# Add source with authentication (in secure script)
export SOURCE_TOKEN="your-token-here"
curl -H "Authorization: Bearer $SOURCE_TOKEN" \
     https://private.example.com/index.json \
     > ~/.gpkg/sources.d/private.json

gpkg add-source file://$HOME/.gpkg/sources.d/private.json --name private
```

## Working with Manifests

### Validate Manifests

Validate before installing:

```bash
# Validate local manifest
gpkg validate ./manifest.yaml

# Check for common issues
gpkg validate ./manifest.yaml --strict
```

### Generate Manifests

Create manifests from templates:

```bash
#!/bin/bash
# generate-manifest.sh

cat > manifest.yaml <<EOF
name: ${PACKAGE_NAME}
version: ${VERSION}
description: "${DESCRIPTION}"
homepage: "https://github.com/${OWNER}/${REPO}"
license: "${LICENSE}"

release_assets:
  - url: "https://github.com/${OWNER}/${REPO}/releases/download/${VERSION}/${PACKAGE_NAME}-linux-amd64.tar.gz"
    platform: "linux/amd64"
    checksum:
      sha256: "${CHECKSUM}"
EOF
```

### Multi-Platform Manifests

Support multiple platforms:

```yaml
release_assets:
  - url: "https://example.com/tool-linux-amd64.tar.gz"
    platform: "linux/amd64"
    checksum:
      sha256: "abc123..."
      
  - url: "https://example.com/tool-linux-arm64.tar.gz"
    platform: "linux/arm64"
    checksum:
      sha256: "def456..."
      
  - url: "https://example.com/tool-darwin-amd64.tar.gz"
    platform: "darwin/amd64"
    checksum:
      sha256: "ghi789..."
      
  - url: "https://example.com/tool-windows-amd64.zip"
    platform: "windows/amd64"
    checksum:
      sha256: "jkl012..."
```

## Dependency Management

### Manual Dependency Resolution

Install dependencies manually:

```bash
# Check dependencies
gpkg info package --deps-tree

# Install dependencies first
gpkg install dependency1
gpkg install dependency2
gpkg install package --no-deps
```

### Dependency Ordering

Ensure correct installation order:

```bash
#!/bin/bash
# install-with-deps.sh

declare -A installed

install_with_deps() {
    local pkg="$1"
    
    # Skip if already installed
    [[ ${installed[$pkg]} ]] && return
    
    # Get dependencies
    deps=$(gpkg info "$pkg" --json | jq -r '.dependencies[]')
    
    # Install dependencies first
    for dep in $deps; do
        install_with_deps "$dep"
    done
    
    # Install package
    gpkg install "$pkg" --yes
    installed[$pkg]=1
}

install_with_deps "my-package"
```

## Performance Optimization

### Increase Parallel Downloads

```toml
# ~/.gpkg/config.toml
parallel_downloads = 16
```

### Use Local Cache

```bash
# Set up local package cache
mkdir -p ~/.gpkg/cache

# Preserve downloaded archives
export GPKG_CACHE_DIR=~/.gpkg/cache
```

### Optimize Database

```bash
# Vacuum database periodically
sqlite3 ~/.gpkg/pkgdb.sqlite "VACUUM;"

# Analyze for query optimization
sqlite3 ~/.gpkg/pkgdb.sqlite "ANALYZE;"
```

## Security Best Practices

### Always Verify Checksums

```toml
# Never disable in production!
require_checksums = true
```

### Audit Package Sources

```bash
# Review configured sources
gpkg list-sources

# Remove untrusted sources
gpkg remove-source untrusted-source --yes
```

### Verify Package Integrity

```bash
# Check installed package integrity
gpkg verify package-name

# Verify all installed packages
gpkg verify --all
```

### Secure Configuration

```bash
# Restrict config file permissions
chmod 600 ~/.gpkg/config.toml

# Protect package database
chmod 600 ~/.gpkg/pkgdb.sqlite
```

## Advanced Troubleshooting

### Enable Debug Logging

```bash
# Verbose output
gpkg install package -vvv

# Debug log level
gpkg install package --log-level debug

# Log to file
gpkg install package --log-level debug 2>&1 | tee install.log
```

### Database Inspection

```bash
# Inspect database directly
sqlite3 ~/.gpkg/pkgdb.sqlite

# List installed packages
sqlite3 ~/.gpkg/pkgdb.sqlite "SELECT name, version FROM packages;"

# Check file tracking
sqlite3 ~/.gpkg/pkgdb.sqlite "SELECT * FROM files WHERE package='my-tool';"
```

### Network Debugging

```bash
# Test source connectivity
curl -I https://packages.example.com/index.json

# Download with verbose output
gpkg install package --verbose --log-level debug

# Use offline mode to test local operations
gpkg install ./local-manifest.yaml --offline
```

### Recovery Procedures

```bash
# Backup before risky operations
tar czf gpkg-backup.tar.gz ~/.gpkg

# Reset to clean state
rm -rf ~/.gpkg
mkdir -p ~/.gpkg

# Rebuild from package list
cat installed-packages.txt | while read pkg; do
    gpkg install "$pkg"
done
```

## See Also

- [Package Management](Package-Management) - Standard package workflows
- [Configuration](Configuration) - Configuration options
- [Troubleshooting](../Troubleshooting) - Common issues
- [Developer Guide](../Developer-Guide/Contributing) - Contributing to gpkg

---

**Next:** Learn about the [Developer Guide](../Developer-Guide/Contributing) if you want to contribute.
