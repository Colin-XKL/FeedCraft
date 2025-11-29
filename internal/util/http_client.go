package util

import (
	"context"
	"net/http"
	"time"
)

// Timeout Constants
const (
	ExternalRequestTimeout = 1 * time.Minute
	LLMRequestTimeout      = 5 * time.Minute
)

// Retry Configuration
const (
	MaxRetries       = 3
	RetryWaitTime    = 1 * time.Second
	RetryMaxWaitTime = 5 * time.Second
)

// HTTPClientWithTimeout returns an http.Client with the specified timeout.
func HTTPClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

// CreateContextWithTimeout creates a context with a timeout.
func CreateContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
