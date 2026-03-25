# Getting Started with gpkg

Welcome to **gpkg**! This guide will help you get up and running with gpkg in just a few minutes.

## Table of Contents
- [What is gpkg?](#what-is-gpkg)
- [Installation](#installation)
- [First Steps](#first-steps)
- [Your First Package Installation](#your-first-package-installation)
- [Basic Commands](#basic-commands)
- [Next Steps](#next-steps)

---

## What is gpkg?

**gpkg** is a simple, user-focused package manager that installs either release binaries or builds from source into a configurable prefix. Think "pacman for GitHub releases + source builds" — designed to be secure, script-friendly, and easy to use.

### Key Features

- 📦 **Install from GitHub releases** or build from source
- 🔒 **Secure downloads** with checksum verification (SHA256/SHA512)
- 🎯 **Isolated installations** in configurable prefix (default: `~/.gpkg`)
- 🔄 **Package management** - install, upgrade, rollback, uninstall
- 📋 **Source repositories** - add multiple package sources
- 🔍 **Package discovery** - search and view package information
- 🤖 **CI-friendly** - JSON output, non-interactive mode, dry-run support

---

## Installation

### Quick Install (Recommended)

For detailed installation instructions, see the [Installation Guide](Installation.md).

#### Linux/macOS (from release)

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

**Requirements:** Go 1.20+, SQLite3, Git

```bash
git clone https://github.com/grave0x/gpkg.git
cd gpkg
go build -o bin/gpkg .
sudo mv bin/gpkg /usr/local/bin/
gpkg --version
```

---

## First Steps

### 1. Initialize gpkg

After installing gpkg, it will automatically create its configuration directory on first run. By default, gpkg uses `~/.gpkg` for all its data:

```bash
# View your current configuration
gpkg config show
```

The default configuration includes:
- **Install prefix**: `~/.gpkg` (where packages are installed)
- **Package database**: `~/.gpkg/pkgdb.sqlite` (SQLite database tracking installed packages)
- **Sources directory**: `~/.gpkg/sources.d` (package source definitions)

### 2. Add Your PATH

To use installed packages, add gpkg's bin directory to your PATH. Add this line to your `~/.bashrc`, `~/.zshrc`, or equivalent:

```bash
export PATH="$HOME/.gpkg/bin:$PATH"
```

Then reload your shell:

```bash
source ~/.bashrc  # or source ~/.zshrc
```

### 3. Add a Package Source (Optional)

gpkg can install packages from GitHub releases directly, or you can add package repositories for curated collections:

```bash
# Add a package source
gpkg add-source https://packages.example.com/index.json --name example

# Update package indices
gpkg update
```

For now, you can skip this step and install packages directly from GitHub repositories.

---

## Your First Package Installation

Let's install a package! Here are a few common scenarios:

### Install from GitHub Release

Install a package directly from a GitHub repository (gpkg will use the latest release):

```bash
# Install from a GitHub repository
gpkg install owner/repo

# Example: Install a specific tool
gpkg install cli/cli
```

gpkg will:
1. 📥 Fetch the latest release from GitHub
2. 🔍 Select the correct binary for your platform (linux/amd64, darwin/arm64, etc.)
3. ✅ Verify checksums for security
4. 📦 Extract and install to `~/.gpkg/`
5. 💾 Record the installation in its database

### Install from a Local Manifest

If you have a package manifest file (YAML format):

```bash
# Install from a local manifest
gpkg install ./path/to/manifest.yaml

# Example with the included example
gpkg install ./examples/manifest.yaml
```

### Install with Source Build

Force gpkg to build from source instead of using a pre-built release:

```bash
gpkg install owner/repo --from-source
```

This is useful when:
- No pre-built binary is available for your platform
- You want to compile with custom flags
- You prefer to build from source for security reasons

### Preview Before Installing (Dry Run)

Want to see what will happen without actually installing?

```bash
gpkg install owner/repo --dry-run
```

This will show you the installation plan without making any changes.

---

## Basic Commands

Here's a quick reference of the most common gpkg commands:

### Package Installation & Management

```bash
# Install a package
gpkg install owner/repo

# List installed packages
gpkg list --installed

# List available packages (from sources)
gpkg list --available

# Upgrade all packages
gpkg upgrade

# Upgrade specific packages
gpkg upgrade package1 package2

# Uninstall a package
gpkg uninstall package-name
```

### Package Information & Discovery

```bash
# Search for packages
gpkg search tool

# View detailed package information
gpkg info owner/repo

# Validate a manifest file
gpkg validate ./manifest.yaml
```

### Source Management

```bash
# Add a package source
gpkg add-source https://packages.example.com/index.json

# List configured sources
gpkg list-sources

# Remove a source
gpkg remove-source example

# Update package indices from all sources
gpkg update
```

### Configuration

```bash
# Show current configuration
gpkg config show

# Get a specific config value
gpkg config get install_prefix

# Set a config value
gpkg config set require_checksums true
```

### Helpful Flags

All commands support these useful flags:

```bash
--dry-run          # Preview actions without making changes
--json             # Output in JSON format (great for scripts)
-y, --yes          # Skip confirmations (non-interactive mode)
--verbose, -v      # Increase output detail (use -vv or -vvv for more)
--quiet            # Minimal output
--help, -h         # Show help for any command
```

**Example:**

```bash
# Preview an upgrade with detailed output
gpkg upgrade --dry-run --verbose

# Install without prompts (for scripts)
gpkg install owner/repo --yes

# Get machine-readable output
gpkg list --installed --json
```

---

## Next Steps

### Learn More

- **[User Guide](User-Guide/README.md)** - Comprehensive guide to all gpkg features
- **[Reference](Reference/README.md)** - Complete command reference and manifest format
- **[Developer Guide](Developer-Guide/README.md)** - Contributing and development information

### Common Tasks

- **[Writing Package Manifests](Reference/Manifest-Format.md)** - Create your own package definitions
- **[Configuration Guide](User-Guide/Configuration.md)** - Customize gpkg to your needs
- **[CI/CD Integration](User-Guide/CI-CD.md)** - Use gpkg in automated workflows
- **[Troubleshooting](User-Guide/Troubleshooting.md)** - Common issues and solutions

### Get Involved

- 🐛 **Found a bug?** [Report it on GitHub](https://github.com/grave0x/gpkg/issues)
- 💡 **Have an idea?** [Start a discussion](https://github.com/grave0x/gpkg/discussions)
- 🤝 **Want to contribute?** Check out [CONTRIBUTING.md](../CONTRIBUTING.md)

### Example Workflows

#### Set up a development machine

```bash
# Add package sources
gpkg add-source https://packages.example.com/dev-tools.json

# Update indices
gpkg update

# Install your favorite tools
gpkg install cli/cli
gpkg install sharkdp/fd
gpkg install BurntSushi/ripgrep

# Verify installations
gpkg list --installed
```

#### Keep your tools up to date

```bash
# Update package indices
gpkg update

# Preview what will be upgraded
gpkg upgrade --dry-run

# Upgrade everything
gpkg upgrade -y
```

#### Create and install a custom package

```bash
# Create a manifest file (see examples/manifest.yaml)
vim my-tool.yaml

# Validate the manifest
gpkg validate my-tool.yaml

# Install it
gpkg install ./my-tool.yaml
```

---

## Quick Tips

💡 **Tip #1**: Use `--dry-run` before making changes to preview what will happen.

💡 **Tip #2**: gpkg stores everything in `~/.gpkg` by default - it won't touch your system directories.

💡 **Tip #3**: All downloads are verified with checksums for security. If a checksum fails, gpkg will refuse to install.

💡 **Tip #4**: Use `--json` flag with any command for machine-readable output - perfect for scripts!

💡 **Tip #5**: Tab completion coming soon! Star the repo to stay updated.

---

**Ready to dive deeper?** Check out the [User Guide](User-Guide/README.md) for advanced features and workflows!

---

[← Back to Wiki Home](Home.md) | [Installation Guide →](Installation.md) | [User Guide →](User-Guide/README.md)
