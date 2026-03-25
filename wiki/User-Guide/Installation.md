# Installation Guide

This guide provides detailed instructions for installing **gpkg** on Linux, macOS, and Windows systems.

## Table of Contents

- [System Requirements](#system-requirements)
- [Installation from Release Binaries](#installation-from-release-binaries)
  - [Linux](#linux)
  - [macOS](#macos)
  - [Windows](#windows)
- [Installation from Source](#installation-from-source)
- [Verification Steps](#verification-steps)
- [Post-Installation Setup](#post-installation-setup)
- [Troubleshooting](#troubleshooting)
- [Uninstallation](#uninstallation)

---

## System Requirements

### Minimum Requirements

- **Operating System**: Linux (amd64/arm64), macOS (10.15+), Windows 10/11
- **Disk Space**: 20 MB for binary installation, 100+ MB for source installation
- **Network**: Internet connection for downloading packages (offline mode available after initial setup)

### Runtime Dependencies

- **Git**: Required for source builds and Git-based package installations
  - Version: 2.20 or higher recommended
- **SQLite3**: Built into the binary (no separate installation needed for release binaries)

### Build Dependencies (Source Installation Only)

- **Go**: Version 1.20 or higher
- **GCC/Clang**: C compiler for CGO (SQLite3 support)
- **SQLite3 Development Libraries**: Required for building with CGO
- **Git**: For cloning the repository

---

## Installation from Release Binaries

Installing from pre-built release binaries is the **recommended method** for most users. Binaries are available for common platforms and architectures.

### Linux

#### Method 1: Automated Installation Script (Recommended)

```bash
# Download and run the install script
curl -fsSL https://github.com/grave0x/gpkg/releases/latest/download/install.sh | bash

# Or download, inspect, then run
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/install.sh
chmod +x install.sh
./install.sh
```

The script automatically detects your platform and installs gpkg to `/usr/local/bin`.

#### Method 2: Manual Installation

**For AMD64 (x86_64):**

```bash
# Download the latest release
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg-linux-amd64.tar.gz

# Verify the checksum (optional but recommended)
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/checksums.txt
sha256sum --check --ignore-missing checksums.txt

# Extract the archive
tar -xzf gpkg-linux-amd64.tar.gz

# Install to system location
sudo install -m 755 gpkg /usr/local/bin/gpkg

# Or install to user location (no sudo required)
mkdir -p ~/.local/bin
install -m 755 gpkg ~/.local/bin/gpkg
```

**For ARM64 (aarch64):**

```bash
# Download the ARM64 release
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg-linux-arm64.tar.gz

# Extract and install
tar -xzf gpkg-linux-arm64.tar.gz
sudo install -m 755 gpkg /usr/local/bin/gpkg
```

#### Method 3: Install via Package Manager

**For Arch Linux (AUR):**

```bash
# Using yay
yay -S gpkg-bin

# Or using paru
paru -S gpkg-bin
```

**For Debian/Ubuntu (via .deb package):**

```bash
# Download the .deb package
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg_amd64.deb

# Install
sudo dpkg -i gpkg_amd64.deb

# Fix dependencies if needed
sudo apt-get install -f
```

**For Fedora/RHEL (via .rpm package):**

```bash
# Download the .rpm package
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg_amd64.rpm

# Install
sudo rpm -ivh gpkg_amd64.rpm

# Or using dnf
sudo dnf install gpkg_amd64.rpm
```

### macOS

#### Method 1: Using Homebrew (Recommended)

```bash
# Add the tap (once)
brew tap grave0x/gpkg

# Install gpkg
brew install gpkg

# Update gpkg
brew upgrade gpkg
```

#### Method 2: Manual Installation

**For Intel Macs (x86_64):**

```bash
# Download the latest release
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg-darwin-amd64.tar.gz

# Extract
tar -xzf gpkg-darwin-amd64.tar.gz

# Install
sudo install -m 755 gpkg /usr/local/bin/gpkg

# Or install to user location
mkdir -p ~/.local/bin
install -m 755 gpkg ~/.local/bin/gpkg
```

**For Apple Silicon Macs (ARM64):**

```bash
# Download the ARM64 release
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg-darwin-arm64.tar.gz

# Extract and install
tar -xzf gpkg-darwin-arm64.tar.gz
sudo install -m 755 gpkg /usr/local/bin/gpkg
```

**macOS Security Note:**

On first run, macOS may block the binary because it's not signed. To allow it:

```bash
# Remove the quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/gpkg

# Or allow it through System Preferences:
# System Preferences > Security & Privacy > General > Click "Allow Anyway"
```

### Windows

#### Method 1: Using Scoop (Recommended)

```powershell
# Add the bucket
scoop bucket add gpkg https://github.com/grave0x/scoop-gpkg

# Install gpkg
scoop install gpkg

# Update gpkg
scoop update gpkg
```

#### Method 2: Using Chocolatey

```powershell
# Install
choco install gpkg

# Update
choco upgrade gpkg
```

#### Method 3: Manual Installation

**For AMD64 (64-bit Windows):**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/grave0x/gpkg/releases/latest/download/gpkg-windows-amd64.zip" -OutFile "gpkg-windows-amd64.zip"

# Extract (using PowerShell 5.0+)
Expand-Archive -Path gpkg-windows-amd64.zip -DestinationPath .

# Move to a directory in your PATH
# Option 1: System-wide (requires admin)
Move-Item gpkg.exe "C:\Program Files\gpkg\gpkg.exe"

# Option 2: User-level
$userBin = "$env:LOCALAPPDATA\Programs\gpkg"
New-Item -ItemType Directory -Force -Path $userBin
Move-Item gpkg.exe "$userBin\gpkg.exe"

# Add to PATH (user-level)
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$userBin",
    "User"
)
```

**Using Command Prompt:**

```cmd
REM Download using curl (Windows 10+)
curl -LO https://github.com/grave0x/gpkg/releases/latest/download/gpkg-windows-amd64.zip

REM Extract using tar (Windows 10+)
tar -xf gpkg-windows-amd64.zip

REM Move to Program Files (requires admin cmd)
move gpkg.exe "C:\Program Files\gpkg\gpkg.exe"
```

---

## Installation from Source

Building from source gives you the latest features and allows customization.

### Prerequisites

#### Linux

**Debian/Ubuntu:**

```bash
sudo apt-get update
sudo apt-get install -y git gcc make sqlite3 libsqlite3-dev golang-go
```

**Fedora/RHEL:**

```bash
sudo dnf install -y git gcc make sqlite sqlite-devel golang
```

**Arch Linux:**

```bash
sudo pacman -S git gcc make sqlite go
```

**Alpine Linux:**

```bash
sudo apk add git gcc make musl-dev sqlite sqlite-dev go
```

#### macOS

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Go using Homebrew
brew install go git

# SQLite is pre-installed on macOS
```

#### Windows

1. Install [Git for Windows](https://git-scm.com/download/win)
2. Install [Go](https://golang.org/dl/) (version 1.20+)
3. Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)
4. Ensure `gcc` and `go` are in your PATH

### Build Instructions

#### Step 1: Clone the Repository

```bash
git clone https://github.com/grave0x/gpkg.git
cd gpkg
```

#### Step 2: Install Go Dependencies

```bash
go mod download
```

#### Step 3: Build the Binary

**Standard Build:**

```bash
# Build to bin/gpkg
go build -o bin/gpkg .

# Or use the build script
./scripts/build.sh
```

**Optimized Build (Smaller Binary):**

```bash
# Build with optimizations
go build -ldflags="-s -w" -o bin/gpkg .

# Optional: Compress with UPX (if installed)
upx --best --lzma bin/gpkg
```

**Static Build (No CGO Dependencies):**

```bash
# Build a fully static binary (note: some features may be limited)
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/gpkg .
```

**With Version Information:**

```bash
# Build with version info
VERSION=$(git describe --tags --always --dirty)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
go build -ldflags="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" -o bin/gpkg .
```

#### Step 4: Install the Binary

**Linux/macOS:**

```bash
# System-wide installation (requires sudo)
sudo install -m 755 bin/gpkg /usr/local/bin/gpkg

# User installation
mkdir -p ~/.local/bin
install -m 755 bin/gpkg ~/.local/bin/gpkg

# Make sure ~/.local/bin is in your PATH
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**Windows:**

```powershell
# Create directory
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\Programs\gpkg"

# Copy binary
Copy-Item bin\gpkg.exe "$env:LOCALAPPDATA\Programs\gpkg\gpkg.exe"

# Add to PATH
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:LOCALAPPDATA\Programs\gpkg",
    "User"
)
```

#### Step 5: Verify the Build

```bash
gpkg --version
```

### Building for Multiple Platforms

Use the provided build script to create binaries for all supported platforms:

```bash
# Build for all platforms
./scripts/build.sh

# Specify version
./scripts/build.sh 1.0.0

# Binaries will be in bin/ directory:
# - gpkg-linux-amd64
# - gpkg-linux-arm64
# - gpkg-darwin-amd64
# - gpkg-darwin-arm64
# - gpkg-windows-amd64.exe
```

---

## Verification Steps

After installation, verify that gpkg is working correctly:

### Step 1: Check Version

```bash
gpkg --version
```

**Expected Output:**
```
gpkg version 1.0.0 (commit: abc1234) built at 2024-01-15_10:30:00
```

### Step 2: Check Help

```bash
gpkg --help
```

**Expected Output:** A list of available commands and global flags.

### Step 3: Verify Installation Path

```bash
# Linux/macOS
which gpkg

# Windows (PowerShell)
Get-Command gpkg | Select-Object -ExpandProperty Source

# Windows (Command Prompt)
where gpkg
```

### Step 4: Test Basic Functionality

```bash
# Check configuration
gpkg config show

# List sources (should be empty initially)
gpkg list-sources
```

### Step 5: Verify Database Creation

```bash
# Install a test package (optional)
gpkg install --dry-run grave0x/example-package

# Check that the database was created
ls -lh ~/.gpkg/pkgdb.sqlite  # Linux/macOS
dir %USERPROFILE%\.gpkg\pkgdb.sqlite  # Windows
```

---

## Post-Installation Setup

### Configure PATH (If Not Already Done)

**Linux/macOS (Bash):**

```bash
# Add to ~/.bashrc or ~/.bash_profile
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**Linux/macOS (Zsh):**

```bash
# Add to ~/.zshrc
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Linux/macOS (Fish):**

```bash
# Add to ~/.config/fish/config.fish
echo 'set -gx PATH $HOME/.local/bin $PATH' >> ~/.config/fish/config.fish
```

**Windows:**

PATH should be automatically configured by the installer or manual steps above.

### Initialize Configuration

```bash
# Create default config file
gpkg config show > ~/.gpkg/config.toml

# Or set custom prefix
gpkg config set install_prefix ~/my-packages
```

### Add Package Sources

```bash
# Add the official package source
gpkg add-source https://packages.gpkg.dev/index.json

# Update package indices
gpkg update
```

### Enable Shell Completions (Optional)

**Bash:**

```bash
gpkg completion bash > /etc/bash_completion.d/gpkg  # System-wide
# Or
gpkg completion bash > ~/.local/share/bash-completion/completions/gpkg  # User
```

**Zsh:**

```bash
gpkg completion zsh > "${fpath[1]}/_gpkg"
```

**Fish:**

```bash
gpkg completion fish > ~/.config/fish/completions/gpkg.fish
```

**PowerShell:**

```powershell
gpkg completion powershell | Out-File -Encoding UTF8 $PROFILE.CurrentUserAllHosts
```

---

## Troubleshooting

### Common Installation Issues

#### Issue 1: "command not found: gpkg"

**Cause:** The gpkg binary is not in your PATH.

**Solution:**

```bash
# Find where gpkg is installed
find / -name gpkg 2>/dev/null

# Add the directory to PATH
export PATH="/path/to/gpkg/directory:$PATH"

# Make permanent by adding to ~/.bashrc or ~/.zshrc
echo 'export PATH="/path/to/gpkg/directory:$PATH"' >> ~/.bashrc
```

#### Issue 2: "permission denied" when running gpkg

**Cause:** The binary is not executable.

**Solution:**

```bash
chmod +x /path/to/gpkg
```

#### Issue 3: Build fails with "sqlite3.h: No such file or directory"

**Cause:** SQLite3 development headers are not installed.

**Solution:**

```bash
# Debian/Ubuntu
sudo apt-get install libsqlite3-dev

# Fedora/RHEL
sudo dnf install sqlite-devel

# macOS
brew install sqlite

# Or build without CGO (limited functionality)
CGO_ENABLED=0 go build -o bin/gpkg .
```

#### Issue 4: macOS "cannot be opened because the developer cannot be verified"

**Cause:** macOS Gatekeeper blocks unsigned binaries.

**Solution:**

```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/gpkg

# Or allow through System Preferences
# System Preferences > Security & Privacy > General > "Allow Anyway"
```

#### Issue 5: Windows Defender or antivirus blocks gpkg

**Cause:** Some antivirus software flags new executables.

**Solution:**

1. Add an exception for gpkg in your antivirus settings
2. Download from official GitHub releases only
3. Verify checksums to ensure binary integrity

#### Issue 6: "go: go.mod file not found" during build

**Cause:** Not in the gpkg repository directory.

**Solution:**

```bash
# Make sure you're in the gpkg directory
cd /path/to/gpkg

# Verify go.mod exists
ls go.mod
```

#### Issue 7: Database initialization fails

**Cause:** Insufficient permissions or corrupted database.

**Solution:**

```bash
# Remove and reinitialize database
rm -f ~/.gpkg/pkgdb.sqlite

# Run gpkg again to recreate
gpkg config show
```

#### Issue 8: TLS/SSL certificate errors on older systems

**Cause:** Outdated CA certificates.

**Solution:**

```bash
# Update CA certificates
# Debian/Ubuntu
sudo apt-get update && sudo apt-get install ca-certificates

# Fedora/RHEL
sudo dnf update ca-certificates

# macOS
brew install ca-certificates
```

### Getting Help

If you encounter issues not covered here:

1. **Check existing issues**: [GitHub Issues](https://github.com/grave0x/gpkg/issues)
2. **Enable debug logging**: `gpkg --log-level debug <command>`
3. **Run with verbose output**: `gpkg -vvv <command>`
4. **Check system logs**: `journalctl -u gpkg` (Linux with systemd)
5. **Create a new issue**: Include OS, gpkg version, and error messages

---

## Uninstallation

### Remove Installed Packages (Optional)

Before uninstalling gpkg, you may want to remove packages it has installed:

```bash
# List installed packages
gpkg list --installed

# Uninstall all packages
gpkg list --installed --json | jq -r '.[] | .name' | xargs -I {} gpkg uninstall {}
```

### Uninstall gpkg Binary

**Linux/macOS (Manual Installation):**

```bash
# Remove binary
sudo rm -f /usr/local/bin/gpkg

# Or if installed to user directory
rm -f ~/.local/bin/gpkg
```

**macOS (Homebrew):**

```bash
brew uninstall gpkg
brew untap grave0x/gpkg
```

**Linux (Package Manager):**

```bash
# Debian/Ubuntu
sudo dpkg -r gpkg

# Fedora/RHEL
sudo rpm -e gpkg

# Arch Linux
yay -R gpkg-bin
# or
paru -R gpkg-bin
```

**Windows (Scoop):**

```powershell
scoop uninstall gpkg
```

**Windows (Chocolatey):**

```powershell
choco uninstall gpkg
```

**Windows (Manual):**

```powershell
# Remove binary
Remove-Item "$env:LOCALAPPDATA\Programs\gpkg\gpkg.exe"

# Remove from PATH
# Go to: System Properties > Environment Variables > Edit PATH
# Remove the gpkg entry
```

### Remove Configuration and Data

**Complete Removal (Including All Data):**

```bash
# Linux/macOS
rm -rf ~/.gpkg

# Windows (PowerShell)
Remove-Item -Recurse -Force "$env:USERPROFILE\.gpkg"
```

**Selective Removal:**

```bash
# Remove only the database
rm ~/.gpkg/pkgdb.sqlite

# Remove only installed packages (keep config)
rm -rf ~/.gpkg/bin ~/.gpkg/lib ~/.gpkg/share

# Remove only package sources
rm -rf ~/.gpkg/sources.d
```

### Verify Uninstallation

```bash
# Verify binary is removed
which gpkg  # Should return nothing

# Verify data is removed
ls ~/.gpkg  # Should not exist
```

---

## Next Steps

After successful installation:

1. **Read the [User Guide](../User-Guide.md)** to learn how to use gpkg
2. **Configure package sources**: `gpkg add-source <url>`
3. **Search for packages**: `gpkg search <term>`
4. **Install your first package**: `gpkg install <package>`
5. **Explore the [CLI Reference](../CLI-Reference.md)** for advanced usage

---

## Additional Resources

- **Official Documentation**: [GitHub Wiki](https://github.com/grave0x/gpkg/wiki)
- **Source Code**: [GitHub Repository](https://github.com/grave0x/gpkg)
- **Issue Tracker**: [GitHub Issues](https://github.com/grave0x/gpkg/issues)
- **Releases**: [GitHub Releases](https://github.com/grave0x/gpkg/releases)
- **Build Scripts**: [`scripts/build.sh`](../../scripts/build.sh)
- **Example Manifests**: [`examples/`](../../examples/)

---

**Last Updated**: 2024-01-15  
**gpkg Version**: 1.0.0
