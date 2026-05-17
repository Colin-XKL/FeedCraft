package controller

import (
	"errors"
	"net/http"
	"testing"
)

func TestFormatFeedViewerValidationErrorPreservesUserFacingCapitalization(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "invalid URL",
			err:  errors.New("please enter a valid http(s) feed URL"),
			want: "Please enter a valid http(s) feed URL",
		},
		{
			name: "private IP",
			err:  errors.New("access to private IP 127.0.0.1 is forbidden"),
			want: "Access to private IP 127.0.0.1 is forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFeedViewerValidationError(tt.err)
			if got != tt.want {
				t.Fatalf("formatFeedViewerValidationError() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClassifyFeedViewerErrorHandlesLowercaseResolveMessage(t *testing.T) {
	status, msg := classifyFeedViewerError(errors.New("unable to resolve this URL: lookup example.invalid: no such host"))

	if status != http.StatusOK {
		t.Fatalf("status = %d, want %d", status, http.StatusOK)
	}
	const want = "Unable to fetch this URL. Please check the address and try again."
	if msg != want {
		t.Fatalf("msg = %q, want %q", msg, want)
	}
}

func TestClassifyFeedViewerErrorExplainsEmbeddingConfiguration(t *testing.T) {
	status, msg := classifyFeedViewerError(errors.New("[embedding-filter] failed to compute anchor vectors: failed to load embedding config: FC_EMBEDDING_API_MODEL must be set when using FC_EMBEDDING_API_TYPE='ollama'"))

	if status != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", status, http.StatusBadRequest)
	}
	const want = "Embedding filter is not configured correctly: FC_EMBEDDING_API_MODEL must be set when using FC_EMBEDDING_API_TYPE='ollama'"
	if msg != want {
		t.Fatalf("msg = %q, want %q", msg, want)
	}
}

func TestClassifyFeedViewerErrorDoesNotExposeEmbeddingRuntimeError(t *testing.T) {
	status, msg := classifyFeedViewerError(errors.New("[embedding-filter] all article embeddings failed: embedding call failed after retries (batch [0-1]): provider returned 500 with token detail"))

	if status != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", status, http.StatusInternalServerError)
	}
	const want = "Failed to preview this feed due to an internal error."
	if msg != want {
		t.Fatalf("msg = %q, want %q", msg, want)
	}
}
