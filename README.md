# gpkg

gpkg is a small, user-focused package manager that installs either release binaries or builds from source into a configurable prefix. It aims to be simple, secure, and script-friendly — think "pacman for GitHub releases + source builds".

Status
- Language: Go (recommended)
- Repo: https://github.com/grave0x/gpkg
- Key saved issues:
  - Source manager: https://github.com/grave0x/gpkg/issues/2
  - Repo viewer / info: https://github.com/grave0x/gpkg/issues/6

Goals
- Install release binaries or build from source using a manifest.
- Isolated, configurable install prefix (default: ~/.gpkg).
- Source management (add/remove/list).
- Update / upgrade flows and a repo viewer (gpkg info).
- Secure downloads (checksums) and atomic installs.

Quick start (developer)
1. Clone:
   git clone https://github.com/grave0x/gpkg.git
2. Build (Go 1.20+ recommended):
   cd gpkg
   go build ./...
3. Run tests:
   go test ./...

User quick-start (when binaries exist)
- Add a source:
  gpkg add-source https://packages.example.com/index.json
- Refresh sources:
  gpkg update
- Install from release:
  gpkg install owner/repo --from-release
- Install from local manifest:
  gpkg install ./examples/manifest.yaml --from-source
- Show package info:
  gpkg info owner/repo
- Remove:
  gpkg uninstall mypkg

Commands (MVP)
- gpkg add-source <uri>
- gpkg remove-source <uri|id>
- gpkg list-sources
- gpkg update
- gpkg install <pkg|manifest> [--from-release|--from-source] [--prefix=<path>]
- gpkg upgrade [pkg...]
- gpkg uninstall <pkg>
- gpkg info <repo|pkg|manifest> [--deps-tree] [--raw|--parsed]
- gpkg list --installed

Config
- User config path: ~/.gpkg/config.toml
- Default install prefix: ~/.gpkg/
- Package DB: SQLite at ~/.gpkg/pkgdb.sqlite

Manifest (example)
See examples/manifest.yaml in this repo for a complete sample. Manifests are YAML and (for MVP) should include:
- name, version, description
- release_assets (or repo/source)
- checksums (sha256)
- build_steps (for source builds)
- files_to_install

Security & installs
- Downloads must include checksums (sha256/sha512). gpkg will verify before extracting.
- Installs are staged in a temp dir and moved atomically to the prefix on success. Installed files are recorded in the pkgdb for uninstall/rollback.

Development notes
- Keep CLI ergonomics consistent and script-friendly. Provide --json output for CI.
- Start with downloader + installer + pkgdb + gpkg info as MVP priorities.
- Use small, testable PRs. I can open or split issues for you when needed.

Contributing
See CONTRIBUTING.md for contribution guidelines and how to pick issues. For immediate tasks, check the saved issues in this Space:
- Implement source manager — https://github.com/grave0x/gpkg/issues/2
- “gpkg info <repo>” — https://github.com/grave0x/gpkg/issues/6

Roadmap (high level)
Sprint 1 (MVP): CLI, install from release, custom install dir, manifest parsing, local pkgdb, gpkg info, add/remove/list sources, update/upgrade, uninstall.
Sprint 2: Checksums, resumeable downloads, atomic installs, rollback support, manifest linter.
Sprint 3: Source builds with sandboxing, signatures, background updater, GUI/viewer.

If you want, I can:
- Draft a more detailed README (usage examples, exit codes, output formats).
- Add example manifests and a local source index for testing.
- Scaffold CLI code in Go with Cobra/urfave/cli.
Tell me which piece to produce next.