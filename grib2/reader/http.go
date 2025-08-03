package reader

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPReaderAt wraps an HTTP client to provide ReaderAt functionality with Range requests
type HTTPReaderAt struct {
	url    string
	client *http.Client
	size   int64
}

func NewHTTPReaderAt(url string) (*HTTPReaderAt, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Get content length with HEAD request
	resp, err := client.Head(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get content length: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP HEAD request failed: %s", resp.Status)
	}

	return &HTTPReaderAt{
		url:    url,
		client: client,
		size:   resp.ContentLength,
	}, nil
}

func (h *HTTPReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= h.size {
		return 0, io.EOF
	}

	// Calculate the range to request
	end := off + int64(len(p)) - 1
	if end >= h.size {
		end = h.size - 1
	}

	// Create HTTP request with Range header
	req, err := http.NewRequest("GET", h.url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", off, end))

	resp, err := h.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP range request failed: %s", resp.Status)
	}

	return io.ReadFull(resp.Body, p[:end-off+1])
}

func (h *HTTPReaderAt) Size() int64 {
	return h.size
}
