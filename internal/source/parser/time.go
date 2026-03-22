package parser

import "time"

var parserTimeLayouts = []string{
	time.RFC3339,
	"2006-01-02",
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC850,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
}

func parseFlexibleTime(value string) (time.Time, bool) {
	for _, layout := range parserTimeLayouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}
