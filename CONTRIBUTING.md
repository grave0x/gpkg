# Contributing to gpkg

Thanks for your interest in contributing to gpkg! We welcome contributions of all kinds: bug fixes, new features, documentation improvements, and more.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Workflow](#workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Code Review](#code-review)
- [Issue Workflow](#issue-workflow)
- [Communication](#communication)

---

## Getting Started

### Prerequisites

- **Go 1.20+** (see [go.mod](go.mod) for exact version)
- **SQLite3** (required for CGO)
- **Git**
- **Make** (optional, for convenience commands)

### Fork and Clone

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gpkg.git
   cd gpkg
   ```
3. **Add upstream** remote:
   ```bash
   git remote add upstream https://github.com/grave0x/gpkg.git
   ```

---

## Development Setup

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed setup instructions, including:
- Building from source
- Running tests locally
- Setting up your development environment
- IDE configuration

**Quick start:**
```bash
# Install dependencies
go mod download

# Build the binary
go build -o bin/gpkg .

# Run tests
go test ./...

# Run linter
golangci-lint run
```

---

## Workflow

### Branch Naming

Use descriptive branch names with prefixes:
- `feat/` - New features (e.g., `feat/add-retry-logic`)
- `fix/` - Bug fixes (e.g., `fix/checksum-validation`)
- `docs/` - Documentation updates (e.g., `docs/update-install-guide`)
- `refactor/` - Code refactoring (e.g., `refactor/cleanup-resolver`)
- `test/` - Test improvements (e.g., `test/add-integration-tests`)
- `chore/` - Maintenance tasks (e.g., `chore/update-dependencies`)

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Examples:**
```
feat(install): add retry logic for failed downloads

Implements exponential backoff for network errors during
package downloads. Configurable via --max-retries flag.

Closes #123
```

```
fix(checksum): handle SHA-512 checksums correctly

SHA-512 checksums were being truncated. This fix ensures
full checksum comparison for all supported algorithms.

Fixes #456
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `refactor` - Code refactoring (no behavior change)
- `test` - Adding or updating tests
- `chore` - Maintenance (dependencies, CI, etc.)
- `perf` - Performance improvements

### Keeping Your Fork Updated

```bash
# Fetch upstream changes
git fetch upstream

# Rebase your branch on latest main
git rebase upstream/main

# Push to your fork (may need force push)
git push origin feat/my-feature --force-with-lease
```

---

## Coding Standards

### Go Code Style

- **Use `gofmt`** - All code must be formatted with `gofmt`
- **Run `golangci-lint`** - Passes all linters in CI
- **Follow Go best practices** - See [Effective Go](https://golang.org/doc/effective_go.html)

### Code Organization

- **Internal packages**: Business logic goes in `internal/`
- **Commands**: CLI commands in `cmd/gpkg/cmd/`
- **Tests**: Co-located with source files (`*_test.go`)
- **Examples**: Manifest examples in `examples/`

### Error Handling

- Use **wrapped errors** with context: `fmt.Errorf("failed to download: %w", err)`
- Return errors, don't panic (except for programmer errors)
- Log at appropriate levels (debug, info, warn, error)

### Comments

- Document **exported** functions, types, and constants
- Use **godoc** format for documentation comments
- Include **examples** for complex APIs

### CLI Design

- **JSON output** for machine consumption (`--json` flag)
- **Non-interactive mode** for scripts (no prompts)
- **Dry-run support** (`--dry-run` flag where applicable)
- **Consistent flag naming** across commands
- **Clear error messages** with actionable suggestions

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with race detection (as CI does)
go test -race ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/resolver/...

# Run specific test
go test -run TestInstallBinary ./cmd/gpkg/cmd/
```

### Test Requirements

- **Unit tests** for all new logic
- **Integration tests** for end-to-end workflows
- **Table-driven tests** for multiple scenarios
- **Mocks** for external dependencies (HTTP, filesystem)

### Writing Good Tests

```go
func TestInstallBinary(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid binary install",
            input:   "example-package",
            want:    "installed successfully",
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

---

## Documentation

### What to Document

- **README.md** - Update for new features or changes
- **Wiki** - Comprehensive guides in `wiki/` directory
- **Code comments** - Godoc for exported APIs
- **Examples** - Add manifest examples for new features
- **Changelog** - Add entry for significant changes

### Wiki Structure

```
wiki/
├── Home.md
├── Getting-Started.md
├── Troubleshooting.md
├── User-Guide/
│   ├── Installation.md
│   ├── Basic-Commands.md
│   └── ...
├── Developer-Guide/
│   ├── Contributing.md
│   ├── Architecture.md
│   └── ...
└── Reference/
    ├── CLI-Reference.md
    ├── API-Reference.md
    └── ...
```

### Documentation Standards

- Use **clear, concise language**
- Include **code examples** where helpful
- Add **screenshots** for UI/visual features
- Keep **table of contents** updated
- Link **related documentation**

---

## Pull Request Process

### Before Submitting

1. ✅ **Tests pass** locally (`go test ./...`)
2. ✅ **Linter passes** (`golangci-lint run`)
3. ✅ **Documentation updated** (if applicable)
4. ✅ **Commits are clean** and follow conventions
5. ✅ **Branch is up-to-date** with `main`

### Creating a Pull Request

1. **Push** your branch to your fork
2. **Open a PR** against `grave0x/gpkg:main`
3. **Fill out the PR template** completely
4. **Link related issues** (e.g., "Closes #123")
5. **Request review** from maintainers

### PR Title Format

Use the same format as commit messages:

```
feat(install): add retry logic for failed downloads
fix(checksum): handle SHA-512 checksums correctly
docs(wiki): update installation guide for macOS
```

### PR Description

Include:
- **What** - What does this PR do?
- **Why** - Why is this change needed?
- **How** - How does it work?
- **Testing** - How was it tested?
- **Screenshots** - For UI changes
- **Breaking changes** - If any

**Example:**
```markdown
## What
Adds retry logic for failed package downloads with exponential backoff.

## Why
Network errors during downloads cause installation failures. Users have to manually retry.

## How
- Implements `RetryDownload()` with exponential backoff
- Configurable max retries (default: 3)
- New flag: `--max-retries`

## Testing
- Added unit tests for retry logic
- Tested manually with flaky network
- All existing tests pass

## Breaking Changes
None

Closes #123
```

---

## Code Review

### Review Process

1. **Automated checks** run first (tests, lints, security)
2. **Maintainer review** - Code quality and design
3. **Address feedback** - Make requested changes
4. **Approval** - PR is approved and merged

### Addressing Feedback

- **Be responsive** - Reply to comments promptly
- **Ask questions** - If feedback is unclear
- **Make changes** - Push additional commits
- **Re-request review** - After addressing all feedback

### Review Timeline

- Initial review: **Within 2-3 business days**
- Follow-up reviews: **Within 1-2 business days**
- Urgent fixes: **Same day** (if possible)

---

## Issue Workflow

### Finding Issues

- Browse [open issues](https://github.com/grave0x/gpkg/issues)
- Look for `good-first-issue` label for beginner-friendly tasks
- Look for `help-wanted` label for issues seeking contributors

### Claiming an Issue

1. **Comment** on the issue expressing interest
2. **Wait for assignment** (or self-assign if allowed)
3. **Ask questions** if anything is unclear
4. **Start working** after assignment

### Creating Issues

**Bug Reports** should include:
- Clear title and description
- Steps to reproduce
- Expected vs actual behavior
- Environment (OS, Go version, gpkg version)
- Logs or error messages

**Feature Requests** should include:
- Clear use case
- Proposed solution (if any)
- Alternatives considered
- Willingness to contribute

---

## Communication

### Where to Communicate

- **GitHub Issues** - Bug reports, feature requests
- **Pull Requests** - Code review and discussion
- **Discussions** - General questions and ideas

### Communication Guidelines

- **Be respectful** - Follow the [Code of Conduct](CODE_OF_CONDUCT.md)
- **Be concise** - Get to the point quickly
- **Be clear** - Avoid ambiguity
- **Be patient** - Maintainers are volunteers

### Getting Help

- **Check documentation** - README, wiki, API docs
- **Search issues** - Your question may already be answered
- **Ask in Discussions** - For general questions
- **Open an issue** - For specific problems

---

## Additional Resources

- **[Development Guide](DEVELOPMENT.md)** - Detailed development setup
- **[Architecture Guide](wiki/Developer-Guide/Architecture.md)** - System design
- **[Testing Guide](wiki/Developer-Guide/Testing.md)** - Testing practices
- **[Release Process](wiki/Developer-Guide/Release-Process.md)** - How releases work
- **[Workflow Documentation](.github/workflows/README.md)** - CI/CD workflows

---

## Recognition

Contributors are recognized in:
- **Release notes** - Credited for their contributions
- **CONTRIBUTORS.md** - Hall of fame
- **Commit history** - Co-authored-by trailers

---

## License

By contributing to gpkg, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

**Thank you for contributing to gpkg! 🎉**

If you want help picking a starter issue, open a discussion or comment on an issue marked `good-first-issue`.