package download_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/grave0x/gpkg/internal/download"
)

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

func TestParallelDownloaderSuccess(t *testing.T) {
	mock := &MockDownloader{}
	pd := download.NewParallelDownloader(mock, 2)

	items := []*download.DownloadItem{
		{URL: "https://example.com/pkg1.tar.gz", Dest: "/tmp/pkg1.tar.gz"},
		{URL: "https://example.com/pkg2.tar.gz", Dest: "/tmp/pkg2.tar.gz"},
		{URL: "https://example.com/pkg3.tar.gz", Dest: "/tmp/pkg3.tar.gz"},
	}

	ctx := context.Background()
	err := pd.DownloadMultiple(ctx, items)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mock.downloaded) != 3 {
		t.Errorf("expected 3 downloads, got %d", len(mock.downloaded))
	}
}

func TestParallelDownloaderMaxWorkers(t *testing.T) {
	mock := &MockDownloader{}
	pd := download.NewParallelDownloader(mock, 1) // Single worker

	items := []*download.DownloadItem{
		{URL: "https://example.com/file1", Dest: filepath.Join(t.TempDir(), "file1")},
		{URL: "https://example.com/file2", Dest: filepath.Join(t.TempDir(), "file2")},
	}

	ctx := context.Background()
	err := pd.DownloadMultiple(ctx, items)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(mock.downloaded) != 2 {
		t.Errorf("expected 2 downloads with single worker, got %d", len(mock.downloaded))
	}
}

func TestParallelDownloaderEmpty(t *testing.T) {
	mock := &MockDownloader{}
	pd := download.NewParallelDownloader(mock, 2)

	ctx := context.Background()
	err := pd.DownloadMultiple(ctx, []*download.DownloadItem{})

	if err != nil {
		t.Errorf("unexpected error for empty list: %v", err)
	}

	if len(mock.downloaded) != 0 {
		t.Errorf("expected no downloads for empty list")
	}
}

func TestDownloadItemPool(t *testing.T) {
	item1 := download.GetDownloadItem()
	if item1 == nil {
		t.Errorf("expected item from pool")
	}

	item1.URL = "https://example.com/test"
	download.PutDownloadItem(item1)

	item2 := download.GetDownloadItem()
	if item2.URL != "" {
		t.Errorf("expected reset item from pool, but URL was not cleared")
	}
}
