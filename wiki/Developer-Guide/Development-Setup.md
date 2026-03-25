# Development Setup

A comprehensive step-by-step guide for setting up your development environment to contribute to gpkg.

## Prerequisites

Before you begin, ensure you have the following tools installed on your system:

### Required Software

#### 1. Go 1.20 or Higher
gpkg requires Go 1.20+ for building and development.

**Download:** [https://golang.org/dl](https://golang.org/dl)

**Verify installation:**
```bash
go version
# Expected output: go version go1.20.x or higher
```

**Set up Go environment:**
```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:$(go env GOPATH)/bin
export GOPATH=$(go env GOPATH)
```

#### 2. Git
Git is required for cloning the repository and version control.

**Install:**
```bash
# Linux (Debian/Ubuntu)
sudo apt-get install git

# Linux (Fedora/RHEL)
sudo dnf install git

# macOS
brew install git
```

**Verify installation:**
```bash
git --version
# Expected output: git version 2.x.x
```

#### 3. SQLite3
SQLite is required for the package database functionality.

**Install:**
```bash
# Linux (Debian/Ubuntu)
sudo apt-get install sqlite3 libsqlite3-dev

# Linux (Fedora/RHEL)
sudo dnf install sqlite sqlite-devel

# macOS
brew install sqlite3
```

**Verify installation:**
```bash
sqlite3 --version
# Expected output: 3.x.x
```

#### 4. Build Tools

**Linux:**
```bash
# Debian/Ubuntu
sudo apt-get install build-essential

# Fedora/RHEL
sudo dnf groupinstall "Development Tools"
```

**macOS:**
```bash
xcode-select --install
```

### Optional Tools

#### Make (Recommended)
For using build scripts and shortcuts.

```bash
# Linux
sudo apt-get install make  # Debian/Ubuntu
sudo dnf install make      # Fedora/RHEL

# macOS
brew install make
```

#### golangci-lint (Recommended for Code Quality)
```bash
# Install via go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Or via script (Linux/macOS)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

## Cloning the Repository

### 1. Fork the Repository (for Contributors)

If you plan to contribute, first fork the repository on GitHub:
1. Navigate to [https://github.com/grave0x/gpkg](https://github.com/grave0x/gpkg)
2. Click the "Fork" button in the top right
3. Clone your fork (replace `YOUR_USERNAME` with your GitHub username):

```bash
git clone https://github.com/YOUR_USERNAME/gpkg.git
cd gpkg
```

### 2. Clone Directly (for Testing/Exploring)

```bash
git clone https://github.com/grave0x/gpkg.git
cd gpkg
```

### 3. Set Up Git Remotes

If you forked the repository, add the upstream remote:

```bash
# Add the original repository as upstream
git remote add upstream https://github.com/grave0x/gpkg.git

# Verify remotes
git remote -v
# Expected output:
# origin    https://github.com/YOUR_USERNAME/gpkg.git (fetch)
# origin    https://github.com/YOUR_USERNAME/gpkg.git (push)
# upstream  https://github.com/grave0x/gpkg.git (fetch)
# upstream  https://github.com/grave0x/gpkg.git (push)
```

### 4. Branch Setup

Always work on a feature branch, not directly on `main`:

```bash
# Update your local main branch
git checkout main
git pull upstream main

# Create a new feature branch
git checkout -b feature/my-awesome-feature

# Or for bug fixes
git checkout -b fix/bug-description
```

## Dependencies

### 1. Download Go Dependencies

gpkg uses Go modules for dependency management. Download all required dependencies:

```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify
# Expected output: all modules verified
```

### 2. View Dependencies

```bash
# List all dependencies
go list -m all

# Check for available updates
go list -u -m all
```

### 3. Vendor Setup (Optional)

For air-gapped development or CI optimization:

```bash
# Create vendor directory
go mod vendor

# Verify vendor is up to date
go mod tidy
```

### 4. Update Dependencies

```bash
# Update all dependencies to latest minor/patch versions
go get -u ./...

# Update a specific dependency
go get -u github.com/spf13/cobra

# Clean up go.mod and go.sum
go mod tidy
```

## Building

### 1. Standard Development Build

Build the gpkg binary for local development:

```bash
# Build to bin/gpkg
go build -o bin/gpkg ./cmd/gpkg

# Or build and install to $GOPATH/bin
go build -o bin/gpkg .
```

**Verify the build:**
```bash
./bin/gpkg --version
./bin/gpkg --help
```

### 2. Development Build with Debugging

Build with debugging symbols and no optimizations:

```bash
go build -gcflags="all=-N -l" -o bin/gpkg-dev ./cmd/gpkg
```

### 3. Release Build with Version Info

Build with version information embedded:

```bash
VERSION=0.1.0
go build \
  -ldflags "-X main.version=${VERSION}" \
  -o bin/gpkg-${VERSION} \
  ./cmd/gpkg

# Verify version
./bin/gpkg-${VERSION} --version
```

### 4. Optimized Build

Build optimized binary (smaller size, stripped symbols):

```bash
go build \
  -ldflags "-s -w" \
  -trimpath \
  -o bin/gpkg \
  ./cmd/gpkg
```

### 5. Cross-Platform Builds

Build for different operating systems and architectures:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o bin/gpkg-linux-amd64 ./cmd/gpkg

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o bin/gpkg-linux-arm64 ./cmd/gpkg

# macOS AMD64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o bin/gpkg-darwin-amd64 ./cmd/gpkg

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/gpkg-darwin-arm64 ./cmd/gpkg

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o bin/gpkg.exe ./cmd/gpkg
```

### 6. Install System-Wide (Optional)

```bash
# Install to $GOPATH/bin (usually ~/go/bin)
go install ./cmd/gpkg

# Or copy to system path
sudo cp bin/gpkg /usr/local/bin/

# Verify installation
which gpkg
gpkg --version
```

## Running Tests

### 1. Run All Tests

```bash
# Run all tests in the project
go test ./...

# Run with verbose output
go test -v ./...
```

### 2. Run Specific Package Tests

```bash
# Test specific package
go test -v ./internal/manifest
go test -v ./internal/config
go test -v ./internal/pkgdb
go test -v ./internal/source

# Test command implementations
go test -v ./cmd/gpkg/cmd
```

### 3. Run Tests with Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage percentage
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open in browser (Linux)
xdg-open coverage.html

# Open in browser (macOS)
open coverage.html
```

### 4. Run Tests with Race Detector

Detect race conditions in concurrent code:

```bash
go test -race ./...
```

### 5. Run Integration Tests

```bash
# Run integration tests (if they exist in tests/ directory)
go test -v ./tests/integration

# Run with longer timeout for slow tests
go test -v -timeout 5m ./tests/integration
```

### 6. Run Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run benchmarks with memory stats
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkManifestParse ./internal/manifest
```

### 7. Watch Tests (Continuous Testing)

Using a file watcher for automatic test runs:

```bash
# Install gowatch (if not already installed)
go install github.com/silenceper/gowatch@latest

# Run tests automatically on file changes
gowatch -o ./bin/gpkg -p ./ -v ./...
```

## Development Workflow

### 1. Branch Strategy

gpkg follows a feature branch workflow:

- `main` - Stable, production-ready code
- `feature/*` - New features (e.g., `feature/add-rollback`)
- `fix/*` - Bug fixes (e.g., `fix/checksum-validation`)
- `docs/*` - Documentation updates (e.g., `docs/api-guide`)
- `test/*` - Test improvements (e.g., `test/integration-coverage`)

**Creating a feature branch:**
```bash
# Ensure main is up to date
git checkout main
git pull upstream main

# Create and switch to feature branch
git checkout -b feature/my-feature

# Work on your changes...
```

### 2. Making Commits

Follow these commit message conventions:

```bash
# Feature commits
git commit -m "feat: add rollback command for package downgrades"

# Bug fix commits
git commit -m "fix: correct checksum validation for SHA512"

# Documentation commits
git commit -m "docs: update development setup guide"

# Test commits
git commit -m "test: add integration tests for install command"

# Refactor commits
git commit -m "refactor: extract download logic to separate package"
```

**Commit message format:**
```
<type>: <short description>

<optional detailed description>

<optional footer>
```

**Types:** `feat`, `fix`, `docs`, `test`, `refactor`, `style`, `chore`, `perf`

### 3. Code Quality Checks

Before committing, run these checks:

```bash
# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run ./...

# Check for issues
go vet ./...

# Run all tests
go test ./...
```

**Pre-commit checklist:**
- [ ] Code is formatted (`go fmt`)
- [ ] No linter warnings (`golangci-lint run`)
- [ ] No vet issues (`go vet`)
- [ ] All tests pass (`go test ./...`)
- [ ] New code has tests
- [ ] Documentation updated if needed

### 4. Creating Pull Requests

```bash
# Push your branch to your fork
git push origin feature/my-feature

# Go to GitHub and create a Pull Request
# 1. Navigate to your fork: https://github.com/YOUR_USERNAME/gpkg
# 2. Click "Pull Request" button
# 3. Select base: grave0x/gpkg main <- compare: YOUR_USERNAME/gpkg feature/my-feature
# 4. Fill in PR description
# 5. Submit PR
```

**PR Description Template:**
```markdown
## Description
Brief description of what this PR does.

## Changes
- Change 1
- Change 2
- Change 3

## Testing
How to test these changes:
1. Step 1
2. Step 2

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] All tests pass
- [ ] Code formatted and linted
```

### 5. Keeping Your Fork Updated

```bash
# Fetch latest changes from upstream
git fetch upstream

# Merge upstream main into your main
git checkout main
git merge upstream/main

# Push updates to your fork
git push origin main

# Rebase your feature branch (if needed)
git checkout feature/my-feature
git rebase main
```

## IDE Setup

### Visual Studio Code

#### Recommended Extensions

Install these extensions for the best experience:

1. **Go** (golang.go) - Official Go extension
2. **Go Test Explorer** (premparihar.gotestexplorer)
3. **Go Doc** (msyrus.go-doc)
4. **Error Lens** (usernamehw.errorlens)
5. **GitLens** (eamodio.gitlens)

**Install via command line:**
```bash
code --install-extension golang.go
code --install-extension premparihar.gotestexplorer
code --install-extension msyrus.go-doc
code --install-extension usernamehw.errorlens
code --install-extension eamodio.gitlens
```

#### VS Code Settings

Create or update `.vscode/settings.json`:

```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "gofmt",
  "go.useLanguageServer": true,
  "go.buildOnSave": "workspace",
  "go.testOnSave": false,
  "go.coverOnSave": false,
  "go.testFlags": ["-v", "-race"],
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "files.exclude": {
    "**/.git": true,
    "**/.DS_Store": true,
    "**/bin": true,
    "**/coverage.out": true,
    "**/coverage.html": true
  }
}
```

#### Launch Configuration for Debugging

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch gpkg",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/gpkg",
      "args": ["--help"]
    },
    {
      "name": "Launch gpkg install",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/gpkg",
      "args": ["install", "./examples/manifest.yaml"]
    },
    {
      "name": "Test Current File",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${file}"
    },
    {
      "name": "Test Current Package",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${fileDirname}"
    }
  ]
}
```

#### Tasks Configuration

Create `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Build",
      "type": "shell",
      "command": "go build -o bin/gpkg ./cmd/gpkg",
      "group": {
        "kind": "build",
        "isDefault": true
      },
      "problemMatcher": ["$go"]
    },
    {
      "label": "Test",
      "type": "shell",
      "command": "go test -v ./...",
      "group": {
        "kind": "test",
        "isDefault": true
      },
      "problemMatcher": ["$go"]
    },
    {
      "label": "Lint",
      "type": "shell",
      "command": "golangci-lint run ./...",
      "problemMatcher": ["$go"]
    }
  ]
}
```

### GoLand / IntelliJ IDEA

#### Project Setup

1. **Open Project:**
   - File → Open → Select `gpkg` directory
   - GoLand will automatically detect it as a Go module

2. **Configure Go SDK:**
   - File → Settings → Go → GOROOT
   - Select Go 1.20+ installation

3. **Enable Go Modules:**
   - File → Settings → Go → Go Modules
   - Check "Enable Go modules integration"
   - Set "Environment" to `GOPROXY=https://proxy.golang.org,direct`

