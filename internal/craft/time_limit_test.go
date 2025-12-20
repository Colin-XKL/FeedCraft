package craft

import (
	"testing"
	"time"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/assert"
)

func TestOptionTimeLimit(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		days         int
		items        []*feeds.Item
		expectedLen  int
		expectedDesc string // just for description
	}{
		{
			name: "All recent items",
			days: 7,
			items: []*feeds.Item{
				{Title: "Item 1", Created: now},
				{Title: "Item 2", Created: now.AddDate(0, 0, -1)},
			},
			expectedLen: 2,
		},
		{
			name: "Some old items",
			days: 7,
			items: []*feeds.Item{
				{Title: "Recent", Created: now},
				{Title: "Old", Created: now.AddDate(0, 0, -8)}, // 8 days ago
			},
			expectedLen: 1,
		},
		{
			name: "Items with zero date (should be kept)",
			days: 7,
			items: []*feeds.Item{
				{Title: "Recent", Created: now},
				{Title: "Zero", Created: time.Time{}},
			},
			expectedLen: 2,
		},
		{
			name: "All abnormal (1970 or earlier) - should keep all",
			days: 7,
			items: []*feeds.Item{
				{Title: "1970", Created: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Title: "1969", Created: time.Date(1969, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			expectedLen: 2,
		},
		{
			name: "Mixed abnormal and normal - should filter abnormal (assuming normal exists)",
			days: 7,
			items: []*feeds.Item{
				{Title: "Recent", Created: now},
				{Title: "1970", Created: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			expectedLen: 1, // Only Recent should be kept. 1970 is "old" compared to now-7.
		},
		{
			name: "Mixed zero and 1970 - counts as all abnormal/empty? No, logic checks normal date.",
			days: 7,
			items: []*feeds.Item{
				{Title: "Zero", Created: time.Time{}},
				{Title: "1970", Created: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			// Logic: hasNormalDate loop:
			// Zero: IsZero -> continue
			// 1970: !IsZero && Year <= 1970 -> continue
			// hasNormalDate = false.
			// Returns nil (keep all).
			expectedLen: 2,
		},
		{
			name: "Mixed Zero and Recent",
			days: 7,
			items: []*feeds.Item{
				{Title: "Zero", Created: time.Time{}},
				{Title: "Recent", Created: now},
			},
			// Logic: hasNormalDate loop:
			// Zero -> continue
			// Recent -> Normal -> hasNormalDate = true.
			// Filter:
			// Zero -> IsZero -> True (Keep)
			// Recent -> After cutoff -> True (Keep)
			expectedLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed := &feeds.Feed{
				Items: tt.items,
			}
			payload := ExtraPayload{}
			option := OptionTimeLimit(tt.days)

			err := option(feed, payload)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedLen, len(feed.Items))

			// Additional check for "Mixed abnormal and normal" to ensure the right one is kept
			if tt.name == "Mixed abnormal and normal - should filter abnormal (assuming normal exists)" {
				assert.Equal(t, "Recent", feed.Items[0].Title)
			}
		})
	}
}
