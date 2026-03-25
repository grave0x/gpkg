# Release Process

This guide documents the release process for gpkg, including versioning, building releases, and publishing.

## Table of Contents

- [Versioning Strategy](#versioning-strategy)
- [Release Preparation](#release-preparation)
- [Building Releases](#building-releases)
- [Release Workflow](#release-workflow)
- [Automated Release Process](#automated-release-process)
- [Post-Release Tasks](#post-release-tasks)
- [Hotfix Releases](#hotfix-releases)

## Versioning Strategy

gpkg follows [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH
```

- **MAJOR**: Breaking changes, incompatible API changes
- **MINOR**: New features, backward-compatible
- **PATCH**: Bug fixes, backward-compatible

### Version Examples

- `1.0.0` - Initial stable release
- `1.1.0` - New feature added
- `1.1.1` - Bug fix
- `2.0.0` - Breaking change

### Pre-release Versions

For development and testing:

- `1.0.0-alpha.1` - Alpha release
- `1.0.0-beta.1` - Beta release
- `1.0.0-rc.1` - Release candidate

## Release Preparation

### 1. Version Planning

Determine the next version number based on changes:

```bash
# Review commits since last release
git log $(git describe --tags --abbrev=0)..HEAD --oneline

# Categorize changes
# - Breaking changes → MAJOR bump
# - New features → MINOR bump  
# - Bug fixes → PATCH bump
```

### 2. Update Version

Update version in relevant files:

```bash
# main.go
version = "1.2.0"

# README.md
# Update version numbers in examples

# CHANGELOG.md (create if doesn't exist)
# Add new version section
```

### 3. Update Changelog

Create/update `CHANGELOG.md`:

```markdown
## [1.2.0] - 2024-03-25

### Added
- New package signing feature
- Support for private registries

### Changed
- Improved checksum verification performance
- Updated dependency resolution algorithm

### Fixed
- Fixed database lock timeout issue
- Corrected manifest validation error messages

### Security
- Updated dependencies to patch vulnerabilities
```

### 4. Run Full Test Suite

```bash
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Check test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run integration tests
go test ./tests/integration/...

# Run linters
golangci-lint run
```

### 5. Build and Test Locally

```bash
# Build all platforms using release script
./scripts/build.sh 1.2.0

# Test binaries
./dist/gpkg-1.2.0-linux-amd64 --version
./dist/gpkg-1.2.0-darwin-amd64 --version

# Verify checksums
cd dist && sha256sum -c checksums.txt
```

## Building Releases

### Using the Build Script

The `scripts/build.sh` script builds for all platforms:

```bash
# Build version 1.2.0
./scripts/build.sh 1.2.0

# Output in dist/
# - gpkg-1.2.0-linux-amd64
# - gpkg-1.2.0-linux-arm64
# - gpkg-1.2.0-darwin-amd64
# - gpkg-1.2.0-darwin-arm64
# - gpkg-1.2.0-windows-amd64.exe
# - checksums.txt
```

### Manual Cross-Platform Build

```bash
# Set version
VERSION="1.2.0"

# Linux AMD64
GOOS=linux GOARCH=amd64 go build \
  -ldflags "-X main.version=${VERSION} -s -w" \
  -trimpath \
  -o dist/gpkg-${VERSION}-linux-amd64 .

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build \
  -ldflags "-X main.version=${VERSION} -s -w" \
  -trimpath \
  -o dist/gpkg-${VERSION}-darwin-arm64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build \
  -ldflags "-X main.version=${VERSION} -s -w" \
  -trimpath \
  -o dist/gpkg-${VERSION}-windows-amd64.exe .
```

### Generate Checksums

```bash
cd dist
sha256sum gpkg-${VERSION}-* > checksums.txt
```

### Create Archives

```bash
# Linux/macOS - tar.gz
for file in gpkg-${VERSION}-linux-* gpkg-${VERSION}-darwin-*; do
  tar czf ${file}.tar.gz ${file}
done

# Windows - zip
for file in gpkg-${VERSION}-windows-*.exe; do
  zip ${file%.exe}.zip ${file}
done
```

## Release Workflow

### Using GitHub Actions (Automated)

The automated release process is triggered by creating a tag:

```bash
# Create and push tag
git tag -a v1.2.0 -m "Release version 1.2.0"
git push origin v1.2.0
```

GitHub Actions workflow (`.github/workflows/release.yml`) will:
1. Checkout code
2. Set up Go environment
3. Run tests
4. Build multi-platform binaries
5. Generate checksums
6. Create archives
7. Generate release notes
8. Create GitHub Release
9. Upload all artifacts

### Manual Release Process

If automation is unavailable:

#### 1. Create Git Tag

```bash
# Create annotated tag
git tag -a v1.2.0 -m "Release version 1.2.0

## Changes in this release

### Added
- Feature 1
- Feature 2

### Fixed
- Bug fix 1
- Bug fix 2
"

# Push tag
git push origin v1.2.0
```

#### 2. Build Release Binaries

```bash
./scripts/build.sh 1.2.0
```

#### 3. Create GitHub Release

Via GitHub web interface:
1. Go to repository → Releases → Draft a new release
2. Choose tag: `v1.2.0`
3. Release title: `Release v1.2.0`
4. Description: Copy from CHANGELOG.md
5. Upload binaries from `dist/`:
   - All `*.tar.gz` and `*.zip` files
   - `checksums.txt`
6. Publish release

Via GitHub CLI:

```bash
gh release create v1.2.0 \
  --title "Release v1.2.0" \
  --notes-file CHANGELOG.md \
  dist/*.tar.gz \
  dist/*.zip \
  dist/checksums.txt
```

## Automated Release Process

The GitHub Actions release workflow handles:

### Trigger Events

```yaml
# Tag push (v*.*.*)
on:
  push:
    tags:
      - 'v*.*.*'
  
  # Manual dispatch
  workflow_dispatch:
    inputs:
      tag:
        description: 'Tag to release'
        required: true
```

### Build Matrix

Builds for all supported platforms:
- Linux: amd64, arm64, 386
- macOS: amd64, arm64
- Windows: amd64, 386

### Release Assets

Automatically generates:
- Binary archives (`.tar.gz` for Unix, `.zip` for Windows)
- Checksum file (`checksums.txt`)
- Release notes with download links

## Post-Release Tasks

### 1. Verify Release

```bash
# Download and verify
curl -LO https://github.com/grave0x/gpkg/releases/download/v1.2.0/gpkg-1.2.0-linux-amd64.tar.gz
curl -LO https://github.com/grave0x/gpkg/releases/download/v1.2.0/checksums.txt

# Verify checksum
sha256sum -c checksums.txt

# Test binary
tar xzf gpkg-1.2.0-linux-amd64.tar.gz
./gpkg-1.2.0-linux-amd64 --version
```

### 2. Update Documentation

- Update README.md installation instructions
- Update wiki if necessary
- Announce in GitHub Discussions

### 3. Announcement

Consider announcing on:
- GitHub Discussions
- Project website
- Social media
- Relevant communities

### 4. Monitor Issues

Watch for issues related to the new release:
- Installation problems
- Regression bugs
- Platform-specific issues

## Hotfix Releases

For critical bugs in production releases:

### 1. Create Hotfix Branch

```bash
# From the release tag
git checkout -b hotfix/1.2.1 v1.2.0
```

### 2. Apply Fix

```bash
# Make minimal changes
# Only fix the critical bug
git commit -m "fix: critical bug in package installation"
```

### 3. Test Thoroughly

```bash
go test ./...
./scripts/build.sh 1.2.1
# Test affected functionality
```

### 4. Merge and Tag

```bash
# Merge to main
git checkout main
git merge hotfix/1.2.1

# Tag hotfix
git tag -a v1.2.1 -m "Hotfix release 1.2.1"
git push origin main v1.2.1

# Also merge back to develop if using git-flow
git checkout develop
git merge hotfix/1.2.1
git push origin develop
```

### 5. Release Hotfix

Follow the normal release process for the hotfix tag.

## Release Checklist

Use this checklist for each release:

### Pre-Release
- [ ] All tests passing
- [ ] Linters passing
- [ ] Version number updated
- [ ] CHANGELOG.md updated
- [ ] Documentation updated
- [ ] Local build successful
- [ ] Binaries tested locally

### Release
- [ ] Git tag created
- [ ] Tag pushed to GitHub
- [ ] GitHub Actions workflow succeeded
- [ ] Release created on GitHub
- [ ] All artifacts uploaded
- [ ] Release notes published

### Post-Release
- [ ] Release verified (download + test)
- [ ] Documentation updated
- [ ] Announcement made
- [ ] No critical issues reported

## Troubleshooting Releases

### Build Fails

```bash
# Check Go version
go version

# Clean build cache
go clean -cache

# Rebuild
go build -v .
```

### GitHub Actions Fails

1. Check workflow logs in GitHub Actions tab
2. Re-run failed jobs if transient error
3. Fix issues and create new tag if needed

### Checksum Mismatch

```bash
# Regenerate checksums
cd dist
rm checksums.txt
sha256sum gpkg-* > checksums.txt
```

### Missing Assets

Ensure all platforms are built:
- scripts/build.sh builds all platforms
- GitHub Actions workflow includes all platforms

## See Also

- [Development Setup](Development-Setup) - Set up dev environment
- [Testing](Testing) - Testing guidelines
- [Contributing](Contributing) - Contribution process
- [GitHub Actions](../../.github/workflows/README.md) - CI/CD workflows

---

**Questions?** Open an issue or discussion on GitHub.
