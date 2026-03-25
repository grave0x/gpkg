# Manifest Format Reference

## Overview

A **manifest** is a YAML file that describes how to install and manage a package in gpkg. Manifests define package metadata, installation methods (prebuilt binaries or build from source), dependencies, and platform-specific configurations.

Manifests enable:
- **Release-based installation**: Download prebuilt binaries from GitHub releases or other URLs
- **Source builds**: Build packages from source with custom build commands
- **Dependency management**: Declare package dependencies for automatic installation
- **Platform targeting**: Specify platform-specific assets and configurations
- **Checksum verification**: Ensure integrity of downloaded assets

## Format

Manifests are written in **YAML** format and follow a structured schema with three main sections:

```yaml
package:
  # Package metadata

install:
  # Release/binary installation specification

build_source:
  # Optional: build from source specification
```

## Schema Structure

### Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `package` | object | ✓ | Package metadata and identification |
| `install` | object | * | Release/binary installation specification |
| `build_source` | object | * | Build from source specification |
| `dependencies` | array | | List of package dependencies |

\* At least one of `install` or `build_source` must be specified.

## Required Fields

### Package Section

The `package` section defines core package metadata:

```yaml
package:
  name: package-name      # Required: Package identifier
  version: 1.2.0          # Required: Semantic version
  author: "Name <email>"  # Optional: Package author
  url: "https://..."      # Optional: Homepage URL
  license: "MIT"          # Optional: License identifier
```

#### Package Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✓ | Unique package name (lowercase, alphanumeric, hyphens) |
| `version` | string | ✓ | Semantic version (e.g., `1.2.0`, `0.1.0-beta`) |
| `author` | string | | Author name and email |
| `url` | string | | Package homepage or repository URL |
| `license` | string | | License identifier (e.g., `MIT`, `Apache-2.0`, `GPL-3.0`) |

**Example:**
```yaml
package:
  name: cool-tool
  version: 1.2.0
  author: "Owner <owner@example.com>"
  url: "https://github.com/owner/cool-tool"
  license: "MIT"
```

## Optional Fields

### Install Section

The `install` section specifies how to install prebuilt binaries or release assets:

```yaml
install:
  type: release           # Required: Installation type
  source: "https://..."   # Required: Download URL or pattern
  pattern: "..."          # Optional: URL pattern with variables
  checksum:               # Optional: Checksums for verification
    sha256: "abc123..."
  extract_path: "..."     # Optional: Path within archive
  executable: "..."       # Optional: Executable name/path
  post_install: "..."     # Optional: Post-install script
```

#### Install Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✓ | Installation type: `release`, `binary`, or `archive` |
| `source` | string | ✓ | Download URL or source location |
| `pattern` | string | | URL pattern with `{version}`, `{platform}`, `{arch}` variables |
| `checksum` | map | | Checksum values (e.g., `sha256`, `sha512`) |
| `extract_path` | string | | Path to extract from archive (for archives) |
| `executable` | string | | Name of the executable file |
| `post_install` | string | | Command to run after installation |

**Installation Types:**
- `release`: Download from GitHub releases or similar
- `binary`: Direct binary download
- `archive`: Download and extract from archive (`.tar.gz`, `.zip`)

**Example:**
```yaml
install:
  type: archive
  source: "https://github.com/owner/tool/releases/download/v{version}/tool-{platform}-{arch}.tar.gz"
  checksum:
    sha256: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
  extract_path: "bin/tool"
  executable: "tool"
```

### Build Source Section

The `build_source` section specifies how to build the package from source:

```yaml
build_source:
  type: git               # Required: Source type
  source: "https://..."   # Required: Repository or source URL
  tag: "v1.2.0"          # Optional: Git tag to checkout
  branch: "main"         # Optional: Git branch to checkout
  commands:              # Required: Build commands
    - "make"
    - "make install"
  env:                   # Optional: Build environment variables
    CGO_ENABLED: "0"
```

#### Build Source Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ✓ | Source type: `git` or `tarball` |
| `source` | string | ✓ | Repository URL or tarball location |
| `tag` | string | | Git tag to checkout (for git sources) |
| `branch` | string | | Git branch to checkout (for git sources) |
| `commands` | array | ✓ | Build commands to execute in order |
| `env` | map | | Environment variables for build process |

