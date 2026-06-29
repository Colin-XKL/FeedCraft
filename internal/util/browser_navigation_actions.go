package util

import (
	"FeedCraft/internal/config"
	"fmt"
	"strings"
)

func ValidateBrowserNavigationActions(actions []config.BrowserNavigationAction) error {
	for idx, action := range actions {
		actionType := strings.TrimSpace(action.Type)
		switch actionType {
		case config.BrowserNavigationActionClick, config.BrowserNavigationActionWaitForSelector:
			if strings.TrimSpace(action.Selector) == "" {
				return fmt.Errorf("navigation action %d %s selector is required", idx+1, actionType)
			}
		case config.BrowserNavigationActionWait:
			if action.DurationMs <= 0 {
				return fmt.Errorf("navigation action %d wait duration_ms must be greater than 0", idx+1)
			}
		default:
			return fmt.Errorf("navigation action %d has unsupported type %q", idx+1, action.Type)
		}
	}
	return nil
}
