package download

import "context"

// Downloader handles secure downloads and checksum validation.
type Downloader interface {
	// Download fetches a file from URL and saves to destination.
	Download(ctx context.Context, url, dest string) error

	// DownloadWithChecksum downloads and validates checksum.
	DownloadWithChecksum(ctx context.Context, url, dest, expectedChecksum, algorithm string) error

	// ValidateChecksum verifies file integrity.
	ValidateChecksum(filePath, expectedChecksum, algorithm string) (bool, error)
}

// ChecksumAlgorithm defines supported checksum types.
type ChecksumAlgorithm string

const (
	SHA256 ChecksumAlgorithm = "sha256"
	SHA1   ChecksumAlgorithm = "sha1"
	MD5    ChecksumAlgorithm = "md5"
)

// AtomicInstaller handles atomic filesystem operations for installs.
type AtomicInstaller interface {
	// Install performs atomic installation (temp -> final).
	Install(ctx context.Context, sourceFile, destPath string) error

	// Uninstall safely removes an installed package.
	Uninstall(ctx context.Context, pkgName string) error

	// Rollback reverts a failed installation.
	Rollback(ctx context.Context, installID string) error
}
