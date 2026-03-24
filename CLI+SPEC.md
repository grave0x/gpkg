# gpkg CLI specification

This document defines the initial CLI surface and UX for gpkg: subcommands, flags, output formats, config precedence, example workflows, dry-run semantics, and suggested exit codes. It is intended to be precise enough to wire the CLI (recommendation: Cobra for Go).

---

## High-level design
- CLI style: git/pacman-like subcommands (gpkg <command> [flags] [args]).
- Parsing: POSIX-style flags; both long (`--flag`) and short (`-f`) where practical.
- Subcommands are first-class; each subcommand supports `--help`.
- Libraries: implement using Cobra / urfave/cli (Cobra recommended for subcommand ergonomics).
- Behavior: sensible defaults for interactive use; explicit flags for non-interactive / CI.

---

## Global options (applies to all commands)
- -h, --help
- -V, --version: print gpkg version and exit.
- -c, --config <path>: path to config file (overrides default locations).
- --json: machine-readable JSON output where supported (info, list, install summary).
- -y, --yes: assume "yes" for prompts (non-interactive).
- --simulate, --dry-run: plan actions, do not modify filesystem/pkgdb. (See semantics section.)
- --log-level <level>: override logging level (error,warn,info,debug). (Also honor LOG_LEVEL env.)
- -v, --verbose: increase verbosity (stackable).
- --quiet: minimal output (conflicts with --verbose).
- --no-color: disable colorized output for CI.
- --offline: disallow network fetches (fail when remote data required).

Precedence note: CLI flags > environment variables > user config file > system config file.

---

## Config file locations and precedence
1. CLI `--config <path>` (highest)
2. Environment variables (e.g., GPKG_INSTALL_PREFIX, GPKG_PKGDB_PATH)
3. User config: `~/.gpkg/config.toml`
4. System config: `/etc/gpkg/config.toml` (lowest)

Default runtime values:
- install_prefix: `~/.gpkg`
- pkgdb_path: `~/.gpkg/pkgdb.sqlite`
- sources_dir: `~/.gpkg/sources.d`
- require_checksums: true

Document config keys in `config.example.toml`.

---

## Commands (MVP)

### gpkg add-source <uri>
Add a package source.
- Arguments: `<uri>` (HTTP/HTTPS JSON index URL, local directory, or `github:owner/org` shorthand).
- Flags: `--name <id>` optional id, `--priority <n>`, `--force` (overwrite).
- Output: success message or JSON `{ "added": "<id>", "uri": "<uri>" }`.

### gpkg remove-source <uri|id>
Remove a source by URI or id.
- Flags: `--yes` to skip confirmation.

### gpkg list-sources
List configured sources (shows id, uri, priority, last-seen).
- Flags: `--json`.

### gpkg update
Refresh indices from all configured sources.
- Behavior: fetch manifests/indexes and store cached copies in sources_dir/cache.
- Flags: `--concurrency N`, `--force`.
- Output: list of updated sources and counts (JSON supported).

### gpkg search <term>
Search package indices for `term`.
- Flags: `--json`, `--source <id>`.

### gpkg info <repo|pkg|manifest-path>
Show package/repo information (repo can be owner/repo, package name, or local manifest).
- Flags: `--raw` (show raw manifest), `--parsed` (structured view), `--deps-tree`.
- Output (human): summary, latest version, release assets, checksums, build steps, install methods, installed version if present.
- JSON schema: `{ "name": "...", "version": "...", "installed": { ... }, "releases": [...], "manifest": {...} }`.

### gpkg install <pkg|manifest-path> [--from-release|--from-source]
Install package.
- Positional: package identifier or local manifest path.
- Behavior:
  - Default preference: release if available; `--from-source` forces source build.
  - Downloads release asset(s) and verifies checksums before extract.
  - Source build: clone repo, run `build_steps` from manifest in a controlled env.
- Flags:
  - `--prefix <path>` override install prefix for this operation.
  - `--arch <arch>` select platform-specific asset.
  - `--force` overwrite existing install/version.
  - `--no-deps` skip dependency install (dangerous).
  - `--simulate` / `--dry-run` (see semantics).
- Output:
  - Human: step-by-step progress and final success message.
  - JSON: `{ "package":"...", "version":"...", "status":"installed", "files":[...], "prefix":"/home/user/.gpkg" }`.
- Side effects:
  - On success, record package metadata & file list in pkgdb.

### gpkg upgrade [pkg...]
Upgrade packages.
- Without args: upgrade all installed packages where newer version available.
- With args: upgrade only given packages.
- Flags: `--simulate`, `--yes`, `--concurrency`.
- Output: list of upgraded packages and status.

### gpkg uninstall <pkg> [--purge]
Uninstall package and remove recorded files.
- Flags: `--purge` remove all data (pkgdb record and caches).
- Safety: will refuse to remove files outside install_prefix. Confirm unless `--yes`.
- Output: removed files list.

### gpkg list [--installed|--available]
List installed packages (default) or available packages (from sources).
- Flags: `--json`, `--filter`, `--sort`.

