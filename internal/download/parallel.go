package download

import (
	"context"
	"fmt"
	"sync"
)

// ParallelDownloader manages multiple concurrent downloads
type ParallelDownloader struct {
	base       Downloader
	maxWorkers int
}

// NewParallelDownloader creates a parallel downloader
func NewParallelDownloader(base Downloader, maxWorkers int) *ParallelDownloader {
	if maxWorkers <= 0 {
		maxWorkers = 4 // Default to 4 concurrent downloads
	}
	return &ParallelDownloader{
		base:       base,
		maxWorkers: maxWorkers,
	}
}

// DownloadItem represents a file to download
type DownloadItem struct {
	URL    string
	Dest   string
	Hash   string
	Algo   string
	Error  error
}

// DownloadMultiple downloads multiple files in parallel
func (pd *ParallelDownloader) DownloadMultiple(ctx context.Context, items []*DownloadItem) error {
	if len(items) == 0 {
		return nil
	}

	// Channel for work and results
	workChan := make(chan *DownloadItem, len(items))
	resultChan := make(chan *DownloadItem, len(items))

	// Determine worker count
	workers := pd.maxWorkers
	if len(items) < workers {
		workers = len(items)
	}

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go pd.worker(ctx, &wg, workChan, resultChan)
	}

	// Send work
	go func() {
		for _, item := range items {
			workChan <- item
		}
		close(workChan)
	}()

	// Wait for workers
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var errors error
	for result := range resultChan {
		if result.Error != nil {
			errors = fmt.Errorf("download error: %w", result.Error)
		}
	}

	return errors
}

// worker processes download items from the work channel
func (pd *ParallelDownloader) worker(ctx context.Context, wg *sync.WaitGroup, work <-chan *DownloadItem, results chan<- *DownloadItem) {
	defer wg.Done()

	for item := range work {
		if item.Hash != "" && item.Algo != "" {
			item.Error = pd.base.DownloadWithChecksum(ctx, item.URL, item.Dest, item.Hash, item.Algo)
		} else {
			item.Error = pd.base.Download(ctx, item.URL, item.Dest)
		}
		results <- item
	}
}

// DownloadItem pool for reuse (optional optimization)
var downloadItemPool = sync.Pool{
	New: func() interface{} {
		return &DownloadItem{}
	},
}

// GetDownloadItem gets a download item from pool
func GetDownloadItem() *DownloadItem {
	return downloadItemPool.Get().(*DownloadItem)
}

// PutDownloadItem returns a download item to pool
func PutDownloadItem(item *DownloadItem) {
	*item = DownloadItem{} // Reset
	downloadItemPool.Put(item)
}
