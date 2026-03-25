# GitHub Actions Workflows

This directory contains all CI/CD workflows for the gpkg project. These workflows automate testing, code quality, security scanning, releases, and documentation.

---

## 📋 Table of Contents

- [Overview](#overview)
- [Workflows](#workflows)
  - [Test Workflow](#test-workflow)
  - [Code Quality Workflow](#code-quality-workflow)
  - [Security Workflow](#security-workflow)
  - [Dependency Management Workflow](#dependency-management-workflow)
  - [Release Workflow](#release-workflow)
  - [Wiki Sync Workflow](#wiki-sync-workflow)
- [Triggers](#triggers)
- [Secrets and Permissions](#secrets-and-permissions)
- [Manual Workflow Execution](#manual-workflow-execution)
- [Troubleshooting](#troubleshooting)

---

## Overview

The gpkg project uses GitHub Actions for continuous integration and deployment. All workflows are designed to:

- ✅ Run automatically on relevant code changes
- ✅ Support manual triggering for testing
- ✅ Provide clear, actionable feedback
- ✅ Fail fast to catch issues early
- ✅ Generate artifacts and reports

**Workflow Status Badges:**
- [![Tests](https://github.com/grave0x/gpkg/workflows/Test/badge.svg)](https://github.com/grave0x/gpkg/actions/workflows/test.yml)
- [![Code Quality](https://github.com/grave0x/gpkg/workflows/Code%20Quality/badge.svg)](https://github.com/grave0x/gpkg/actions/workflows/code-quality.yml)
- [![Security](https://github.com/grave0x/gpkg/workflows/Security/badge.svg)](https://github.com/grave0x/gpkg/actions/workflows/security.yml)

---

## Workflows

### Test Workflow

**File:** [`test.yml`](test.yml)

**Purpose:** Comprehensive testing across platforms and Go versions

**Triggers:**
- Push to `main` branch
- Pull requests to `main`
- Manual dispatch

**Jobs:**
1. **Unit Tests** - Run on matrix of OS (Linux, macOS, Windows) and Go versions (1.20, 1.21, 1.22)
   - Executes all unit tests with race detection
   - Generates coverage reports
   - Uploads coverage to Codecov (Ubuntu + Go 1.22 only)
   
2. **Linting** - Run golangci-lint with comprehensive checks
   - Version: golangci-lint v6
   - Timeout: 10 minutes
   - Caching enabled for performance
   
3. **Coverage Report** - Generate coverage summary for PRs
   - Runs on latest Go 1.22
   - Comments on PRs with coverage percentage
   - Includes per-package breakdown
   
4. **Benchmarks** - Performance benchmarking
   - Runs benchmarks with `-benchmem`
   - Compares with previous results on PRs
   - Stores results as artifacts
   
5. **Integration Tests** - Full end-to-end testing
   - Builds gpkg binary
   - Runs integration test suite
   - Tests on Linux and macOS
   
6. **Build All Platforms** - Verify multi-platform builds
   - Builds for 7 platform combinations
   - Ensures CGO compatibility with SQLite

**Artifacts:**
- Coverage reports (HTML and text)
- Benchmark results
- Test logs

---

### Code Quality Workflow

**File:** [`code-quality.yml`](code-quality.yml)

**Purpose:** Enforce code quality standards beyond linting

**Triggers:**
- Push to `main` branch
- Pull requests to `main`
- Manual dispatch

**Checks:**
1. **gofmt** - Verify code formatting
2. **goimports** - Check import organization
3. **go vet** - Static analysis for suspicious constructs
4. **staticcheck** - Advanced static analysis
5. **ineffassign** - Detect ineffectual assignments
6. **misspell** - Spell checking in code and comments
7. **Quality Metrics** - LOC, cyclomatic complexity
8. **Quality Report** - Generate markdown report

**Artifacts:**
- Quality report (markdown)
- Metrics summary

---

### Security Workflow

**File:** [`security.yml`](security.yml)

**Purpose:** Continuous security scanning and vulnerability detection

**Triggers:**
- Push to `main` branch
- Pull requests to `main`
- Weekly schedule (Mondays at midnight UTC)
- Manual dispatch

**Scanners:**
1. **GoSec** - Go security scanner
   - Detects common security issues
   - Checks for hardcoded credentials
   - Analyzes SQL injection risks
   
2. **Trivy** - Multi-purpose security scanner
   - Vulnerability scanning
   - Secret detection
   - Misconfiguration detection
   - License compliance
   
3. **Dependency Vulnerability Check** - Using govulncheck
   - Checks Go dependencies for known vulnerabilities
   - Uses official Go vulnerability database
   
4. **GitHub CodeQL** - Advanced semantic analysis
   - Detects complex security patterns
   - Customizable query sets

**Reports:**
- SARIF reports uploaded to GitHub Security tab
- Artifacts stored for 30 days
- Security summary with all scan results

---

### Dependency Management Workflow

**File:** [`dependency-update.yml`](dependency-update.yml)

**Purpose:** Automate dependency updates and maintenance

**Triggers:**
- Weekly schedule (Mondays at 10:00 UTC)
- Pull requests modifying `go.mod` or `go.sum`
- Manual dispatch

**Jobs:**
1. **Check Outdated** - Find outdated dependencies
   - Lists available updates
   - Creates report artifact
   - Comments on PRs with outdated info
   
2. **Generate Changelog** - Document dependency changes
   - Runs on dependency PRs
   - Shows version diffs
   - Lists all updated modules
   
3. **Auto-merge Minor/Patch** - Automatically merge safe updates
   - Only merges minor/patch updates (not major)
   - Runs full test suite first
   - Auto-approves and merges if tests pass
   - Requires `dependabot[bot]` actor
   
4. **Security Audit** - Verify dependency integrity
   - Runs `go mod verify`
   - Executes govulncheck
   - Generates security report

**Artifacts:**
- Outdated dependencies report
- Dependency changelog
- Security audit report

---

### Release Workflow

**File:** [`release.yml`](release.yml)

**Purpose:** Automated multi-platform releases with assets

**Triggers:**
- Git tags matching `v*.*.*` (e.g., `v0.1.0`)
- Manual dispatch with tag input

**Jobs:**
1. **Build and Release** - Main release job
   - Runs full test suite
   - Generates changelog from git commits
   - Builds binaries for 7 platforms:
     - Linux: amd64, arm64, 386
     - macOS: amd64 (Intel), arm64 (Apple Silicon)
     - Windows: amd64, 386
   - Calculates SHA256 checksums
   - Creates archives (`.tar.gz` for Unix, `.zip` for Windows)
   - Generates comprehensive release notes
   - Creates GitHub release with all assets
   - Updates `latest` tag for stable releases
   
2. **Update Homebrew** - Homebrew formula automation
   - Downloads release archives
   - Calculates SHA256 checksums
   - Generates updated Homebrew formula
   - Creates PR with formula update
   - Includes installation instructions
   
3. **Build and Publish Docker** - Docker images
   - Builds multi-arch images (amd64, arm64)
   - Tags with semantic versioning
   - Publishes to GitHub Container Registry
   - Creates `latest` tag for stable releases

**Artifacts:**
- Release binaries (7 platforms)
- Checksums file
- Release notes
- Updated Homebrew formula
- Docker images

**Release Assets:**
```
gpkg-v0.1.0-linux-amd64.tar.gz
gpkg-v0.1.0-linux-arm64.tar.gz
gpkg-v0.1.0-linux-386.tar.gz
gpkg-v0.1.0-darwin-amd64.tar.gz
gpkg-v0.1.0-darwin-arm64.tar.gz
gpkg-v0.1.0-windows-amd64.zip
gpkg-v0.1.0-windows-386.zip
checksums.txt
```

---

### Wiki Sync Workflow

**File:** [`wiki-sync.yml`](wiki-sync.yml)

**Purpose:** Automatically sync documentation from repository to GitHub Wiki

**Triggers:**
- Push to `main` branch affecting `wiki/**` files
- Changes to workflow file
- Manual dispatch

**Process:**
1. Checkout main repository
2. Checkout wiki repository
3. Sync all wiki files from `wiki/` directory
4. Commit and push changes to wiki repository
5. Validate markdown files
6. Generate wiki statistics

**Features:**
- Preserves wiki git history
- Includes source commit information
- Validates wiki content
- Provides detailed sync summary

---

## Triggers

### Automatic Triggers

| Workflow | Push to main | Pull Request | Schedule | Tags |
|----------|--------------|--------------|----------|------|
| Test | ✅ | ✅ | ❌ | ❌ |
| Code Quality | ✅ | ✅ | ❌ | ❌ |
| Security | ✅ | ✅ | Weekly | ❌ |
| Dependencies | ❌ | ✅ (go.mod) | Weekly | ❌ |
| Release | ❌ | ❌ | ❌ | ✅ |
| Wiki Sync | ✅ (wiki/) | ❌ | ❌ | ❌ |

### Scheduled Triggers

- **Security Workflow**: Every Monday at 00:00 UTC
- **Dependency Management**: Every Monday at 10:00 UTC

---

## Secrets and Permissions

### Required Permissions

All workflows use `GITHUB_TOKEN` with appropriate permissions:

- **Test Workflow**: `contents: read`
- **Code Quality**: `contents: read`
- **Security**: `contents: read, security-events: write`
- **Dependencies**: `contents: write, pull-requests: write, issues: write`
- **Release**: `contents: write, packages: write`
- **Wiki Sync**: `contents: write`

### Optional Secrets

- **`CODECOV_TOKEN`**: For uploading coverage to Codecov (optional, works without for public repos)

### No Additional Secrets Needed

All workflows use `GITHUB_TOKEN` which is automatically provided by GitHub Actions. No manual secret configuration required!

---

## Manual Workflow Execution

All workflows support manual triggering via `workflow_dispatch`. To run manually:

### Via GitHub UI

1. Navigate to **Actions** tab
2. Select the workflow from the left sidebar
3. Click **Run workflow** button
4. Select branch (usually `main`)
5. Fill in any required inputs (e.g., tag for release)
6. Click **Run workflow**

### Via GitHub CLI

```bash
# Run test workflow
gh workflow run test.yml

# Run code quality checks
gh workflow run code-quality.yml

# Run security scan
gh workflow run security.yml

# Sync wiki
gh workflow run wiki-sync.yml

# Create a release (requires tag input)
gh workflow run release.yml -f tag=v0.2.0

# Check workflow status
gh run list --workflow=test.yml

# View workflow run details
gh run view <run-id>
```

---

## Troubleshooting

### Common Issues

#### Test Failures

**Problem:** Tests fail in CI but pass locally

**Solutions:**
1. Check Go version consistency
2. Verify CGO is enabled (`CGO_ENABLED=1`)
3. Ensure SQLite3 is available
4. Check for race conditions (tests use `-race` flag)

```bash
# Run tests with same settings as CI
CGO_ENABLED=1 go test -v -race ./...
```

#### Lint Failures

**Problem:** Linting passes locally but fails in CI

**Solutions:**
1. Update golangci-lint to match CI version
2. Run with same configuration:

```bash
# Install matching version
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run with same args as CI
golangci-lint run --verbose --print-issued-lines=true --print-linter-name=true
```

#### Release Failures

**Problem:** Release workflow fails to build for specific platform

**Solutions:**
1. Verify CGO is enabled for the platform
2. Check cross-compilation requirements
3. Test locally with matching GOOS/GOARCH:

```bash
# Example: Build for macOS ARM64 from Linux
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o gpkg-darwin-arm64 .
```

#### Homebrew Formula Update Fails

**Problem:** Homebrew formula PR creation fails

**Solutions:**
1. Verify release assets were uploaded successfully
2. Check that checksums.txt contains all platforms
3. Ensure `GITHUB_TOKEN` has write permissions

#### Docker Build Fails

**Problem:** Multi-arch Docker build fails

**Solutions:**
1. Check Dockerfile exists and is valid
2. Verify QEMU setup for cross-platform builds
3. Check GitHub Container Registry permissions

### Workflow Debugging

Enable debug logging:

```bash
# Set repository secrets
gh secret set ACTIONS_RUNNER_DEBUG --body "true"
gh secret set ACTIONS_STEP_DEBUG --body "true"
```

Then re-run the workflow to see detailed logs.

### Getting Help

- **Workflow logs**: Check the Actions tab for detailed execution logs
- **GitHub Issues**: Report workflow issues at https://github.com/grave0x/gpkg/issues
- **GitHub Actions docs**: https://docs.github.com/en/actions

---

## Workflow Maintenance

### Updating Workflows

When modifying workflows:

1. **Test locally first** using `act` (if possible):
   ```bash
   # Install act
   brew install act  # macOS
   # or download from https://github.com/nektos/act
   
   # Run workflow locally
   act -W .github/workflows/test.yml
   ```

2. **Use workflow_dispatch** for testing changes in CI

3. **Update this documentation** when making significant changes

4. **Review security implications** of any permission changes

### Keeping Actions Up to Date

Dependabot automatically creates PRs for action updates. Review and merge these PRs regularly.

To manually update:

```bash
# Check for outdated actions
gh api repos/:owner/:repo/actions/workflows | jq '.workflows[].path'

# Review and update action versions in workflow files
```

---

## Related Documentation

- **[Contributing Guide](../../CONTRIBUTING.md)** - How to contribute to gpkg
- **[Development Guide](../../DEVELOPMENT.md)** - Development setup and practices
- **[Release Process](../../wiki/Developer-Guide/Release-Process.md)** - How releases are managed
- **[Security Policy](../../SECURITY.md)** - Security practices and reporting

---

*Workflows documentation last updated for gpkg v0.1.0*