### gpkg config (subcommands)
- `gpkg config get <key>`
- `gpkg config set <key> <value>`
- `gpkg config show` (prints merged config: system < user < CLI)

### gpkg validate <manifest>
Validate a manifest against schema.
- Output: success or list of errors/warnings. `--fix` can auto-fix trivial issues.

### gpkg rollback <pkg> --to-version <ver>
Roll back package to prior installed version (if available).
- Requires prior versions recorded in pkgdb.

---

## UX and help output
- Each subcommand exposes `--help` with synopsis, examples, and common flags.
- Use examples in help (2-3 short examples).
- `gpkg --help` lists commands and global options; `gpkg <cmd> --help` lists command options and examples.
- Provide human-friendly formatting by default (colors, aligned columns), and `--json` for CI.

Suggested help snippet (install):
gpkg install <pkg|manifest> [--from-release|--from-source] [--prefix PATH] [--simulate] [-y]
Install a package from a release or build from source.
Examples:
  gpkg install owner/cool-tool
  gpkg install ./examples/cool-tool.yaml --from-source --prefix /opt/gpkg

---

## Output formats: human and JSON
- Use `--json` to produce structured output on stdout; errors still go to stderr with JSON error structure when possible.
- Error JSON schema:
  { "error": { "code": <int>, "message": "<string>", "detail": "<optional>" } }
- Success JSON examples:
  - install:
    {
      "package":"cool-tool",
      "version":"1.2.0",
      "status":"installed",
      "files":["bin/cool-tool","share/man/..."],
      "prefix":"/home/alice/.gpkg"
    }
  - info:
    { "name":"cool-tool", "version":"1.2.0", "installed": {"version":"1.1.0"}, "releases":[ ... ] }

---

## Dry-run / simulate semantics
`--simulate` / `--dry-run` should:
- Validate inputs and manifests.
- Resolve dependencies and produce a full plan of actions (downloads, builds, file changes).
- Perform read-only network operations necessary to inspect manifests and release metadata (unless `--offline`).
- NOT download large release assets, NOT extract to install prefix, NOT update pkgdb, NOT write persistent caches except for safe read-only caches.
- Return a non-zero exit if plan reveals fatal issues (checksum mismatch in index, missing manifest fields, unresolved deps).
- Display a clear "Plan" summary: actions with their type (download, extract, build, write), sizes, and estimated time.

---

## Command workflows (step-by-step)

Install (release-first)
1. gpkg resolves package identifier via sources.
2. Fetch manifest (remote or local).
3. Validate manifest schema.
4. Select release asset for platform/arch.
5. Download asset (resume support), verify checksum.
6. Extract asset into staging dir.
7. Record list of installed files.
8. Atomically move staged tree into install_prefix/<pkg>-<version> (or layout decided).
9. Update pkgdb with metadata and files.
10. Report success.

Install (from source)
1. Fetch manifest and validate.
2. Clone repo (or use local path).
3. Set up build env (templated variables).
4. Run build steps (capture logs).
5. Stage installation artifacts.
6. Continue as release install (atomically swap, pkgdb update).

Info
1. Resolve input (owner/repo or local manifest).
2. Fetch manifest/release metadata.
3. Display manifest summary, latest release, install status, deps tree.

Update/Upgrade
- `update` refreshes indices.
- `upgrade` resolves new versions and runs install workflow for updated packages.

Uninstall
1. Lookup pkgdb for file list.
2. Confirm removal and remove files (safeguard to not remove outside prefix).
3. Update pkgdb.

---

## Suggested exit codes
- 0: success
- 1: general failure
- 2: usage/argument error
- 3: network error / remote unavailable
- 4: checksum/verification failed
- 5: install/build failed
- 6: manifest validation error
- 7: package not found
- 8: pkgdb error
- 9: simulated-only (e.g., dry-run success distinct) — avoid if possible; prefer 0

---

## Security & safety notes (UX)
- By default require checksums; warn and abort if missing.
- Confirm destructive operations (uninstall, remove-source) unless `--yes`.
- Build-from-source should warn when provenance is missing and respect `allow_unverified_source_builds` config.

---

## Examples
- Install a release binary:
  gpkg install owner/cool-tool
- Force a source build:
  gpkg install owner/cool-tool --from-source
- Install a local manifest:
  gpkg install ./examples/cool-tool.yaml --prefix /opt/gpkg
- Show info:
  gpkg info owner/cool-tool --deps-tree
- Simulate an upgrade:
  gpkg upgrade --simulate
- Add a source:
  gpkg add-source https://packages.example.com/index.json --name example

---

## Implementation recommendations
- Use Cobra for subcommands and well-structured help text.
- Centralize config merging (system < user < env < CLI).
- Create a "planner" component that both simulate and real installers consume (same plan generation).
- Provide JSON schemas for manifests and for gpkg's JSON outputs.
- Add extensive unit tests for plan generation and dry-run behavior.

---

End of CLI spec. Implementing this will allow the project to wire subcommands, create consistent help, and produce machine-friendly outputs for CI automation.
