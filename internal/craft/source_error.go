package craft

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

func httpStatusForCraftSourceError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	msg := err.Error()
	lower := strings.ToLower(msg)

	switch {
	case errors.Is(err, context.DeadlineExceeded) || strings.Contains(lower, "context deadline exceeded") || strings.Contains(lower, "client.timeout"):
		return http.StatusGatewayTimeout, msg
	case strings.Contains(lower, "parse failed:") || strings.Contains(lower, "failed to detect feed type") || strings.Contains(lower, "invalid xml"):
		return http.StatusUnprocessableEntity, msg
	case strings.Contains(lower, "http status not ok:"),
		strings.Contains(lower, "http get failed:"),
		strings.Contains(lower, "failed to read response body:"),
		strings.Contains(lower, "response body exceeds"):
		return http.StatusBadGateway, msg
	default:
		return http.StatusInternalServerError, msg
	}
}
