# Testing Guide

Comprehensive guide for writing, running, and maintaining tests in the gpkg project.

## Table of Contents
- [Test Organization](#test-organization)
- [Unit Tests](#unit-tests)
- [Integration Tests](#integration-tests)
- [Test Coverage](#test-coverage)
- [Benchmarking](#benchmarking)
- [Mocking](#mocking)
- [CI/CD Testing](#cicd-testing)
- [Best Practices](#best-practices)

---

## Test Organization

### Directory Structure

Tests are organized alongside the code they test:

```
gpkg/
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   ├── config_test.go        # Unit tests for config package
│   │   └── loader.go
│   ├── download/
│   │   ├── download.go
│   │   ├── parallel.go
│   │   └── parallel_test.go      # Unit tests for download package
│   ├── manifest/
│   │   ├── manifest.go
│   │   └── manifest_test.go      # Unit tests for manifest package
│   ├── pkgdb/
│   │   ├── sqlite.go
│   │   └── sqlite_test.go        # Unit tests for database package
│   └── ...
└── tests/
    └── integration/
        └── workflows_test.go      # Integration/workflow tests
```

### Naming Conventions

**Test Files:**
- Use `<filename>_test.go` pattern
- Place in same package directory as code under test
- Use `package <name>_test` for external (black-box) testing

**Test Functions:**
- `Test<FunctionName>` - for standard unit tests
- `Test<Feature>` - for feature tests
- `Benchmark<Operation>` - for benchmark tests (when needed)
- `Example<Function>` - for example code (documentation)

**Examples:**
```go
func TestDefaultConfig(t *testing.T)           // Tests DefaultConfig function
func TestYAMLConfigLoading(t *testing.T)       // Tests YAML loading feature
func TestSQLiteManager(t *testing.T)           // Tests SQLiteManager type
func BenchmarkParallelDownload(b *testing.B)   // Benchmark test
```

---

## Unit Tests

Unit tests focus on testing individual functions and components in isolation.

### Basic Test Structure

```go
package config_test

import (
    "testing"
    
    "github.com/grave0x/gpkg/internal/config"
)

func TestDefaultConfig(t *testing.T) {
    cfg := config.DefaultConfig()

    if cfg == nil {
        t.Fatal("expected config but got nil")
    }

    if cfg.LogLevel != "info" {
        t.Errorf("expected log level 'info', got %s", cfg.LogLevel)
    }

    if cfg.Color != true {
        t.Errorf("expected color enabled, got %v", cfg.Color)
    }
}
```

### Table-Driven Tests

Use table-driven tests for testing multiple scenarios:

```go
func TestParseGitHubIdentifier(t *testing.T) {
    r := resolver.NewGitHubResolver()

    tests := []struct {
        input   string
        owner   string
        repo    string
        pkgType string
        fail    bool
    }{
        {"owner/repo", "owner", "repo", "github", false},
        {"github:owner/repo", "owner", "repo", "github", false},
        {"https://github.com/owner/repo", "owner", "repo", "github", false},
        {"invalid", "", "", "", true},
        {"owner/repo/extra", "", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            _, owner, repo, err := r.ParseIdentifier(tt.input)

            if tt.fail {
                if err == nil {
                    t.Errorf("expected error for %s", tt.input)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if owner != tt.owner || repo != tt.repo {
                    t.Errorf("expected %s/%s, got %s/%s", 
                        tt.owner, tt.repo, owner, repo)
                }
            }
        })
    }
}
```

### Testing with Temporary Files

Use `t.TempDir()` for file-based tests:

```go
func TestYAMLConfigLoading(t *testing.T) {
    tmpDir := t.TempDir()
    configFile := filepath.Join(tmpDir, "config.yaml")

    // Create test config
    configContent := `
prefix: /custom/prefix
cache_dir: /custom/cache
log_level: debug
color: false
`

    if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
        t.Fatalf("failed to write test config: %v", err)
    }

    loader := config.NewYAMLLoader("")
    cfg, err := loader.LoadFrom(configFile)

    if err != nil {
        t.Fatalf("failed to load config: %v", err)
    }

    if cfg.Prefix != "/custom/prefix" {
        t.Errorf("expected prefix '/custom/prefix', got %s", cfg.Prefix)
    }
}
```

### Testing Database Operations

Example from `sqlite_test.go`:

```go
func TestSQLiteManager(t *testing.T) {
    tmpDir := t.TempDir()
    dbPath := filepath.Join(tmpDir, "test.db")

    manager, err := pkgdb.NewSQLiteManager(dbPath)
    if err != nil {
        t.Fatalf("failed to create manager: %v", err)
    }
    defer manager.Close()

    // Test AddPackage
    pkg := &pkgdb.PackageRecord{
        Name:         "test-pkg",
        Version:      "1.0.0",
        Source:       "release",
        Prefix:       "/home/user/.gpkg",
        Author:       "Test Author",
        URL:          "https://example.com",
        License:      "MIT",
        Checksums:    map[string]string{"sha256": "abc123"},
        Dependencies: []string{"dep1"},
    }

    id, err := manager.AddPackage(pkg)
    if err != nil {
        t.Fatalf("failed to add package: %v", err)
    }

    if id == 0 {
        t.Errorf("expected non-zero package ID")
    }

    // Test GetPackage
    retrieved, err := manager.GetPackage("test-pkg")
    if err != nil {
        t.Fatalf("failed to get package: %v", err)
    }

    if retrieved.Name != "test-pkg" {
        t.Errorf("expected name 'test-pkg', got %s", retrieved.Name)
    }
}
```

### Running Unit Tests

```bash
# Run all tests
go test ./...

# Run tests in a specific package
go test ./internal/config

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestDefaultConfig ./internal/config

# Run tests matching a pattern
go test -v -run TestSQLite ./...

# Run with race detector
go test -race ./...
```

---

## Integration Tests

Integration tests verify complete workflows and component interactions.

### Integration Test Structure

Located in `tests/integration/`:

```go
package integration_test

import (
    "os"
    "testing"
)

// TestInstallWorkflow tests the complete install workflow
func TestInstallWorkflow(t *testing.T) {
    // This would:
    // 1. Create temp directory
    // 2. Add a test source
    // 3. Search for a package
    // 4. Install the package
    // 5. Verify installation in pkgdb
    // 6. Verify files exist
    t.Skip("Integration test - requires running gpkg binary")
}

// TestSourceManagement tests source add/remove/list workflow
func TestSourceManagement(t *testing.T) {
    tmpDir := t.TempDir()
    oldHome := os.Getenv("HOME")
    os.Setenv("HOME", tmpDir)
    defer os.Setenv("HOME", oldHome)

    // Test:
    // 1. No sources initially
    // 2. Add source
    // 3. List sources shows it
    // 4. Remove source
    // 5. List sources is empty
}
```

### Workflow Tests

Integration tests verify complete user workflows:

- **Install Workflow**: Add source → Search → Install → Verify
- **Upgrade Workflow**: Install v1.0.0 → Upgrade to v2.0.0 → Verify
- **Uninstall Workflow**: Install → Uninstall → Verify removal
- **Rollback Workflow**: Install v1 → Upgrade v2 → Rollback to v1
- **Dependency Install**: Install package with dependencies → Verify all installed

### Running Integration Tests

```bash
# Run all integration tests
go test ./tests/integration/...

# Run specific integration test
go test -v -run TestInstallWorkflow ./tests/integration

# Run with timeout (for long-running tests)
go test -timeout 10m ./tests/integration/...
```

### Environment Setup for Integration Tests

```go
func setupTestEnvironment(t *testing.T) (cleanup func()) {
    tmpDir := t.TempDir()
    
    // Override environment variables
    oldHome := os.Getenv("HOME")
    oldPrefix := os.Getenv("GPKG_PREFIX")
    
    os.Setenv("HOME", tmpDir)
    os.Setenv("GPKG_PREFIX", filepath.Join(tmpDir, ".gpkg"))
    
    return func() {
        os.Setenv("HOME", oldHome)
        os.Setenv("GPKG_PREFIX", oldPrefix)
    }
}

func TestWithCleanEnvironment(t *testing.T) {
    cleanup := setupTestEnvironment(t)
    defer cleanup()
    
    // Your test code here
}
```

---

## Test Coverage

### Generating Coverage Reports

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out

# Generate coverage with mode atomic (for race detector)
go test -coverprofile=coverage.out -covermode=atomic ./...

# Coverage for specific package
go test -coverprofile=coverage.out ./internal/config
go tool cover -html=coverage.out
```

### Coverage Requirements

**Project Coverage Goals:**
- **Overall**: ≥ 70% coverage
- **Core packages** (config, pkgdb, planner): ≥ 80% coverage
- **Utility packages** (download, manifest): ≥ 75% coverage
- **Integration tests**: Workflow coverage

### Checking Coverage by Package

```bash
# View coverage per package
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep -E '^github.com/grave0x/gpkg'

# Example output:
# github.com/grave0x/gpkg/internal/config/config.go:12:        DefaultConfig      100.0%
# github.com/grave0x/gpkg/internal/config/loader.go:15:        LoadFrom           85.7%
# github.com/grave0x/gpkg/internal/pkgdb/sqlite.go:25:         NewSQLiteManager   90.5%
```

### Coverage in CI/CD

Coverage is automatically uploaded to Codecov on:
- Ubuntu-latest
- Go 1.22
- Main branch and PRs

```yaml
- name: Run tests with coverage
  run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v4
  with:
    file: ./coverage.out
    flags: unittests
```

---

## Benchmarking

### Writing Benchmark Tests

Benchmark tests measure performance of critical operations.

```go
func BenchmarkParallelDownload(b *testing.B) {
    mock := &MockDownloader{}
    pd := download.NewParallelDownloader(mock, 4)
    
    items := []*download.DownloadItem{
        {URL: "https://example.com/pkg1.tar.gz", Dest: "/tmp/pkg1.tar.gz"},
        {URL: "https://example.com/pkg2.tar.gz", Dest: "/tmp/pkg2.tar.gz"},
        {URL: "https://example.com/pkg3.tar.gz", Dest: "/tmp/pkg3.tar.gz"},
    }
    
    ctx := context.Background()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pd.DownloadMultiple(ctx, items)
    }
}

func BenchmarkSQLiteInsert(b *testing.B) {
    tmpDir := b.TempDir()
    dbPath := filepath.Join(tmpDir, "bench.db")
    
    manager, _ := pkgdb.NewSQLiteManager(dbPath)
    defer manager.Close()
    
    pkg := &pkgdb.PackageRecord{
        Name:    "bench-pkg",
        Version: "1.0.0",
        Source:  "release",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        manager.AddPackage(pkg)
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkParallelDownload ./internal/download

# Run with memory allocation stats
go test -bench=. -benchmem ./...

# Run benchmarks multiple times for accuracy
go test -bench=. -benchtime=10s -count=5 ./...

# Compare benchmarks (requires benchstat)
go test -bench=. -count=10 | tee old.txt
# Make changes...
go test -bench=. -count=10 | tee new.txt
benchstat old.txt new.txt
```

### Benchmark Output

```
BenchmarkParallelDownload-8      1000000      1234 ns/op      512 B/op      10 allocs/op
BenchmarkSQLiteInsert-8           500000      2345 ns/op     1024 B/op      15 allocs/op
```

- **1000000**: Number of iterations
- **1234 ns/op**: Nanoseconds per operation
- **512 B/op**: Bytes allocated per operation
- **10 allocs/op**: Number of allocations per operation

---

## Mocking

### Mock Implementations

Mocks allow testing without external dependencies.

#### Example: Mock Downloader

```go
// MockDownloader is a mock implementation for testing
type MockDownloader struct {
    downloaded []string
}

func (m *MockDownloader) Download(ctx context.Context, url, dest string) error {
    m.downloaded = append(m.downloaded, url)
    return nil
}

func (m *MockDownloader) DownloadWithChecksum(ctx context.Context, url, dest, hash, algo string) error {
    m.downloaded = append(m.downloaded, url)
    return nil
}

func (m *MockDownloader) ValidateChecksum(filePath, hash, algo string) (bool, error) {
    return true, nil
}
```

#### Using Mocks in Tests

```go
func TestParallelDownloaderSuccess(t *testing.T) {
    mock := &MockDownloader{}
    pd := download.NewParallelDownloader(mock, 2)

    items := []*download.DownloadItem{
        {URL: "https://example.com/pkg1.tar.gz", Dest: "/tmp/pkg1.tar.gz"},
        {URL: "https://example.com/pkg2.tar.gz", Dest: "/tmp/pkg2.tar.gz"},
    }

    ctx := context.Background()
    err := pd.DownloadMultiple(ctx, items)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    if len(mock.downloaded) != 2 {
        t.Errorf("expected 2 downloads, got %d", len(mock.downloaded))
    }
}
```

### Interface-Based Design

The project uses interfaces to enable mocking:

```go
// Internal packages define interfaces
type Downloader interface {
    Download(ctx context.Context, url, dest string) error
    DownloadWithChecksum(ctx context.Context, url, dest, hash, algo string) error
    ValidateChecksum(filePath, hash, algo string) (bool, error)
}

type Registry interface {
    AddSource(ctx context.Context, src *Source) error
    RemoveSource(ctx context.Context, id string) error
    GetSource(ctx context.Context, id string) (*Source, error)
    ListSources(ctx context.Context) ([]*Source, error)
}

type Parser interface {
    Parse(path string) (*Manifest, error)
    ParseBytes(data []byte) (*Manifest, error)
}
```

### Testing Error Conditions

```go
// Mock with error injection
type ErrorDownloader struct {
    err error
}

func (e *ErrorDownloader) Download(ctx context.Context, url, dest string) error {
    return e.err
}

func TestDownloadError(t *testing.T) {
    mock := &ErrorDownloader{err: fmt.Errorf("network error")}
    pd := download.NewParallelDownloader(mock, 2)
    
    items := []*download.DownloadItem{
        {URL: "https://example.com/pkg.tar.gz", Dest: "/tmp/pkg.tar.gz"},
    }
    
    err := pd.DownloadMultiple(context.Background(), items)
    if err == nil {
        t.Error("expected error but got nil")
    }
}
```

---

## CI/CD Testing

### GitHub Actions Workflow

The project uses GitHub Actions for continuous testing across platforms.

**Workflow file:** `.github/workflows/test.yml`

#### Test Matrix

Tests run on:
- **Operating Systems**: Ubuntu, macOS, Windows
- **Go Versions**: 1.20, 1.21, 1.22
- **Total**: 9 test configurations (3 OS × 3 Go versions)

#### Workflow Steps

```yaml
jobs:
  test:
    name: Test on ${{ matrix.os }} with Go ${{ matrix.go-version }}
    runs-on: ${{ matrix.os }}
    strategy:
      fail-fast: false
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.20', '1.21', '1.22']
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ matrix.go-version }}
        cache: true

    - name: Download dependencies
      run: go mod download

    - name: Verify dependencies
      run: go mod verify

    - name: Run go vet
      run: go vet ./...

    - name: Run go fmt check
      run: |
        if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
          echo "The following files need formatting:"
          gofmt -s -l .
          exit 1
        fi

    - name: Run tests with coverage
      run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v4
      with:
        file: ./coverage.out
        flags: unittests
      if: matrix.os == 'ubuntu-latest' && matrix.go-version == '1.22'

    - name: Build
      run: go build -v ./...
```

#### Linting Job

Separate linting job using golangci-lint:

```yaml
lint:
  name: Lint
  runs-on: ubuntu-latest
  
  steps:
  - name: Checkout code
    uses: actions/checkout@v4

  - name: Set up Go
    uses: actions/setup-go@v5
    with:
      go-version: '1.22'

  - name: Run golangci-lint
    uses: golangci/golangci-lint-action@v4
    with:
      version: latest
      args: --timeout=5m
```

#### Build Verification

```yaml
build-all-platforms:
  name: Build for all platforms
  runs-on: ubuntu-latest
  needs: [test, lint]
  
  steps:
  - name: Build for multiple platforms
    run: |
      GOOS=linux GOARCH=amd64 go build -o dist/gpkg-linux-amd64 ./cmd/gpkg
      GOOS=linux GOARCH=arm64 go build -o dist/gpkg-linux-arm64 ./cmd/gpkg
      GOOS=darwin GOARCH=amd64 go build -o dist/gpkg-darwin-amd64 ./cmd/gpkg
      GOOS=darwin GOARCH=arm64 go build -o dist/gpkg-darwin-arm64 ./cmd/gpkg
      GOOS=windows GOARCH=amd64 go build -o dist/gpkg-windows-amd64.exe ./cmd/gpkg
```

### Local CI Simulation

Run the same checks locally before pushing:

```bash
# Run all CI checks locally
./scripts/ci-check.sh

# Or manually:
go mod download
go mod verify
go vet ./...
gofmt -s -l .
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go build -v ./...
```

### Pre-commit Hooks

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash

echo "Running pre-commit checks..."

# Format check
UNFORMATTED=$(gofmt -s -l .)
if [ -n "$UNFORMATTED" ]; then
    echo "❌ The following files need formatting:"
    echo "$UNFORMATTED"
    exit 1
fi

# Vet
go vet ./...
if [ $? -ne 0 ]; then
    echo "❌ go vet failed"
    exit 1
fi

# Tests
go test ./...
if [ $? -ne 0 ]; then
    echo "❌ Tests failed"
    exit 1
fi

echo "✅ All pre-commit checks passed"
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

## Best Practices

### 1. Test Naming and Organization

**✅ DO:**
- Use descriptive test names: `TestYAMLConfigLoading`
- Group related tests using subtests with `t.Run()`
- Use table-driven tests for multiple scenarios
- Keep test files in the same package directory

**❌ DON'T:**
- Use generic names like `TestFunction1`
- Mix unrelated tests in one function
- Create separate test directories far from code

### 2. Test Independence

**✅ DO:**
- Make tests independent and runnable in any order
- Use `t.TempDir()` for temporary files
- Clean up resources with `defer`
- Reset global state between tests

**❌ DON'T:**
- Depend on test execution order
- Use shared global state
- Leave temporary files after tests
- Share database connections between tests

```go
// Good: Independent test with cleanup
func TestDatabaseOperations(t *testing.T) {
    tmpDir := t.TempDir()  // Auto-cleaned
    db := setupTestDB(t, tmpDir)
    defer db.Close()
    
    // Test code
}

// Bad: Shared state
var sharedDB *sql.DB  // ❌ Tests will interfere

func TestInsert(t *testing.T) {
    sharedDB.Insert(...)  // ❌ Depends on execution order
}
```

### 3. Error Handling in Tests

**✅ DO:**
- Use `t.Fatal()` for setup failures
- Use `t.Error()` for test assertion failures
- Provide clear error messages with context

**❌ DON'T:**
- Ignore errors silently
- Use `panic()` in tests

```go
// Good: Clear error messages
func TestConfig(t *testing.T) {
    cfg, err := LoadConfig("config.yaml")
    if err != nil {
        t.Fatalf("failed to load config: %v", err)
    }
    
    if cfg.Port != 8080 {
        t.Errorf("expected port 8080, got %d", cfg.Port)
    }
}

// Bad: Unclear errors
func TestConfig(t *testing.T) {
    cfg, _ := LoadConfig("config.yaml")  // ❌ Ignored error
    
    if cfg.Port != 8080 {
        t.Error("wrong")  // ❌ Not descriptive
    }
}
```

### 4. Test Coverage Guidelines

**✅ DO:**
- Test happy paths and error conditions
- Test edge cases and boundary conditions
- Test concurrent access where applicable
- Focus on critical business logic

**❌ DON'T:**
- Chase 100% coverage for its own sake
- Test trivial getters/setters
- Test third-party library code
- Duplicate tests unnecessarily

```go
// Good: Test both success and failure
func TestParseVersion(t *testing.T) {
    tests := []struct {
        input   string
        want    string
        wantErr bool
    }{
        {"1.0.0", "1.0.0", false},           // Happy path
        {"v1.0.0", "1.0.0", false},          // With prefix
        {"", "", true},                       // Empty input
        {"invalid", "", true},                // Invalid format
        {"1.0.0-beta", "1.0.0-beta", false}, // With suffix
    }
    
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            got, err := ParseVersion(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseVersion(%q) error = %v, wantErr %v", 
                    tt.input, err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ParseVersion(%q) = %v, want %v", 
                    tt.input, got, tt.want)
            }
        })
    }
}
```

### 5. Using Subtests

**✅ DO:**
- Use `t.Run()` for related test cases
- Name subtests descriptively
- Use subtests for table-driven tests

```go
func TestManifestValidation(t *testing.T) {
    t.Run("valid manifest with install", func(t *testing.T) {
        // Test code
    })
    
    t.Run("missing package name", func(t *testing.T) {
        // Test code
    })
    
    t.Run("missing both install and build_source", func(t *testing.T) {
        // Test code
    })
}

// Run specific subtest:
// go test -v -run TestManifestValidation/valid_manifest
```

### 6. Testing Concurrency

**✅ DO:**
- Use `-race` flag to detect race conditions
- Test goroutine synchronization
- Use channels for coordination in tests

```go
func TestParallelDownload(t *testing.T) {
    mock := &MockDownloader{}
    pd := download.NewParallelDownloader(mock, 4)
    
    items := generateTestItems(100)
    
    ctx := context.Background()
    err := pd.DownloadMultiple(ctx, items)
    
    if err != nil {
        t.Errorf("parallel download failed: %v", err)
    }
    
    // Verify all downloaded
    if len(mock.downloaded) != 100 {
        t.Errorf("expected 100 downloads, got %d", len(mock.downloaded))
    }
}

// Run with race detector:
// go test -race ./internal/download
```

### 7. Test Data Management

**✅ DO:**
- Use realistic test data
- Create test fixtures for complex data
- Document test data requirements

```go
func newTestPackageRecord() *pkgdb.PackageRecord {
    return &pkgdb.PackageRecord{
        Name:         "test-pkg",
        Version:      "1.0.0",
        Source:       "release",
        Prefix:       "/home/user/.gpkg",
        Author:       "Test Author",
        URL:          "https://example.com",
        License:      "MIT",
        Checksums:    map[string]string{"sha256": "abc123"},
        Dependencies: []string{"dep1", "dep2"},
    }
}

func TestPackageOperations(t *testing.T) {
    pkg := newTestPackageRecord()
    // Use pkg in tests
}
```

### 8. Performance Testing

**✅ DO:**
- Benchmark critical paths
- Use `b.ResetTimer()` to exclude setup time
- Run benchmarks multiple times for accuracy
- Compare benchmarks before/after optimizations

```go
func BenchmarkDatabaseInsert(b *testing.B) {
    db := setupBenchDB(b)
    defer db.Close()
    
    pkg := newTestPackageRecord()
    
    b.ResetTimer()  // Don't count setup time
    for i := 0; i < b.N; i++ {
        db.AddPackage(pkg)
    }
}
```

### 9. Documentation

**✅ DO:**
- Document complex test scenarios
- Explain non-obvious test logic
- Add comments for future maintainers

```go
// TestVersionHistory verifies that package version history
// is correctly stored and retrieved in chronological order.
// This is critical for rollback functionality.
func TestVersionHistory(t *testing.T) {
    // Test implementation
}
```

### 10. Continuous Improvement

**✅ DO:**
- Review test failures in CI
- Add tests for bug fixes
- Refactor tests with code
- Monitor test execution time

**❌ DON'T:**
- Ignore flaky tests
- Skip failing tests without fixing
- Let test suite become slow

---

## Quick Reference

### Common Test Commands

```bash
# Run all tests
go test ./...

# Verbose output
go test -v ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# With race detector
go test -race ./...

# Specific package
go test ./internal/config

# Specific test
go test -v -run TestDefaultConfig ./internal/config

# Integration tests only
go test ./tests/integration/...

# Short tests only (skip long-running)
go test -short ./...

# Parallel execution
go test -parallel 4 ./...

# With timeout
go test -timeout 30s ./...

# Benchmarks
go test -bench=. ./...
go test -bench=. -benchmem ./...
```

### Test Assertion Patterns

```go
// Fatal vs Error
t.Fatal("stops test immediately")    // Use for setup failures
t.Error("continues test execution")  // Use for assertion failures

// Checking errors
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}

// Checking values
if got != want {
    t.Errorf("got %v, want %v", got, want)
}

// Checking nil
if result == nil {
    t.Fatal("expected result, got nil")
}

// Checking booleans
if !condition {
    t.Error("expected condition to be true")
}

// Checking slices/arrays
if len(items) != expected {
    t.Errorf("expected %d items, got %d", expected, len(items))
}
```

---

## Resources

### Go Testing Documentation
- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Blog: Table Driven Tests](https://go.dev/blog/table-driven-tests)
- [Go Blog: Coverage](https://go.dev/blog/cover)

### Tools
- **golangci-lint**: Comprehensive linter
- **go-test-coverage**: Enhanced coverage visualization
- **gotestsum**: Prettier test output
- **benchstat**: Statistical benchmark comparison

### Project-Specific
- [CI/CD Configuration](.github/workflows/test.yml)
- [Development Guide](DEVELOPMENT.md)
- [Build Guide](BUILD.md)
- [Contributing Guidelines](CONTRIBUTING.md)

---

## Getting Help

For questions about testing:
1. Check existing tests for examples
2. Review this guide
3. Check GitHub Actions logs for CI failures
4. Open an issue for test infrastructure problems
5. Discuss in pull requests for test design questions

**Remember:** Good tests are as important as good code. Invest time in writing clear, maintainable tests that verify correctness and catch regressions.
