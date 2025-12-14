package fetcher

import "context"

// Fetcher handles the I/O, just retrieving the raw binary data.
type Fetcher interface {
	Fetch(ctx context.Context) ([]byte, error)
	BaseURL() string
}
