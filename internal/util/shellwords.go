package util

import (
	"errors"
)

// ParseShellWords parses a command line string into arguments, handling quotes and escapes.
// This is a simplified implementation inspired by github.com/mattn/go-shellwords,
// but without environment variable expansion or shell execution (backticks) for security.
func ParseShellWords(line string) ([]string, error) {
	var args []string
	var buf string
	var escaped, doubleQuoted, singleQuoted bool
	gotArg := false

	for _, r := range line {
		if escaped {
			buf += string(r)
			escaped = false
			gotArg = true
			continue
		}

		if r == '\\' {
			if singleQuoted {
				buf += string(r)
				gotArg = true
			} else {
				escaped = true
			}
			continue
		}

		if isSpace(r) {
			if singleQuoted || doubleQuoted {
				buf += string(r)
			} else if gotArg {
				args = append(args, buf)
				buf = ""
				gotArg = false
			}
			continue
		}

		switch r {
		case '"':
			if !singleQuoted {
				doubleQuoted = !doubleQuoted
				gotArg = true
				continue
			}
		case '\'':
			if !doubleQuoted {
				singleQuoted = !singleQuoted
				gotArg = true
				continue
			}
		}

		buf += string(r)
		gotArg = true
	}

	if gotArg {
		args = append(args, buf)
	}

	if escaped || singleQuoted || doubleQuoted {
		return nil, errors.New("invalid command line string: unmatched quotes or escape")
	}

	return args, nil
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}
