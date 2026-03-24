# Build & Installation Guide

## Quick Start

### Prerequisites

- **Go 1.20+** - [Download Go](https://golang.org/dl)
- **Git** - Required for cloning and source builds
- **Make** (optional) - For build commands
- **bash/sh** - For running build scripts

### Build from Source

```bash
# Clone the repository
git clone https://github.com/grave0x/gpkg.git
cd gpkg

# Fetch dependencies
go mod download

# Build the binary
go build -o bin/gpkg ./cmd/gpkg

# Verify the build
./bin/gpkg --version
```

### Install to System

```bash
# Install to $GOPATH/bin (usually ~/go/bin)
go install ./cmd/gpkg

# Copy to system path (requires sudo)
sudo cp bin/gpkg /usr/local/bin/

# Verify installation
gpkg --version
gpkg --help
```

## Build Variants

### Development Build
```bash
go build -o bin/gpkg-dev ./cmd/gpkg
```

### Release Build with Version Info
```bash
VERSION=0.1.0
go build \
  -ldflags "-X main.version=${VERSION}" \
  -o bin/gpkg-${VERSION} \
  ./cmd/gpkg
```

### Optimized Build
```bash
go build \
  -ldflags "-s -w" \
  -trimpath \
  -o bin/gpkg \
  ./cmd/gpkg
```

### Build for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/gpkg-linux-amd64 ./cmd/gpkg

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/gpkg-darwin-amd64 ./cmd/gpkg
GOOS=darwin GOARCH=arm64 go build -o bin/gpkg-darwin-arm64 ./cmd/gpkg

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/gpkg.exe ./cmd/gpkg
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test -v ./internal/manifest

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Development Workflow

### Add Dependencies

```bash
# Add a new dependency
go get github.com/some/package

# Update go.mod and go.sum
go mod tidy

# Check for vulnerabilities
go list -u -m all
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run ./...

# Check for issues
go vet ./...
```

### Building a Release

1. Update version in relevant files
2. Run tests: `go test ./...`
3. Build release: `go build -ldflags "-X main.version=X.Y.Z" ...`
4. Tag release: `git tag vX.Y.Z && git push --tags`

## Docker Build (Optional)

If you want to provide Docker builds:

```dockerfile
# Build stage
FROM golang:1.20 as builder
WORKDIR /app
COPY . .
RUN go build -o bin/gpkg ./cmd/gpkg

# Runtime stage
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/bin/gpkg /usr/local/bin/
ENTRYPOINT ["gpkg"]
```

Build with:
```bash
docker build -t gpkg:latest .
docker run gpkg --help
```

## Installation Verification

After installation, verify everything works:

```bash
# Show version
gpkg --version

# Show help
gpkg --help

# Test commands
gpkg list-sources          # Should work (empty if no sources)
gpkg add-source https://example.com/packages  # Should succeed
gpkg list-sources          # Should show the added source
```

## Troubleshooting

### "go: command not found"
- Go is not installed or not in PATH
- Install from https://golang.org/dl
- Add `$GOPATH/bin` (usually ~/go/bin) to PATH

### "git: command not found"
- Git is not installed or not in PATH
- Install git from https://git-scm.com

### Go Module Issues
```bash
# Clear cache
go clean -modcache

# Verify dependencies
go mod verify

# Update to latest compatible versions
go mod tidy
```

### Build Fails on Windows
- Use MinGW or Git Bash
- Ensure paths don't have spaces
- Try building with `-v` flag for verbose output

## Performance Tips

1. **Faster Rebuilds**: Use `-trimpath` flag
2. **Smaller Binaries**: Use `-ldflags "-s -w"` to strip symbols
3. **Parallel Builds**: Go automatically parallelizes; set `GOMAXPROCS`

## Contributing

To set up development environment:

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gpkg.git
cd gpkg

# Add upstream remote
git remote add upstream https://github.com/grave0x/gpkg.git

# Create feature branch
git checkout -b feature/my-feature

# Make changes, test
go test ./...

# Push and create PR
git push origin feature/my-feature
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.
