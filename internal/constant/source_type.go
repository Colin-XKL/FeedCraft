package constant

// SourceType is a custom type for source identifiers to avoid magic strings.
type SourceType string

const (
	SourceRSS  SourceType = "rss"
	SourceHTML SourceType = "html"
	SourceJSON SourceType = "json"
	// Add other source types here as they are implemented
)
