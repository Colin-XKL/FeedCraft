package craft

import (
	"FeedCraft/internal/engine"
	"FeedCraft/internal/util"

	"gorm.io/gorm"
)

// BuildProcessor compiles a craft DSL string into a runtime processor chain.
func BuildProcessor(db *gorm.DB, craftName string, feedURL string) (engine.FeedProcessor, error) {
	if db == nil {
		db = util.GetDatabase()
	}

	craftOptionList, err := getCraftOptions(db, craftName)
	if err != nil {
		return nil, err
	}
	if len(craftOptionList) == 0 {
		return nil, nil
	}

	payload := ExtraPayload{originalFeedUrl: feedURL}
	processors := make([]engine.FeedProcessor, 0, len(craftOptionList))
	for _, opt := range craftOptionList {
		processors = append(processors, &LegacyOptionAdapter{
			Option: opt,
			Extra:  payload,
		})
	}

	return &engine.FlowCraftProcessor{Processors: processors}, nil
}
