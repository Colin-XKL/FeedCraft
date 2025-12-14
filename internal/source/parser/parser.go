package parser

import "github.com/mmcdole/gofeed"

// Parser focuses on logic, transforming binary data into a Feed object.
type Parser interface {
	Parse(data []byte) (*gofeed.Feed, error)
}
