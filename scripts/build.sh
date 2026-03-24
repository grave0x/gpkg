#!/bin/bash
# Build script for gpkg - cross-platform builds

set -e

PROJECT_ROOT=$(dirname "$0")
BINARY_NAME="gpkg"
VERSION="${1:-0.1.0}"
OUT_DIR="$PROJECT_ROOT/dist"

echo "Building gpkg v${VERSION}..."

mkdir -p "$OUT_DIR"

# Helper function to build for a platform
build_platform() {
    local os=$1
    local arch=$2
    local output="$3"
    
    echo "  Building for ${os}/${arch}..."
    GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build \
        -ldflags "-X main.version=${VERSION} -s -w" \
        -trimpath \
        -o "$output" \
        ./cmd/gpkg
    
    # Make executable on Unix platforms
    if [[ "$os" != "windows" ]]; then
        chmod +x "$output"
    fi
}

# Linux builds
build_platform "linux" "amd64" "$OUT_DIR/gpkg-${VERSION}-linux-amd64"
build_platform "linux" "arm64" "$OUT_DIR/gpkg-${VERSION}-linux-arm64"
build_platform "linux" "386" "$OUT_DIR/gpkg-${VERSION}-linux-386"

# macOS builds
build_platform "darwin" "amd64" "$OUT_DIR/gpkg-${VERSION}-darwin-amd64"
build_platform "darwin" "arm64" "$OUT_DIR/gpkg-${VERSION}-darwin-arm64"

# Windows builds
build_platform "windows" "amd64" "$OUT_DIR/gpkg-${VERSION}-windows-amd64.exe"
build_platform "windows" "386" "$OUT_DIR/gpkg-${VERSION}-windows-386.exe"

# Generate checksums
echo "Generating checksums..."
cd "$OUT_DIR"
if command -v sha256sum &> /dev/null; then
    sha256sum gpkg-${VERSION}-* > SHA256SUMS
elif command -v shasum &> /dev/null; then
    shasum -a 256 gpkg-${VERSION}-* > SHA256SUMS
fi

echo "Build complete! Binaries in $OUT_DIR"
echo ""
ls -lh "$OUT_DIR"/gpkg-${VERSION}-*
