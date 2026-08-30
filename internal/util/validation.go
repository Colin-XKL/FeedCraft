package util

import "regexp"

var idRegex = regexp.MustCompile(`^[a-z0-9-_]+$`)

// IsValidID checks if the given string consists of only lowercase letters, digits, hyphens, and underscores.
func IsValidID(id string) bool {
	return idRegex.MatchString(id)
}