**Build Commands:**
- Commands are executed sequentially
- Use `{install_prefix}` variable for installation directory
- Exit code 0 indicates success; non-zero indicates failure

**Example:**
```yaml
build_source:
  type: git
  source: "https://github.com/owner/cool-tool.git"
  tag: "v1.2.0"
  commands:
    - "make build"
    - "make install PREFIX={install_prefix}"
  env:
    CGO_ENABLED: "0"
    GOFLAGS: "-buildmode=pie"
```

### Dependencies

Declare packages that must be installed before this package:

```yaml
dependencies:
  - package-one
  - package-two
  - package-three
```

#### Dependency Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `dependencies` | array | | List of package names (strings) |

Dependencies are installed in order before the current package. gpkg automatically resolves and installs dependencies recursively.

**Example:**
```yaml
dependencies:
  - libssl
  - zlib
  - curl
```

## Platform Support

### Platform Variables

Use platform variables in URL patterns and commands for cross-platform support:

| Variable | Description | Example Values |
|----------|-------------|----------------|
| `{platform}` | Operating system | `linux`, `darwin`, `windows` |
| `{arch}` | CPU architecture | `amd64`, `arm64`, `386` |
| `{version}` | Package version | `1.2.0` |
| `{install_prefix}` | Installation directory | `/home/user/.local/gpkg/packages/tool/1.2.0` |

**Example with platform variables:**
```yaml
install:
  type: archive
  source: "https://releases.example.com/{version}/tool-{platform}-{arch}.tar.gz"
```

This expands to different URLs based on the platform:
- Linux AMD64: `https://releases.example.com/1.2.0/tool-linux-amd64.tar.gz`
- macOS ARM64: `https://releases.example.com/1.2.0/tool-darwin-arm64.tar.gz`
- Windows AMD64: `https://releases.example.com/1.2.0/tool-windows-amd64.tar.gz`

### Platform-Specific Assets

For multiple platform releases, the manifest can specify different checksums per platform or use conditional logic in scripts.

**Example:**
```yaml
install:
  type: release
  pattern: "https://github.com/owner/tool/releases/download/v{version}/tool-{platform}-{arch}.tar.gz"
  # Checksum verification happens automatically per platform
```

## Checksums

### Supported Hash Algorithms

| Algorithm | Field Name | Example |
|-----------|------------|---------|
| SHA-256 | `sha256` | `0123456789abcdef...` (64 hex chars) |
| SHA-512 | `sha512` | `0123456789abcdef...` (128 hex chars) |

### Checksum Verification

Checksums ensure downloaded files haven't been tampered with:

```yaml
install:
  source: "https://example.com/tool.tar.gz"
  checksum:
    sha256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
```

**Best Practices:**
- Always include checksums for release assets
- Use SHA-256 or SHA-512 for security
- Generate checksums with: `sha256sum file.tar.gz` or `shasum -a 256 file.tar.gz`

## Complete Examples

### Example 1: Simple Binary Release

```yaml
package:
  name: simple-tool
  version: 2.1.0
  author: "Developer <dev@example.com>"
  url: "https://github.com/dev/simple-tool"
  license: "MIT"

install:
  type: binary
  source: "https://github.com/dev/simple-tool/releases/download/v2.1.0/simple-tool"
  checksum:
    sha256: "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2"
  executable: "simple-tool"
```

### Example 2: Archive with Platform Support

```yaml
package:
  name: cross-platform-app
  version: 3.0.1
  author: "Team <team@example.com>"
  url: "https://cross-platform-app.dev"
  license: "Apache-2.0"

install:
  type: archive
  source: "https://releases.example.com/v{version}/app-{platform}-{arch}.tar.gz"
  checksum:
    sha256: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
  extract_path: "bin/app"
  executable: "app"
```

### Example 3: Build from Source

```yaml
package:
  name: build-tool
  version: 1.5.0
  author: "Builder <builder@example.com>"
  url: "https://github.com/builder/build-tool"
  license: "GPL-3.0"

build_source:
  type: git
  source: "https://github.com/builder/build-tool.git"
  tag: "v1.5.0"
  commands:
    - "./configure --prefix={install_prefix}"
    - "make -j4"
    - "make install"
  env:
    CC: "gcc"
    CFLAGS: "-O2 -march=native"
```