#### Run Configurations

**Build Configuration:**
1. Run → Edit Configurations → + → Go Build
2. **Name:** Build gpkg
3. **Package path:** `./cmd/gpkg`
4. **Output directory:** `bin`
5. **Working directory:** Project root

**Test Configuration:**
1. Run → Edit Configurations → + → Go Test
2. **Name:** All Tests
3. **Test kind:** Directory
4. **Directory:** Project root
5. **Pattern:** (empty for all)

**Debug Configuration:**
1. Run → Edit Configurations → + → Go Build
2. **Name:** Debug gpkg
3. **Package path:** `./cmd/gpkg`
4. **Program arguments:** `--help` (or other commands)
5. **Working directory:** Project root

#### Recommended Plugins

- **Markdown** - For editing documentation
- **YAML/Ansible** - For manifest editing
- **.ignore** - For .gitignore support
- **GitToolBox** - Enhanced Git integration

## Debugging

### Command-Line Debugging with Delve

[Delve](https://github.com/go-delve/delve) is the Go debugger.

#### Install Delve

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

#### Debug the Application

```bash
# Debug the main application
dlv debug ./cmd/gpkg -- --help

# Debug with specific arguments
dlv debug ./cmd/gpkg -- install ./examples/manifest.yaml

# Debug a specific test
dlv test ./internal/manifest -- -test.run TestManifestParse
```

#### Common Delve Commands

```
(dlv) break main.main          # Set breakpoint at main
(dlv) break manifest.go:42     # Set breakpoint at file:line
(dlv) breakpoints              # List all breakpoints
(dlv) continue                 # Continue execution
(dlv) next                     # Step to next line
(dlv) step                     # Step into function
(dlv) print myVar              # Print variable
(dlv) locals                   # Print local variables
(dlv) goroutines               # List goroutines
(dlv) help                     # Show help
```

### Debugging in VS Code

1. **Set Breakpoints:** Click in the gutter next to line numbers
2. **Start Debugging:** Press `F5` or Run → Start Debugging
3. **Use Debug Console:** View variables, call stack, and execute commands
4. **Step Through Code:**
   - `F10` - Step Over
   - `F11` - Step Into
   - `Shift+F11` - Step Out
   - `F5` - Continue

### Debugging in GoLand

1. **Set Breakpoints:** Click in the gutter next to line numbers
2. **Debug:** Right-click on file → Debug 'go build main.go'
3. **Debug Panel:** View variables, watches, call stack
4. **Step Controls:** Use toolbar or keyboard shortcuts

### Debugging Tips

```bash
# Print debug output
log.Printf("Debug: value = %v", value)

# Add temporary debugging code
fmt.Fprintf(os.Stderr, "DEBUG: %+v\n", myStruct)

# Enable verbose logging
export LOG_LEVEL=debug
./bin/gpkg install ./examples/manifest.yaml

# Use -v flag
./bin/gpkg -v install package-name
```

### Troubleshooting Common Issues

**Issue: "dlv: command not found"**
```bash
# Ensure $GOPATH/bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

**Issue: Breakpoints not hitting**
- Ensure you're building with `-gcflags="all=-N -l"` (no optimizations)
- Check that source code matches the binary

**Issue: Cannot inspect variables**
- Build with debug symbols (no `-ldflags="-s -w"`)
- Variables might be optimized out; rebuild with `-gcflags="all=-N -l"`

## Local Testing

### 1. Test Installation Locally

Create a test prefix directory to avoid affecting your system:

```bash
# Create test environment
mkdir -p ~/gpkg-test
export GPKG_INSTALL_PREFIX=~/gpkg-test

# Build and test
go build -o bin/gpkg ./cmd/gpkg
./bin/gpkg --prefix ~/gpkg-test install ./examples/manifest.yaml

# Verify installation
ls -la ~/gpkg-test/
```

### 2. Test with Local Manifests

Create test manifests in the `examples/` directory:

```bash
# Test manifest validation
./bin/gpkg validate ./examples/manifest.yaml

# Test installation from local manifest
./bin/gpkg install ./examples/binary-release.yaml

# Test source build (if applicable)
./bin/gpkg install --from-source ./examples/source-build.yaml
```

### 3. Test Source Management

```bash
# Test adding a source
./bin/gpkg add-source https://raw.githubusercontent.com/grave0x/gpkg-registry/main/index.json

# List sources
./bin/gpkg list-sources

# Update indices
./bin/gpkg update

# Search packages
./bin/gpkg search tool

# View package info
./bin/gpkg info package-name
```

### 4. Test Package Database

```bash
# View installed packages
./bin/gpkg list --installed

# Check package database
sqlite3 ~/.gpkg/pkgdb.sqlite "SELECT * FROM packages;"

# View package files
sqlite3 ~/.gpkg/pkgdb.sqlite "SELECT * FROM package_files WHERE package_id='package-name';"
```

### 5. Test Error Handling

```bash
# Test invalid manifest
./bin/gpkg validate ./examples/invalid-manifest.yaml

# Test missing package
./bin/gpkg install non-existent-package

# Test corrupted download (modify checksum in manifest)
./bin/gpkg install ./examples/bad-checksum.yaml
```

### 6. Test CLI Output Formats

```bash
# Test JSON output
./bin/gpkg --json list --installed

# Test verbose output
./bin/gpkg -v install package-name

# Test quiet mode
./bin/gpkg --quiet install package-name

# Test dry-run
./bin/gpkg --dry-run install package-name
```

### 7. Integration Testing Script

Create a test script `scripts/test-local.sh`:

```bash
#!/bin/bash
set -e

echo "=== gpkg Local Testing ==="

# Setup test environment
TEST_PREFIX="/tmp/gpkg-test-$$"
mkdir -p "$TEST_PREFIX"
export GPKG_INSTALL_PREFIX="$TEST_PREFIX"

# Build
echo "Building gpkg..."
go build -o bin/gpkg ./cmd/gpkg

# Test commands
echo "Testing version..."
./bin/gpkg --version

echo "Testing help..."
./bin/gpkg --help

echo "Testing add-source..."
./bin/gpkg add-source https://example.com/index.json

echo "Testing list-sources..."
./bin/gpkg list-sources

# Cleanup
echo "Cleaning up..."
rm -rf "$TEST_PREFIX"

echo "=== All tests passed ==="
```

### 8. Clean Up Test Environment

```bash
# Remove test prefix
rm -rf ~/gpkg-test

# Clean build artifacts
rm -f bin/gpkg bin/gpkg-dev
rm -f coverage.out coverage.html

# Clean Go cache (if needed)
go clean -cache -testcache -modcache
```

## Next Steps

Now that you have your development environment set up:

1. **Read the Documentation:**
   - [Architecture Overview](./Architecture.md)
   - [Contributing Guidelines](../../CONTRIBUTING.md)
   - [Code Standards](./Code-Standards.md)

2. **Explore the Codebase:**
   - Review the project structure in `DEVELOPMENT.md`
   - Read through existing code in `internal/` packages
   - Check out `examples/` for sample manifests

3. **Pick Your First Issue:**
   - Browse [GitHub Issues](https://github.com/grave0x/gpkg/issues)
   - Look for `good first issue` labels
   - Comment on an issue to claim it

4. **Join the Community:**
   - Follow the repository for updates
   - Join discussions in issues and PRs
   - Share your ideas and improvements

## Getting Help

If you encounter issues during setup:

- **Documentation:** Check [BUILD.md](../../BUILD.md) and [DEVELOPMENT.md](../../DEVELOPMENT.md)
- **Issues:** Search existing [GitHub Issues](https://github.com/grave0x/gpkg/issues)
- **Questions:** Open a new issue with the `question` label
- **Wiki:** Browse the [project wiki](https://github.com/grave0x/gpkg/wiki) for guides

Happy coding! 🚀
