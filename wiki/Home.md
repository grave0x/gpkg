# Welcome to the gpkg Wiki!

**gpkg** is a simple, user-focused package manager that installs either release binaries or builds from source into a configurable prefix. Think "pacman for GitHub releases + source builds."

## 🚀 Quick Navigation

### Getting Started
- **[Getting Started Guide](Getting-Started)** - Installation and first steps
- **[Installation](User-Guide/Installation)** - Detailed installation instructions for all platforms

### User Documentation
- **[Basic Commands](User-Guide/Basic-Commands)** - Learn the essential gpkg commands
- **[Configuration](User-Guide/Configuration)** - Set up and customize gpkg
- **[Package Management](User-Guide/Package-Management)** - Installing, upgrading, and managing packages
- **[Advanced Usage](User-Guide/Advanced-Usage)** - Power user features and tips

### Developer Documentation
- **[Contributing](Developer-Guide/Contributing)** - How to contribute to gpkg
- **[Architecture](Developer-Guide/Architecture)** - Understanding the codebase structure
- **[Development Setup](Developer-Guide/Development-Setup)** - Set up your development environment
- **[Testing](Developer-Guide/Testing)** - Testing guidelines and best practices
- **[Release Process](Developer-Guide/Release-Process)** - How releases are created

### Reference Documentation
- **[CLI Reference](Reference/CLI-Reference)** - Complete command-line reference
- **[API Reference](Reference/API-Reference)** - Internal API documentation
- **[Manifest Format](Reference/Manifest-Format)** - Package manifest specification
- **[Configuration Reference](Reference/Configuration-Reference)** - All configuration options

### Help & Support
- **[Troubleshooting](Troubleshooting)** - Common issues and solutions
- **[GitHub Issues](https://github.com/grave0x/gpkg/issues)** - Report bugs or request features
- **[Main Repository](https://github.com/grave0x/gpkg)** - View source code

## 📦 What is gpkg?

gpkg is designed to make package management simple and reliable:

- **✅ Install from GitHub releases** - Download and install pre-built binaries
- **🔨 Build from source** - Compile packages from source with custom build steps
- **🔒 Secure** - Checksum verification for all downloads (SHA256/SHA512)
- **🎯 Isolated** - Packages install to `~/.gpkg` by default, separate from system
- **📊 Tracked** - SQLite database tracks all installed packages and files
- **🤖 CI-friendly** - JSON output, non-interactive mode, dry-run support

## 🎯 Quick Start Example

```bash
# Add a package source
gpkg add-source https://packages.example.com/index.json

# Update package indices
gpkg update

# Search for a package
gpkg search tool

# Install a package
gpkg install owner/repo

# List installed packages
gpkg list --installed

# Upgrade all packages
gpkg upgrade
```

## 📚 Documentation Sections

### For Users
Start with the [Getting Started Guide](Getting-Started) to install gpkg and learn the basics. Then explore the [User Guide](User-Guide/Installation) for detailed usage instructions.

### For Developers
Read the [Contributing Guide](Developer-Guide/Contributing) to learn how to contribute, then check out the [Architecture Guide](Developer-Guide/Architecture) to understand the codebase.

### For Reference
Use the [CLI Reference](Reference/CLI-Reference) for complete command documentation and the [Configuration Reference](Reference/Configuration-Reference) for all config options.

## 🔗 External Links

- **Repository**: [https://github.com/grave0x/gpkg](https://github.com/grave0x/gpkg)
- **Issues**: [Report a bug or request a feature](https://github.com/grave0x/gpkg/issues)
- **Releases**: [Download latest release](https://github.com/grave0x/gpkg/releases)

## 📝 Recent Updates

Check the [Release Notes](https://github.com/grave0x/gpkg/releases) for the latest changes and improvements.

---

**Need help?** Check the [Troubleshooting](Troubleshooting) page or [open an issue](https://github.com/grave0x/gpkg/issues).