### Example 4: Package with Dependencies

```yaml
package:
  name: dependent-app
  version: 2.0.0
  author: "Developer <dev@example.com>"
  url: "https://github.com/dev/dependent-app"
  license: "MIT"

install:
  type: archive
  source: "https://github.com/dev/dependent-app/releases/download/v{version}/app-{platform}-{arch}.tar.gz"
  checksum:
    sha256: "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
  extract_path: "app"
  executable: "app"

dependencies:
  - libcurl
  - openssl
  - zlib
```

### Example 5: Complex Multi-Step Build

```yaml
package:
  name: complex-project
  version: 4.2.1
  author: "Open Source Team <oss@example.com>"
  url: "https://github.com/oss/complex-project"
  license: "BSD-3-Clause"

build_source:
  type: git
  source: "https://github.com/oss/complex-project.git"
  branch: "release-4.2"
  tag: "v4.2.1"
  commands:
    - "git submodule update --init --recursive"
    - "mkdir -p build && cd build"
    - "cmake -DCMAKE_INSTALL_PREFIX={install_prefix} .."
    - "cmake --build . --parallel"
    - "cmake --install ."
  env:
    CMAKE_BUILD_TYPE: "Release"
    BUILD_SHARED_LIBS: "ON"

dependencies:
  - cmake
  - gcc
  - pkg-config
```

### Example 6: Post-Install Hook

```yaml
package:
  name: config-tool
  version: 1.0.0
  author: "Admin <admin@example.com>"
  url: "https://config-tool.example.com"
  license: "MIT"

install:
  type: archive
  source: "https://releases.example.com/config-tool-{version}.tar.gz"
  checksum:
    sha256: "fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"
  extract_path: "config-tool"
  executable: "config-tool"
  post_install: "config-tool --init-config"
```

## Validation Rules

The manifest parser validates the following rules:

### Required Validations

1. **Package name**: Must be present and non-empty
2. **Package version**: Must be present and non-empty
3. **Installation method**: At least one of `install` or `build_source` must be specified
4. **Install source**: If `install` is specified, `source` must be present
5. **Build source**: If `build_source` is specified:
   - `source` must be present
   - At least one `commands` entry must be present

### Format Validations

- YAML syntax must be valid
- Field types must match schema (string, array, map)
- Version should follow semantic versioning (recommended)
- URLs should be valid HTTP/HTTPS URLs (recommended)

### Best Practices

- Use semantic versioning (e.g., `1.2.3`, `2.0.0-beta.1`)
- Include checksums for all release assets
- Specify license for legal clarity
- Document dependencies explicitly
- Test build commands before publishing
- Use platform variables for cross-platform support

## Schema Reference

### Complete YAML Schema

```yaml
# Required package metadata
package:
  name: string              # Required: package identifier
  version: string           # Required: semantic version
  author: string            # Optional: author information
  url: string               # Optional: homepage URL
  license: string           # Optional: license identifier

# Optional: install from release/binary
install:
  type: string              # Required: "release" | "binary" | "archive"
  source: string            # Required: download URL
  pattern: string           # Optional: URL pattern with variables
  checksum:                 # Optional: checksum verification
    sha256: string          # SHA-256 hash
    sha512: string          # SHA-512 hash
  extract_path: string      # Optional: archive extraction path
  executable: string        # Optional: executable name
  post_install: string      # Optional: post-install command

# Optional: build from source
build_source:
  type: string              # Required: "git" | "tarball"
  source: string            # Required: repository URL or tarball
  tag: string               # Optional: git tag
  branch: string            # Optional: git branch
  commands:                 # Required: array of build commands
    - string
  env:                      # Optional: environment variables
    KEY: value

# Optional: package dependencies
dependencies:               # Optional: array of package names
  - string
```

## See Also

- [Creating Packages](../User-Guide/Creating-Packages.md) - Guide to creating your own package manifests
- [Package Registry](Package-Registry.md) - Publishing packages to registries
- [Getting Started](../Getting-Started.md) - Installation and basic usage
- [Command Reference](Command-Reference.md) - gpkg command documentation
