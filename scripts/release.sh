#!/bin/bash
# Release script for gpkg

set -e

VERSION="${1:-0.1.0}"
PROJECT_ROOT=$(dirname "$0")

echo "Preparing release v${VERSION}..."

# Build binaries
bash "$PROJECT_ROOT/build.sh" "$VERSION"

# Create release notes
echo "Release gpkg v${VERSION}" > RELEASE_NOTES.md
echo "" >> RELEASE_NOTES.md
echo "## Changes" >> RELEASE_NOTES.md
echo "- Phase 7-10 implementation complete" >> RELEASE_NOTES.md
echo "- Full CLI with all subcommands" >> RELEASE_NOTES.md
echo "- Package database (SQLite)" >> RELEASE_NOTES.md
echo "- Planner component for dry-runs" >> RELEASE_NOTES.md
echo "- GitHub package resolver" >> RELEASE_NOTES.md
echo "- Cross-platform builds (Linux, macOS, Windows)" >> RELEASE_NOTES.md

# Tag release
git tag "v${VERSION}"

echo "Release v${VERSION} prepared!"
echo "Binaries: dist/"
echo "To publish: git push origin v${VERSION}"
