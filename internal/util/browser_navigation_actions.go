package util

import (
	"FeedCraft/internal/config"
	"fmt"
	"strings"
)

const (
	MaxBrowserNavigationActions       = 20
	MaxBrowserNavigationTimeoutMs     = 30000
	DefaultBrowserNavigationTimeoutMs = 5000
	MaxBrowserNavigationWaitMs        = 10000
	MaxBrowserNavigationBudgetMs      = 60000
)

func ValidateBrowserNavigationActions(actions []config.BrowserNavigationAction) error {
	if len(actions) > MaxBrowserNavigationActions {
		return fmt.Errorf("navigation actions cannot exceed %d", MaxBrowserNavigationActions)
	}
	var budgetMs int64
	for idx, action := range actions {
		actionType := strings.TrimSpace(action.Type)
		switch actionType {
		case config.BrowserNavigationActionClick, config.BrowserNavigationActionWaitForSelector:
			selector := strings.TrimSpace(action.Selector)
			if selector == "" {
				return fmt.Errorf("navigation action %d %s selector is required", idx+1, actionType)
			}
			if action.TimeoutMs < 0 || action.TimeoutMs > MaxBrowserNavigationTimeoutMs {
				return fmt.Errorf("navigation action %d timeout_ms must be between 0 and %d", idx+1, MaxBrowserNavigationTimeoutMs)
			}
			effectiveTimeoutMs := action.TimeoutMs
			if effectiveTimeoutMs == 0 {
				effectiveTimeoutMs = DefaultBrowserNavigationTimeoutMs
			}
			budgetMs += effectiveTimeoutMs
			actions[idx].Type = actionType
			actions[idx].Selector = selector
			actions[idx].TimeoutMs = effectiveTimeoutMs
		case config.BrowserNavigationActionWait:
			if action.DurationMs <= 0 || action.DurationMs > MaxBrowserNavigationWaitMs {
				return fmt.Errorf("navigation action %d wait duration_ms must be between 1 and %d", idx+1, MaxBrowserNavigationWaitMs)
			}
			budgetMs += action.DurationMs
			actions[idx].Type = actionType
		default:
			return fmt.Errorf("navigation action %d has unsupported type %q", idx+1, action.Type)
		}
	}
	if budgetMs > MaxBrowserNavigationBudgetMs {
		return fmt.Errorf("navigation actions total wait budget cannot exceed %d ms", MaxBrowserNavigationBudgetMs)
	}
	return nil
}
