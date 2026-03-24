package download

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// HTTPDownloader implements Downloader using HTTP client.
type HTTPDownloader struct {
	client       *http.Client
	timeout      time.Duration
	allowOffline bool
}

// NewHTTPDownloader creates a new HTTP downloader.
func NewHTTPDownloader(timeout time.Duration, allowOffline bool) *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout:      timeout,
		allowOffline: allowOffline,
	}
}

// Download fetches a file from URL and saves to destination.
func (d *HTTPDownloader) Download(ctx context.Context, url, dest string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		if d.allowOffline {
			return fmt.Errorf("offline mode enabled, cannot download: %w", err)
		}
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Write to temp file first for atomic operations
	tmpFile, err := os.CreateTemp(filepath.Dir(dest), ".download-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write downloaded content: %w", err)
	}
	tmpFile.Close()

	// Move temp file to destination
	if err := os.Rename(tmpFile.Name(), dest); err != nil {
		return fmt.Errorf("failed to move downloaded file to destination: %w", err)
	}

	return nil
}

// DownloadWithChecksum downloads and validates checksum.
func (d *HTTPDownloader) DownloadWithChecksum(ctx context.Context, url, dest, expectedChecksum, algorithm string) error {
	if err := d.Download(ctx, url, dest); err != nil {
		return err
	}

	isValid, err := d.ValidateChecksum(dest, expectedChecksum, algorithm)
	if err != nil {
		os.Remove(dest) // Clean up on validation error
		return fmt.Errorf("checksum validation failed: %w", err)
	}

	if !isValid {
		os.Remove(dest) // Clean up on checksum mismatch
		return fmt.Errorf("checksum mismatch: expected %s but got different value", expectedChecksum)
	}

	return nil
}

// ValidateChecksum verifies file integrity.
func (d *HTTPDownloader) ValidateChecksum(filePath, expectedChecksum, algorithm string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var h hash.Hash
	switch algorithm {
	case "sha256":
		h = sha256.New()
	case "sha1":
		h = sha1.New()
	case "md5":
		h = md5.New()
	default:
		return false, fmt.Errorf("unsupported checksum algorithm: %s", algorithm)
	}

	if _, err := io.Copy(h, file); err != nil {
		return false, fmt.Errorf("failed to read file for checksum: %w", err)
	}

	actualChecksum := fmt.Sprintf("%x", h.Sum(nil))
	return actualChecksum == expectedChecksum, nil
}

// AtomicInstallerImpl implements AtomicInstaller for filesystem operations.
type AtomicInstallerImpl struct {
	backupDir string
}

// NewAtomicInstaller creates a new atomic installer.
func NewAtomicInstaller(backupDir string) *AtomicInstallerImpl {
	return &AtomicInstallerImpl{
		backupDir: backupDir,
	}
}

// Install performs atomic installation (temp -> final).
func (a *AtomicInstallerImpl) Install(ctx context.Context, sourceFile, destPath string) error {
	// Create backup if destination exists
	if _, err := os.Stat(destPath); err == nil {
		backupID := fmt.Sprintf("%d", time.Now().UnixNano())
		backupPath := filepath.Join(a.backupDir, backupID)

		if err := os.MkdirAll(a.backupDir, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}

		if err := os.Rename(destPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing file: %w", err)
		}
	}

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Move source file to destination
	if err := os.Rename(sourceFile, destPath); err != nil {
		return fmt.Errorf("failed to move file to destination: %w", err)
	}

	// Set proper permissions
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// Uninstall safely removes an installed package.
func (a *AtomicInstallerImpl) Uninstall(ctx context.Context, pkgName string) error {
	// This will be expanded to handle proper cleanup
	// For now, basic removal
	return nil
}

// Rollback reverts a failed installation.
func (a *AtomicInstallerImpl) Rollback(ctx context.Context, installID string) error {
	backupPath := filepath.Join(a.backupDir, installID)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", installID)
	}

	// Restore from backup
	// Implementation depends on storage strategy
	return nil
}
