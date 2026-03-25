# Troubleshooting Guide

This guide covers common issues you may encounter when using gpkg and their solutions.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Package Installation Failures](#package-installation-failures)
- [Configuration Issues](#configuration-issues)
- [Permission Errors](#permission-errors)
- [Database Issues](#database-issues)
- [Common Error Messages](#common-error-messages)
- [Debug Mode](#debug-mode)
- [Reporting Issues](#reporting-issues)

---

## Installation Issues

### Problem: "Go version too old" or Build Fails with Go Version Error

**Symptoms:**
```
go: go.mod requires go >= 1.20 (running go 1.19)
```

**Solution:**
1. Check your Go version:
   ```bash
   go version
   ```

2. Upgrade Go to 1.20 or higher:
   - **Ubuntu/Debian**: Download from [golang.org](https://golang.org/dl/)
   - **macOS**: `brew upgrade go`
   - **Arch Linux**: `sudo pacman -S go`

3. Verify the upgrade:
   ```bash
   go version
   ```

**Prevention:**
- Keep Go updated to the latest stable version
- Check project requirements before building

---

### Problem: "sqlite3.h: No such file or directory" During Build

**Symptoms:**
```
# github.com/mattn/go-sqlite3
In file included from sqlite3-binding.c:1:
./sqlite3.h:8:10: fatal error: sqlite3.h: No such file or directory
```

**Solution:**
Install SQLite development libraries:

- **Ubuntu/Debian**:
  ```bash
  sudo apt-get update
  sudo apt-get install libsqlite3-dev gcc
  ```

- **Fedora/RHEL/CentOS**:
  ```bash
  sudo dnf install sqlite-devel gcc
  ```

- **Arch Linux**:
  ```bash
  sudo pacman -S sqlite gcc
  ```

- **macOS**:
  ```bash
  brew install sqlite3
  ```

Then rebuild:
```bash
go clean -cache
go build -o bin/gpkg .
```

**Prevention:**
- Install build dependencies before attempting to build from source
- See [BUILD.md](../BUILD.md) for complete dependency list

---

### Problem: Build Succeeds But Binary Doesn't Run

**Symptoms:**
```
./gpkg: command not found
# or
bash: /usr/local/bin/gpkg: Permission denied
```

**Solution:**
1. Check if the binary exists:
   ```bash
   ls -lh bin/gpkg
   ```

2. Make it executable:
   ```bash
   chmod +x bin/gpkg
   ```

3. Verify it's in your PATH:
   ```bash
   which gpkg
   ```

4. If not in PATH, either:
   - Move to `/usr/local/bin`: `sudo mv bin/gpkg /usr/local/bin/`
   - Or add to PATH: `export PATH="$HOME/.local/bin:$PATH"`

**Prevention:**
- Follow installation instructions carefully
- Add custom binary locations to your `~/.bashrc` or `~/.zshrc`

---

## Package Installation Failures

### Problem: Checksum Verification Failed

**Symptoms:**
```
Error: Checksum validation failed
exit code: 4
Expected: abc123...
Got: def456...
```

**Solution:**
1. **Corrupted download** - Retry the installation:
   ```bash
   gpkg install owner/repo --force
   ```

2. **Network interference** - If behind a proxy or using a VPN, try:
   ```bash
   # Disable proxy temporarily
   unset HTTP_PROXY HTTPS_PROXY
   gpkg install owner/repo
   ```

3. **Outdated manifest** - Update package sources:
   ```bash
   gpkg update
   gpkg install owner/repo
   ```

4. **Debug checksum issue**:
   ```bash
   gpkg install owner/repo --verbose --log-level debug
   ```

**When to worry:**
- If retrying doesn't fix it, the release asset may have been modified
- Contact package maintainer or report to gpkg issues
- **Never** disable checksum validation (`require_checksums = false`) unless absolutely necessary

**Prevention:**
- Use stable internet connection
- Keep package sources updated with `gpkg update`
- Use `--dry-run` to preview downloads before installing

---

### Problem: Network Timeout or Connection Refused

**Symptoms:**
```
Error: Network error
exit code: 3
Failed to download: connection timeout
```

**Solution:**
1. **Check internet connection**:
   ```bash
   curl -I https://github.com
   ```

2. **Increase timeout** in config (`~/.gpkg/config.toml`):
   ```toml
   download_timeout_seconds = 300
   ```

3. **Use offline mode** if you have cached packages:
   ```bash
   gpkg install package-name --offline
   ```

4. **Check proxy settings**:
   ```bash
   echo $HTTP_PROXY $HTTPS_PROXY
   # If needed, set them:
   export HTTP_PROXY=http://proxy.example.com:8080
   export HTTPS_PROXY=http://proxy.example.com:8080
   ```

5. **Retry with verbose output**:
   ```bash
   gpkg install owner/repo -vv
   ```

**Prevention:**
- Configure appropriate timeouts for your network
- Use `--dry-run` to validate before large downloads
- Consider caching frequently-used packages

---

### Problem: "Package Not Found"

**Symptoms:**
```
Error: Package not found
exit code: 7
Package: owner/repo
```

**Solution:**
1. **Update package sources**:
   ```bash
   gpkg update
   gpkg search repo
   ```

2. **Check if package exists**:
   ```bash
   gpkg info owner/repo
   ```

3. **Verify source is added**:
   ```bash
   gpkg list-sources
   ```

4. **Add the package source** if missing:
   ```bash
   gpkg add-source https://packages.example.com/index.json
   gpkg update
   ```

5. **Install from local manifest** as workaround:
   ```bash
   # Create manifest.yaml for the package
   gpkg install ./manifest.yaml
   ```

**Prevention:**
- Keep sources updated with `gpkg update`
- Verify package name spelling
- Check package availability on GitHub

---

### Problem: Installation Fails During Build From Source

**Symptoms:**
```
Error: Installation failed
exit code: 5
Build step 'make' failed with exit code 2
```

**Solution:**
1. **Check build logs** with verbose mode:
   ```bash
   gpkg install owner/repo --from-source --verbose
   ```

2. **Install build dependencies**:
   - Check the project's README for requirements
   - Common tools: `gcc`, `make`, `cmake`, `pkg-config`

3. **Verify build environment**:
   ```bash
   which gcc make
   gcc --version
   ```

4. **Try with different prefix**:
   ```bash
   gpkg install owner/repo --from-source --prefix=$HOME/.local
   ```

5. **Manual build for debugging**:
   ```bash
   git clone https://github.com/owner/repo
   cd repo
   # Try build steps manually
   make
   ```

**Prevention:**
- Prefer release binaries over source builds when available
- Keep build tools updated
- Review build requirements in package manifest

---

## Configuration Issues

### Problem: Config File Not Found or Ignored

**Symptoms:**
```
Using default configuration
Config file not found: ~/.gpkg/config.toml
```

**Solution:**
1. **Create config directory**:
   ```bash
   mkdir -p ~/.gpkg
   ```

2. **Copy example config**:
   ```bash
   cp examples/config.example.toml ~/.gpkg/config.toml
   ```

3. **Verify config location**:
   ```bash
   gpkg config show
   ```

4. **Use explicit config path**:
   ```bash
   gpkg --config /path/to/config.toml install package
   ```

**Config precedence** (highest to lowest):
1. `--config <path>` flag
2. Environment variables (`GPKG_INSTALL_PREFIX`, etc.)
3. `~/.gpkg/config.toml` (user config)
4. `/etc/gpkg/config.toml` (system config)

**Prevention:**
- Create config file during initial setup
- Use `gpkg config show` to verify current configuration

---

### Problem: Invalid Config Syntax

**Symptoms:**
```
Error: Failed to parse config
Near line 5: invalid TOML syntax
```

**Solution:**
1. **Validate TOML syntax**:
   ```bash
   # Install TOML validator
   pip install toml

   # Validate
   python -c "import toml; toml.load(open('~/.gpkg/config.toml'))"
   ```

2. **Common TOML mistakes**:
   ```toml
   # ❌ Wrong - missing quotes
   install_prefix = ~/.gpkg

   # ✅ Correct
   install_prefix = "~/.gpkg"

   # ❌ Wrong - invalid boolean
   require_checksums = yes

   # ✅ Correct
   require_checksums = true
   ```

3. **Reset to default config**:
   ```bash
   mv ~/.gpkg/config.toml ~/.gpkg/config.toml.backup
   cp examples/config.example.toml ~/.gpkg/config.toml
   ```

**Prevention:**
- Use the example config as a template
- Validate syntax after editing
- Keep backups of working configurations

---

### Problem: Environment Variables Not Working

**Symptoms:**
```bash
export GPKG_INSTALL_PREFIX=/opt/gpkg
gpkg install package
# Still installs to ~/.gpkg
```

**Solution:**
1. **Check variable is set**:
   ```bash
   echo $GPKG_INSTALL_PREFIX
   ```

2. **Verify supported variables**:
   - `GPKG_INSTALL_PREFIX` - Install directory
   - `GPKG_PKGDB_PATH` - Database location
   - `GPKG_CONFIG` - Config file path
   - `LOG_LEVEL` - Logging verbosity

3. **Use CLI flag instead** (highest precedence):
   ```bash
   gpkg install package --prefix=/opt/gpkg
   ```

4. **Export in shell profile** for persistence:
   ```bash
   echo 'export GPKG_INSTALL_PREFIX=/opt/gpkg' >> ~/.bashrc
   source ~/.bashrc
   ```

**Prevention:**
- Use config file for persistent settings
- Document custom environment variables
- Verify with `gpkg config show`

---

## Permission Errors

### Problem: "Permission Denied" When Installing

**Symptoms:**
```
Error: Installation failed
mkdir /opt/gpkg: permission denied
```

**Solution:**
1. **Use user-writable prefix** (recommended):
   ```bash
   gpkg install package --prefix=$HOME/.gpkg
   ```

2. **Or grant permissions to target directory**:
   ```bash
   sudo mkdir -p /opt/gpkg
   sudo chown $USER:$USER /opt/gpkg
   gpkg install package --prefix=/opt/gpkg
   ```

3. **Never run gpkg as root** unless absolutely necessary:
   ```bash
   # ❌ Avoid this
   sudo gpkg install package

   # ✅ Prefer this
   gpkg install package --prefix=$HOME/.local
   ```

**Prevention:**
- Always use user-owned directories (default `~/.gpkg`)
- Configure `install_prefix` in config to a writable location
- Avoid system directories unless required

---

### Problem: Database Permission Error

**Symptoms:**
```
Error: Package database error
exit code: 8
attempt to write a readonly database
```

**Solution:**
1. **Check database permissions**:
   ```bash
   ls -lh ~/.gpkg/pkgdb.sqlite
   ```

2. **Fix ownership**:
   ```bash
   sudo chown $USER:$USER ~/.gpkg/pkgdb.sqlite
   chmod 644 ~/.gpkg/pkgdb.sqlite
   ```

3. **Fix directory permissions**:
   ```bash
   sudo chown -R $USER:$USER ~/.gpkg
   chmod 755 ~/.gpkg
   ```

4. **If database is locked by another process**:
   ```bash
   # Find process using database
   lsof ~/.gpkg/pkgdb.sqlite

   # Wait for other gpkg processes to finish, or kill if stuck
   pkill gpkg
   ```

**Prevention:**
- Avoid running gpkg with `sudo`
- Don't manually edit database files
- Ensure only one gpkg instance runs at a time

---

## Database Issues

### Problem: Database Corruption

**Symptoms:**
```
Error: Package database error
database disk image is malformed
```

**Solution:**

**Option 1: Repair database**
```bash
# Backup corrupted database
cp ~/.gpkg/pkgdb.sqlite ~/.gpkg/pkgdb.sqlite.corrupted

# Try to repair with SQLite
sqlite3 ~/.gpkg/pkgdb.sqlite ".recover" | sqlite3 ~/.gpkg/pkgdb.sqlite.recovered

# Replace with recovered version
mv ~/.gpkg/pkgdb.sqlite.recovered ~/.gpkg/pkgdb.sqlite
```

**Option 2: Rebuild from scratch**
```bash
# Backup current database
mv ~/.gpkg/pkgdb.sqlite ~/.gpkg/pkgdb.sqlite.backup

# Reinstall packages (database will be recreated)
gpkg list --installed > installed_packages.txt
cat installed_packages.txt | xargs -n1 gpkg install
```

**Option 3: Nuclear option**
```bash
# Completely reset gpkg (loses all package tracking)
rm -rf ~/.gpkg
mkdir -p ~/.gpkg
cp examples/config.example.toml ~/.gpkg/config.toml

# Reinstall packages from scratch
```

**Prevention:**
- Keep regular backups: `cp ~/.gpkg/pkgdb.sqlite ~/.gpkg/pkgdb.sqlite.backup`
- Avoid killing gpkg processes abruptly
- Don't run multiple gpkg instances simultaneously
- Ensure filesystem is not full (`df -h`)

---

### Problem: Database Lock Timeout

**Symptoms:**
```
Error: Package database error
database is locked
```

**Solution:**
1. **Wait for current operation to complete** - Another gpkg process may be running

2. **Check for running processes**:
   ```bash
   ps aux | grep gpkg
   ```

3. **Kill hung processes** (if stuck):
   ```bash
   pkill -9 gpkg
   ```

4. **Remove stale lock** (use with caution):
   ```bash
   # Check for SQLite lock files
   ls -la ~/.gpkg/*.sqlite*
   rm ~/.gpkg/pkgdb.sqlite-wal
   rm ~/.gpkg/pkgdb.sqlite-shm
   ```

5. **Increase busy timeout** in future (edit gpkg source or file issue)

**Prevention:**
- Run only one gpkg operation at a time
- Use `--dry-run` for testing without locking database
- Don't forcefully kill gpkg processes

---

## Common Error Messages

### Exit Code Reference

| Code | Meaning | Common Causes |
|------|---------|---------------|
| 0 | Success | Operation completed successfully |
| 1 | General failure | Unexpected error, check logs with `-v` |
| 2 | Usage error | Invalid command or arguments, check `--help` |
| 3 | Network error | Connection timeout, DNS failure, proxy issues |
| 4 | Checksum failed | Corrupted download, modified release asset |
| 5 | Install failed | Build error, missing dependencies, permission issues |
| 6 | Manifest validation | Invalid YAML, missing required fields |
| 7 | Package not found | Package doesn't exist, sources not updated |
| 8 | Database error | Database corruption, lock timeout, permission issues |

### Message: "No release assets found for platform"

**Solution:**
```bash
# Check available platforms
gpkg info owner/repo

# Try source build instead
gpkg install owner/repo --from-source

# Or specify architecture explicitly
gpkg install owner/repo --arch linux/arm64
```

---

### Message: "Manifest validation failed: missing required field 'name'"

**Solution:**
1. If installing from manifest file, fix the YAML:
   ```yaml
   name: my-package  # Required
   version: 1.0.0    # Required
   description: "My package"
   ```

2. Validate before installing:
   ```bash
   gpkg validate ./manifest.yaml
   ```

3. Use `--fix` to auto-correct simple issues:
   ```bash
   gpkg validate ./manifest.yaml --fix
   ```

---

### Message: "disk space quota exceeded"

**Solution:**
1. **Check disk space**:
   ```bash
   df -h ~/.gpkg
   ```

2. **Clean up old packages**:
   ```bash
   gpkg list --installed
   gpkg uninstall unused-package
   ```

3. **Change install prefix** to larger partition:
   ```bash
   gpkg config set install_prefix "/mnt/large-disk/gpkg"
   ```

4. **Clear download cache** (if implemented):
   ```bash
   rm -rf ~/.gpkg/cache/*
   ```

---

## Debug Mode

### Enabling Verbose Output

**Quick debugging:**
```bash
# Single -v for verbose
gpkg install package -v

# Double -vv for more verbose
gpkg install package -vv

# Triple -vvv for maximum verbosity
gpkg install package -vvv
```

**Set log level explicitly:**
```bash
gpkg install package --log-level debug
```

**Environment variable:**
```bash
export LOG_LEVEL=debug
gpkg install package
```

**In config file:**
```toml
log_level = "debug"  # error, warn, info, debug
```

### What to Look For in Debug Output

**Network issues:**
```
[DEBUG] Downloading https://github.com/owner/repo/releases/download/v1.0.0/asset.tar.gz
[DEBUG] HTTP Response: 404 Not Found
```

**Checksum problems:**
```
[DEBUG] Expected checksum: abc123...
[DEBUG] Calculated checksum: def456...
[ERROR] Checksum validation failed
```

**Build failures:**
```
[DEBUG] Running build step: make
[DEBUG] Exit code: 2
[DEBUG] stderr: gcc: command not found
```

**Permission issues:**
```
[DEBUG] Creating directory: /opt/gpkg/bin
[ERROR] mkdir: permission denied
```

### Capturing Debug Output

**Save to file:**
```bash
gpkg install package -vv 2>&1 | tee debug.log
```

**Only errors:**
```bash
gpkg install package 2> errors.log
```

**JSON output for parsing:**
```bash
gpkg install package --json > result.json 2> errors.json
```

---

## Reporting Issues

### Before Reporting

1. **Search existing issues**: [GitHub Issues](https://github.com/grave0x/gpkg/issues)
2. **Update to latest version**:
   ```bash
   gpkg --version
   # Check against latest release
   ```
3. **Try with verbose mode**:
   ```bash
   gpkg <command> -vv --log-level debug
   ```
4. **Test with `--dry-run`**:
   ```bash
   gpkg install package --dry-run
   ```

### Creating a Good Bug Report

**Required Information:**

1. **gpkg version:**
   ```bash
   gpkg --version
   ```

2. **Operating system:**
   ```bash
   uname -a
   # or
   cat /etc/os-release
   ```

3. **Go version** (if built from source):
   ```bash
   go version
   ```

4. **Command that failed:**
   ```bash
   gpkg install owner/repo --verbose
   ```

5. **Full error output** (use debug mode):
   ```bash
   gpkg install owner/repo -vv --log-level debug 2>&1 | tee error.log
   ```

6. **Configuration** (sanitize sensitive data):
   ```bash
   gpkg config show
   ```

7. **Steps to reproduce**

### Issue Template

```markdown
**Describe the bug**
A clear description of what went wrong.

**To Reproduce**
Steps to reproduce the behavior:
1. Run `gpkg ...`
2. See error

**Expected behavior**
What you expected to happen.

**Environment**
- OS: [e.g. Ubuntu 22.04]
- gpkg version: [e.g. 0.1.0]
- Go version (if built from source): [e.g. 1.21.0]

**Logs**
```
Paste debug output here (gpkg <command> -vv --log-level debug)
```

**Configuration**
```toml
Paste relevant config.toml sections (remove sensitive data)
```

**Additional context**
Any other information that might be helpful.
```

### Where to Get Help

- **GitHub Issues**: [https://github.com/grave0x/gpkg/issues](https://github.com/grave0x/gpkg/issues)
  - Bug reports
  - Feature requests
  - Installation problems

- **GitHub Discussions**: [https://github.com/grave0x/gpkg/discussions](https://github.com/grave0x/gpkg/discussions)
  - Questions
  - General help
  - Ideas and suggestions

- **Documentation**:
  - [README](../README.md) - Quick start and overview
  - [CLI Specification](../CLI+SPEC.md) - Command reference
  - [Development Guide](../DEVELOPMENT.md) - Contributing
  - [Wiki](https://github.com/grave0x/gpkg/wiki) - Detailed guides

- **Check existing resources**:
  ```bash
  # Built-in help
  gpkg --help
  gpkg install --help

  # Example configurations
  cat examples/config.example.toml
  cat examples/manifest.yaml
  ```

---

## Quick Reference

### Diagnostic Commands

```bash
# Check version
gpkg --version

# Show configuration
gpkg config show

# List sources
gpkg list-sources

# Test network connectivity
gpkg update --verbose

# Validate manifest
gpkg validate ./manifest.yaml

# Dry-run installation
gpkg install package --dry-run

# Debug installation
gpkg install package -vv --log-level debug
```

### Emergency Recovery

```bash
# Backup everything
cp -r ~/.gpkg ~/.gpkg.backup

# Reset configuration
mv ~/.gpkg/config.toml ~/.gpkg/config.toml.backup
cp examples/config.example.toml ~/.gpkg/config.toml

# Repair database
sqlite3 ~/.gpkg/pkgdb.sqlite "PRAGMA integrity_check;"

# Complete reset (last resort)
rm -rf ~/.gpkg
mkdir -p ~/.gpkg
```

### Useful Flags

- `--dry-run` - Preview without making changes
- `--verbose` or `-v` - Increase output detail
- `--log-level debug` - Maximum debugging output
- `--json` - Machine-readable output
- `--yes` or `-y` - Skip confirmations (for scripts)
- `--offline` - Work without network access
- `--help` - Show command help

---

## Still Having Issues?

If you've tried the solutions in this guide and are still experiencing problems:

1. **Enable maximum debugging**:
   ```bash
   gpkg <command> -vvv --log-level debug 2>&1 | tee full-debug.log
   ```

2. **Gather system information**:
   ```bash
   gpkg --version
   uname -a
   go version
   sqlite3 --version
   ```

3. **File a detailed issue**: [New Issue](https://github.com/grave0x/gpkg/issues/new)

4. **Include**:
   - Complete command that failed
   - Full debug output
   - System information
   - Config file (sanitized)
   - Steps to reproduce

We're here to help! The more information you provide, the faster we can resolve your issue.
