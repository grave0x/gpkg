# Project Scope: gpkg

## 1. Project Overview
**gpkg** is a lightweight, user-focused package manager designed to install software either by downloading release binaries or building from source based on a declarative manifest. It isolates packages in a user-defined directory and manages them via a CLI interface similar to `pacman` or `apt`.

## 2. Objectives
- Provide a unified interface for installing binary releases and source builds.
- Ensure package isolation via custom installation prefixes.
- Manage dependencies and versioning through a simple manifest system.
- Offer a straightforward "update and upgrade" workflow for custom tools and GitHub repositories.

## 3. In-Scope Features (MVP)

### Core Functionality
- **Installation**:
  - Download and install release assets (binaries/archives) from GitHub or generic URLs.
  - Clone repositories and build from source using manifest-defined build steps.
  - Support for configurable installation directories (default: `~/.gpkg/`).
- **Package Management**:
  - Local database tracking (SQLite) for installed packages, versions, and file paths.
  - Atomic installation and uninstallation (with file tracking and cleanup).
- **Source Management**:
  - Commands to add, remove, and list package sources (local paths, HTTP URLs, GitHub organizations).
- **Updates & Maintenance**:
  - `update`: Refresh package indices and manifests from sources.
  - `upgrade`: Reinstall packages to newer versions based on the refreshed index.
- **Information & Discovery**:
  - `info <repo>`: Display repository metadata, dependencies, latest releases, and build instructions.
  - `list`: Display all currently installed packages.

### User Experience (CLI)
- Command structure consistent with standard package managers (e.g., `gpkg install`, `gpkg remove`, `gpkg search`).
- Dry-run / simulation mode (`--simulate`).
- Verbose and quiet output modes.

### Security & Integrity
- Checksum verification (SHA256/SHA512) for downloaded assets.
- Manifest validation logic to ensure provenance and required fields.

## 4. Technical Architecture
- **Language**: Go (Golang) for performance, static binaries, and cross-compilation capabilities.
- **Configuration**: TOML (e.g., `~/.gpkg/config.toml`).
- **Manifest Format**: YAML (for package definitions).
- **Data Storage**: SQLite (for the local package database).

## 5. Phasing Strategy

### Phase 1: Foundation (MVP)
- CLI skeleton and configuration loading.
- Manifest parsing (YAML) and validation.
- Basic downloader for release binaries.
- Local package database (SQLite) setup.
- Implementation of `install`, `remove`, and `list` commands.

### Phase 2: Stability & Sources
- Source manager implementation (`add-source`, `sync`).
- Update/Upgrade workflow logic.
- `info` command implementation.
- Checksum verification and resumeable downloads.

### Phase 3: Advanced Capabilities & Ecosystem
- **Advanced Security**:
  - Sandboxed builds (chroot, user namespaces) for untrusted source builds.
  - Package signing and signature verification.
  - Vulnerability scanning integration.
- **Complex Dependency Management**:
  - Transactional dependency solving with conflict resolution.
  - Channel support (stable, beta, nightly) and pinning.
- **Extended Operations**:
  - Background daemon service for auto-checking updates.
  - Binary delta updates.
  - GUI or Web frontend.
- **Ecosystem Integration**:
  - Container/OCI image export.
  - Cross-compilation helpers and build-grid integration.
  - Package groups and on-disk deduplication.


## 6. Out of Scope (Initial Release)
- Complex transactional dependency solving (SAT solving).
- Binary delta updates.
- Background daemon services.


