# Basic Commands

This guide provides comprehensive documentation for all `gpkg` commands, including practical examples and best practices.

## Table of Contents

- [Package Installation](#package-installation)
  - [install](#install)
  - [uninstall](#uninstall)
  - [upgrade](#upgrade)
  - [rollback](#rollback)
- [Source Management](#source-management)
  - [add-source](#add-source)
  - [remove-source](#remove-source)
  - [list-sources](#list-sources)
  - [update](#update)
- [Package Discovery](#package-discovery)
  - [search](#search)
  - [info](#info)
  - [list](#list)
- [Configuration & Utilities](#configuration--utilities)
  - [config](#config)
  - [validate](#validate)
  - [completion](#completion)
- [Global Flags](#global-flags)

---

## Package Installation

### install

**Purpose:** Install a package from a release binary or build from source.

**Description:** The `install` command supports multiple installation methods: from package names (requires configured sources), local manifest files, or GitHub repositories. You can choose between installing pre-built release binaries or building from source.

**Syntax:**
```bash
gpkg install <pkg|manifest> [--from-release|--from-source] [flags]
```

**Common Flags:**
- `--from-release` - Install from a pre-built release binary (faster)
- `--from-source` - Build and install from source code
- `--prefix <path>` - Override the default installation directory
- `--dry-run` - Preview the installation without making changes
- `-y, --yes` - Skip confirmation prompts (non-interactive mode)

**Examples:**

```bash
# Install a package from release binary
gpkg install owner/repo --from-release

# Install by building from source
gpkg install owner/repo --from-source

# Install from a local manifest file
gpkg install ./examples/manifest.yaml

# Install to a custom directory
gpkg install my-package --prefix=/opt/mytools

# Preview installation without executing
gpkg install my-package --dry-run

# Install with verbose output
gpkg install my-package -vv

# Install non-interactively (assume yes to all prompts)
gpkg install my-package -y
```

**Tips & Best Practices:**
- Use `--from-release` for faster installations when available
- Use `--dry-run` first to preview what will be installed and where
- Keep the default prefix (`~/.gpkg`) for user-local installations to avoid needing sudo
- Check package info with `gpkg info <package>` before installing
- Use `--from-source` when you need custom build configurations or the latest development version

---

### uninstall

**Purpose:** Remove an installed package from the system.

**Description:** The `uninstall` command cleanly removes a package and all its tracked files from the installation prefix. Only files within the gpkg-managed prefix are removed for safety.

**Syntax:**
```bash
gpkg uninstall <pkg> [flags]
```

**Common Flags:**
- `--prefix <path>` - Specify installation prefix if non-default was used
- `--dry-run` - Preview what will be removed without executing
- `-y, --yes` - Skip confirmation prompts

**Examples:**

```bash
# Uninstall a package
gpkg uninstall my-package

# Uninstall from custom prefix
gpkg uninstall my-package --prefix=/opt/mytools

# Preview what will be removed
gpkg uninstall my-package --dry-run

# Uninstall without confirmation
gpkg uninstall my-package -y

# Uninstall with verbose output to see all removed files
gpkg uninstall my-package -v
```

**Tips & Best Practices:**
- Use `--dry-run` to verify which files will be removed before actual uninstallation
- Uninstalling removes only files tracked in the package database (gpkg-managed files)
- If you installed with a custom `--prefix`, you must specify the same prefix when uninstalling
- Check installed packages with `gpkg list --installed` before uninstalling

---

### upgrade

**Purpose:** Upgrade installed packages to their latest available versions.

**Description:** The `upgrade` command updates one or more packages to the newest versions available from configured sources. If no package names are specified, all installed packages are upgraded.

**Syntax:**
```bash
gpkg upgrade [pkg...] [flags]
```

**Common Flags:**
- `--dry-run` - Preview which packages will be upgraded
- `-y, --yes` - Skip confirmation prompts
- `--json` - Output upgrade information in JSON format

**Examples:**

```bash
# Upgrade all installed packages
gpkg upgrade

# Upgrade a specific package
gpkg upgrade my-package

# Upgrade multiple specific packages
gpkg upgrade pkg1 pkg2 pkg3

# Preview available upgrades without installing
gpkg upgrade --dry-run

# Upgrade all packages non-interactively
gpkg upgrade -y

# Upgrade with detailed output
gpkg upgrade -vv

# Get upgrade information in JSON format
gpkg upgrade --json
```

**Tips & Best Practices:**
- Run `gpkg update` before upgrading to ensure you have the latest package indices
- Use `--dry-run` to see which packages have updates available before upgrading
- Consider upgrading packages individually for critical systems to isolate potential issues
- Use `gpkg list --installed` to check current package versions
- Keep old versions in the database so you can rollback if needed

---

### rollback

**Purpose:** Rollback an installed package to a previous version.

**Description:** The `rollback` command reverts a package to a previously installed version. This requires version history to be available in the package database.

**Syntax:**
```bash
gpkg rollback <pkg> [flags]
```

**Common Flags:**
- `--to-version <version>` - **Required.** The target version to rollback to
- `--prefix <path>` - Specify installation prefix if non-default
- `--dry-run` - Preview the rollback operation
- `-y, --yes` - Skip confirmation prompts

**Examples:**

```bash
# Rollback a package to version 1.0.0
gpkg rollback my-tool --to-version 1.0.0

# Rollback with custom prefix
gpkg rollback my-tool --to-version 1.0.0 --prefix /opt/gpkg

# Preview rollback operation
gpkg rollback my-tool --to-version 1.0.0 --dry-run

# Rollback non-interactively
gpkg rollback my-tool --to-version 1.0.0 -y
```

**Tips & Best Practices:**
- Use `gpkg info <package>` to check available version history
- Rollback is useful after problematic upgrades
- Always specify the exact version string (e.g., `1.0.0`, not `1.0` or `v1.0.0`)
- Test the rollback with `--dry-run` first
- Consider documenting why you rolled back for future reference

---

## Source Management

### add-source

**Purpose:** Add a new package source to the configuration.

**Description:** Package sources are repositories or indices that contain package manifests. Adding sources allows `gpkg` to discover and install packages from those locations.

**Syntax:**
```bash
gpkg add-source <uri> [flags]
```

**Common Flags:**
- `--dry-run` - Preview adding the source without modifying configuration
- `-y, --yes` - Skip confirmation prompts

**Examples:**

```bash
# Add a package source
gpkg add-source https://packages.example.com/index.json

# Add a source with verbose output
gpkg add-source https://repo.example.org/packages.json -v

# Preview adding a source
gpkg add-source https://packages.example.com/index.json --dry-run

# Add source non-interactively
gpkg add-source https://packages.example.com/index.json -y
```

**Tips & Best Practices:**
- Verify the source URL is trustworthy before adding
- Use HTTPS URLs for security
- Run `gpkg update` after adding a new source to fetch its package index
- Use `gpkg list-sources` to view all configured sources
- Sources can be removed later with `gpkg remove-source` if needed

---

### remove-source

**Purpose:** Remove a package source from the configuration.

**Description:** Removes a previously added package source by its ID or URI. Packages from removed sources can still be used if already installed, but new packages from that source won't be discoverable.

**Syntax:**
```bash
gpkg remove-source <id|uri> [flags]
```

**Common Flags:**
- `--dry-run` - Preview removal without executing
- `-y, --yes` - Skip confirmation prompts

**Examples:**

```bash
# Remove a source by ID
gpkg remove-source source-1

# Remove a source by URI
gpkg remove-source https://packages.example.com/index.json

# Preview source removal
gpkg remove-source source-1 --dry-run

# Remove source without confirmation
gpkg remove-source source-1 -y
```

**Tips & Best Practices:**
- Use `gpkg list-sources` to find the source ID or URI
- Removing a source doesn't uninstall packages already installed from it
- Removed sources can be re-added later if needed
- Consider keeping sources you might need in the future

---

### list-sources

**Purpose:** Display all registered package sources.

**Description:** Lists all configured package sources with their IDs, URIs, and status information.

**Syntax:**
```bash
gpkg list-sources [flags]
```

**Common Flags:**
- `--json` - Output in JSON format for scripting

**Examples:**

```bash
# List all sources
gpkg list-sources

# List sources in JSON format
gpkg list-sources --json

# List sources with minimal output
gpkg list-sources --quiet
```

**Tips & Best Practices:**
- Use this command to verify sources after adding or removing them
- JSON output is useful for scripting and automation
- Check source status before running `gpkg update`
- Keep track of which sources provide which packages

---

### update

**Purpose:** Refresh package source metadata from all configured sources.

**Description:** The `update` command fetches the latest package lists and version information from all enabled sources. This doesn't install or upgrade any packages, only refreshes the available package metadata.

**Syntax:**
```bash
gpkg update [flags]
```

**Common Flags:**
- `--dry-run` - Preview what would be updated
- `--json` - Output update results in JSON format
- `--offline` - Run in offline mode (skip network operations)

**Examples:**

```bash
# Update all package sources
gpkg update

# Update with verbose output
gpkg update -v

# Preview update operation
gpkg update --dry-run

# Update with JSON output
gpkg update --json
```

**Tips & Best Practices:**
- Run `gpkg update` regularly to keep package indices fresh (similar to `apt update`)
- Always run `update` before `upgrade` to ensure you get the latest versions
- Run `update` after adding new sources
- Use `-v` to see which sources are being updated
- Update is safe to run frequently; it only fetches metadata

---

## Package Discovery

### search

**Purpose:** Search across configured package sources for packages matching a search term.

**Description:** The `search` command queries all configured sources (or a specific source) for packages matching your search term. It searches package names, descriptions, and other metadata.

**Syntax:**
```bash
gpkg search <term> [flags]
```

**Common Flags:**
- `--source <name>` - Search only in a specific source
- `--json` - Output results in JSON format

**Examples:**

```bash
# Search for packages containing "golang"
gpkg search golang

# Search in a specific source
gpkg search --source example golang

# Search with JSON output for scripting
gpkg search --json curl

# Search with verbose output
gpkg search -v tool

# Search quietly (minimal output)
gpkg search --quiet editor
```

**Tips & Best Practices:**
- Run `gpkg update` before searching to ensure you're searching current indices
- Use specific search terms for better results
- Combine with `gpkg info` to get detailed information about found packages
- Use `--source` to narrow searches when you know which repository contains the package
- JSON output is useful for integrating with other tools

---

### info

**Purpose:** Display detailed information about a package.

**Description:** The `info` command shows comprehensive information about a package from a repository, installed packages, or a manifest file. This includes version, description, dependencies, installation files, and more.

**Syntax:**
```bash
gpkg info <repo|pkg|manifest> [flags]
```

**Common Flags:**
- `--deps-tree` - Display dependency tree
- `--raw` - Show raw manifest YAML
- `--parsed` - Show parsed manifest structure
- `--json` - Output in JSON format

**Examples:**

```bash
# Show info for a GitHub repository
gpkg info owner/repo

# Show info for an installed package
gpkg info my-package

# Show info from a local manifest
gpkg info ./manifest.yaml

# Show dependency tree
gpkg info my-package --deps-tree

# Show raw manifest
gpkg info owner/repo --raw

# Show parsed manifest structure
gpkg info ./manifest.yaml --parsed

# Get package info in JSON format
gpkg info my-package --json
```

**Tips & Best Practices:**
- Use `info` before installing to review what will be installed
- Check `--deps-tree` to understand package dependencies
- Use `--raw` when you need to see the exact manifest definition
- Combine with `--json` for parsing package metadata in scripts
- Review checksums and URLs before installation for security

---

### list

**Purpose:** List installed or available packages.

**Description:** The `list` command displays packages either installed on your system or available from configured sources. Results can be filtered and sorted.

**Syntax:**
```bash
gpkg list [--installed|--available] [flags]
```

**Common Flags:**
- `--installed` - Show installed packages (default)
- `--available` - Show available packages from sources
- `--filter <pattern>` - Filter results by name pattern
- `--sort <field>` - Sort by name, version, or date (default: name)
- `--json` - Output in JSON format

**Examples:**

```bash
# List installed packages (default)
gpkg list
gpkg list --installed

# List available packages from all sources
gpkg list --available

# Filter installed packages by name pattern
gpkg list --filter "tool"

# Sort packages by version
gpkg list --sort version

# Sort packages by installation date
gpkg list --sort date

# List packages in JSON format
gpkg list --json

# List available packages matching a pattern
gpkg list --available --filter "go-"

# List with verbose output
gpkg list -v
```

**Tips & Best Practices:**
- Use `--installed` to check what's currently on your system
- Use `--available` to discover packages before installing
- Combine `--filter` and `--sort` for better organization
- JSON output is great for scripting and automation
- Run `gpkg update` before using `--available` for current results

---

## Configuration & Utilities

### config

**Purpose:** Manage gpkg configuration settings.

**Description:** The `config` command allows you to view and modify gpkg configuration. Configuration can be stored in `~/.gpkg/config.toml` or `/etc/gpkg/config.toml`.

**Syntax:**
```bash
gpkg config [command]
```

**Subcommands:**
- `get <key>` - Get a configuration value
- `set <key> <value>` - Set a configuration value
- `show` - Display merged configuration

**Common Flags:**
- `--json` - Output in JSON format (for `show` and `get`)
- `-c, --config <path>` - Use a specific config file

**Examples:**

```bash
# Show all configuration
gpkg config show

# Show configuration in JSON format
gpkg config show --json

# Get a specific configuration value
gpkg config get prefix
gpkg config get log_level

# Set a configuration value
gpkg config set log_level debug
gpkg config set parallel_downloads 8
gpkg config set require_checksums true

# Use a custom config file
gpkg --config /path/to/config.toml config show
```

**Configuration Keys:**
- `install_prefix` - Default installation directory (default: `~/.gpkg`)
- `pkgdb_path` - Package database location (default: `~/.gpkg/pkgdb.sqlite`)
- `sources_dir` - Directory for source configurations (default: `~/.gpkg/sources.d`)
- `require_checksums` - Require checksum verification (default: `true`)
- `parallel_downloads` - Number of parallel downloads (default: `4`)
- `log_level` - Logging verbosity (error, warn, info, debug)

**Tips & Best Practices:**
- Use `config show` to see your current configuration and defaults
- Keep `require_checksums` enabled for security
- Adjust `parallel_downloads` based on your network capacity
- Store user-specific config in `~/.gpkg/config.toml`
- Use environment variables (e.g., `GPKG_INSTALL_PREFIX`) for temporary overrides
- Back up your config file before making significant changes

---

### validate

**Purpose:** Validate a manifest file against the gpkg schema.

**Description:** The `validate` command checks manifest files for syntax errors, required fields, proper structure, and valid values. It can also attempt to fix trivial issues.

**Syntax:**
```bash
gpkg validate <manifest-path> [flags]
```

**Common Flags:**
- `--fix` - Attempt to automatically fix trivial issues
- `--json` - Output validation results in JSON format

**Validation Checks:**
- Required fields (name, version)
- Valid install or build_source specification
- Proper checksum format (SHA256/SHA512)
- Valid build commands and steps
- Correct YAML syntax
- Platform specifications

**Examples:**

```bash
# Validate a manifest file
gpkg validate ./my-tool.yaml

# Validate and attempt to fix issues
gpkg validate ./manifest.yaml --fix

# Validate with JSON output
gpkg validate ./manifest.yaml --json

# Validate with verbose output
gpkg validate ./manifest.yaml -v
```

**Tips & Best Practices:**
- Always validate manifests before using them for installation
- Use `--fix` to automatically correct simple formatting issues
- Check validation output carefully before using `--fix`
- Validate manifests in CI/CD pipelines to catch errors early
- Use JSON output for integration with other validation tools
- Keep manifests in version control and validate on each change

---

### completion

**Purpose:** Generate shell completion scripts for bash, zsh, or fish.

**Description:** The `completion` command generates shell completion scripts that enable tab-completion for gpkg commands, flags, and arguments in your shell.

**Syntax:**
```bash
gpkg completion [bash|zsh|fish] [flags]
```

**Supported Shells:**
- `bash` - Bash shell completion
- `zsh` - Zsh shell completion
- `fish` - Fish shell completion

**Examples:**

```bash
# Generate bash completion
gpkg completion bash

# Install bash completion (system-wide)
gpkg completion bash | sudo tee /usr/share/bash-completion/completions/gpkg

# Install bash completion (user-only)
gpkg completion bash > ~/.local/share/bash-completion/completions/gpkg

# Generate and install zsh completion
gpkg completion zsh | sudo tee /usr/share/zsh/site-functions/_gpkg

# Install zsh completion (user-only, oh-my-zsh)
gpkg completion zsh > ~/.oh-my-zsh/completions/_gpkg

# Generate and install fish completion
gpkg completion fish | sudo tee /usr/share/fish/vendor_completions.d/gpkg.fish

# Install fish completion (user-only)
gpkg completion fish > ~/.config/fish/completions/gpkg.fish
```

**Tips & Best Practices:**
- Install completions after installing gpkg for a better user experience
- For user-only installation, ensure the completion directory exists first
- Restart your shell or source the completion file after installation
- Update completions after upgrading gpkg if commands change
- Use system-wide installation on multi-user systems
- Check shell-specific documentation for custom completion directories

**Shell-Specific Activation:**

**Bash:**
```bash
# For system-wide installation, completions load automatically
# For user-only, add to ~/.bashrc:
source ~/.local/share/bash-completion/completions/gpkg
```

**Zsh:**
```bash
# Ensure completion system is initialized in ~/.zshrc:
autoload -Uz compinit
compinit
```

**Fish:**
```bash
# Fish auto-loads completions from ~/.config/fish/completions/
# No additional configuration needed
```

---

## Global Flags

These flags are available for all `gpkg` commands:

### Configuration & Output

- `-c, --config <path>` - Specify a custom config file (overrides default locations)
- `--json` - Output in machine-readable JSON format (supported by most commands)
- `--no-color` - Disable colored output (useful for logs and scripts)

### Verbosity & Logging

- `-v, --verbose` - Increase output verbosity (stackable: `-v`, `-vv`, `-vvv`)
- `--quiet` - Minimal output (conflicts with `--verbose`)
- `--log-level <level>` - Set logging level: `error`, `warn`, `info`, `debug` (default: `info`)

### Execution Control

- `--dry-run` - Preview actions without modifying filesystem or database
- `-y, --yes` - Non-interactive mode; assume "yes" for all prompts
- `--offline` - Disable network operations (use cached data only)

### Other

- `-h, --help` - Display help information for any command
- `--version` - Display gpkg version

**Examples:**

```bash
# Dry-run installation to preview
gpkg install my-package --dry-run

# Install with JSON output for scripting
gpkg install my-package --json

# Very verbose output for debugging
gpkg install my-package -vvv

# Quiet installation (minimal output)
gpkg install my-package --quiet

# Non-interactive upgrade (CI/CD)
gpkg upgrade -y --json

# Use custom config file
gpkg --config /etc/my-gpkg.toml list

# Offline mode (use cache only)
gpkg search tool --offline

# Set debug logging level
gpkg install my-package --log-level debug

# Combine multiple flags
gpkg install my-package --dry-run --json -v
```

**Tips for Global Flags:**
- Use `--dry-run` extensively when testing or learning gpkg
- Combine `--json` with `jq` for powerful scripting workflows
- Use `-y` in scripts and CI/CD pipelines to avoid hanging on prompts
- `--offline` is useful when working with limited or no network connectivity
- Stack `-v` flags for increased verbosity: `-v` (verbose), `-vv` (more verbose), `-vvv` (very verbose)
- `--quiet` and `--verbose` are mutually exclusive
- Global flags must be placed appropriately (before or after the command, depending on the flag)

---

## Quick Reference

### Most Common Workflows

**Initial Setup:**
```bash
gpkg add-source https://packages.example.com/index.json
gpkg update
```

**Find and Install a Package:**
```bash
gpkg search my-tool
gpkg info owner/repo
gpkg install owner/repo
```

**Keep Packages Updated:**
```bash
gpkg update
gpkg upgrade --dry-run
gpkg upgrade
```

**Manage Installed Packages:**
```bash
gpkg list --installed
gpkg info my-package
gpkg uninstall my-package
```

**Troubleshooting:**
```bash
gpkg validate ./manifest.yaml
gpkg info my-package --deps-tree
gpkg install my-package -vv --dry-run
gpkg rollback my-package --to-version 1.0.0
```

---

## Environment Variables

gpkg respects the following environment variables:

- `GPKG_INSTALL_PREFIX` - Override default installation directory
- `GPKG_PKGDB_PATH` - Override package database location
- `GPKG_CONFIG` - Override config file path
- `LOG_LEVEL` - Set logging level (error, warn, info, debug)

**Example:**
```bash
# Install to a custom location via environment variable
GPKG_INSTALL_PREFIX=/opt/mytools gpkg install my-package

# Use debug logging
LOG_LEVEL=debug gpkg install my-package

# Use a custom config file
GPKG_CONFIG=/etc/my-gpkg.toml gpkg list
```

---

## Additional Resources

- **Main README**: [README.md](../../README.md)
- **API Documentation**: [API.md](../../API.md)
- **CLI Specification**: [CLI+SPEC.md](../../CLI+SPEC.md)
- **Contributing Guide**: [CONTRIBUTING.md](../../CONTRIBUTING.md)
- **Example Manifests**: [examples/](../../examples/)

For more help with a specific command, use:
```bash
gpkg <command> --help
```
