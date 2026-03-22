package parser

import "FeedCraft/internal/model"

// Parser focuses on logic, transforming binary data into a CraftFeed.
type Parser interface {
	Parse(data []byte) (*model.CraftFeed, error)
}
