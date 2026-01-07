package util

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type CurlParseResult struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// ParseCurlCommand parses a curl command string into a structured request object.
// It supports common curl flags like -X, -H, -d, --data-raw, etc.
func ParseCurlCommand(cmd string) (*CurlParseResult, error) {
	args, err := ParseShellWords(cmd)
	if err != nil {
		return nil, fmt.Errorf("parse curl command failed: %w", err)
	}

	if len(args) == 0 || args[0] != "curl" {
		return nil, fmt.Errorf("invalid curl command")
	}

	result := &CurlParseResult{
		Headers: make(map[string]string),
		Method:  "GET", // Default
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]

		// Handle flags
		if strings.HasPrefix(arg, "-") {
			// Check for attached values in short flags (e.g., -XPOST)
			// Note: shellwords might have already split them if there was a space,
			// but if the user typed "curl -XPOST", shellwords returns ["curl", "-XPOST"].
			// We need to handle this.

			// Helper to get value: either from attached part or next arg
			getValue := func(flag string) (string, bool) {
				if arg == flag {
					// Value is in next arg
					if i+1 < len(args) {
						i++
						return args[i], true
					}
					return "", false
				}
				if strings.HasPrefix(arg, flag) {
					// Value is attached (only for short flags usually)
					return arg[len(flag):], true
				}
				return "", false
			}

			// Method
			if val, ok := getValue("-X"); ok {
				result.Method = val
				continue
			}
			if val, ok := getValue("--request"); ok {
				result.Method = val
				continue
			}

			// Header
			if val, ok := getValue("-H"); ok {
				parseHeader(result, val)
				continue
			}
			if val, ok := getValue("--header"); ok {
				parseHeader(result, val)
				continue
			}

			// Data / Body
			// -d, --data, --data-ascii, --data-raw, --data-binary, --data-urlencode
			// Note: Check longer flags first to avoid partial prefix matching (e.g. --data matching --data-raw)
			dataFlags := []string{"--data-raw", "--data-binary", "--data-ascii", "--data-urlencode", "--data", "-d"}
			matchedData := false
			for _, flag := range dataFlags {
				if val, ok := getValue(flag); ok {
					// If multiple data flags are used, curl concatenates them with "&".
					// For simplicity, we just take the last one or append?
					// Curl docs say: "If you start the data with the letter @, the rest should be a file name...". We ignore file references for security/complexity.
					if result.Body != "" {
						result.Body += "&" + val
					} else {
						result.Body = val
					}

					// Implicit POST if not set
					if result.Method == "GET" {
						result.Method = "POST"
					}
					matchedData = true
					break
				}
			}
			if matchedData {
				continue
			}

			// User/Basic Auth (-u user:password)
			if val, ok := getValue("-u"); ok {
				encoded := base64.StdEncoding.EncodeToString([]byte(val))
				result.Headers["Authorization"] = "Basic " + encoded
				continue
			}
			if val, ok := getValue("--user"); ok {
				encoded := base64.StdEncoding.EncodeToString([]byte(val))
				result.Headers["Authorization"] = "Basic " + encoded
				continue
			}

			// URL explicit flag
			if val, ok := getValue("--url"); ok {
				result.URL = val
				continue
			}

			// Compressed
			if arg == "--compressed" {
				// We don't need to do anything, resty handles compression usually, or we just ignore it so it doesn't break
				continue
			}

			// Ignore other flags
			continue
		}

		// Positional argument (URL)
		// Only if URL is empty (first positional is usually URL)
		if result.URL == "" {
			result.URL = arg
		}
	}

	return result, nil
}

func parseHeader(result *CurlParseResult, headerStr string) {
	parts := strings.SplitN(headerStr, ":", 2)
	if len(parts) == 2 {
		result.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	} else if len(parts) == 1 {
		// Empty header or malformed
		// curl -H "X-Empty;" removes the header.
		// curl -H "Key" sends "Key;" (empty value) logic varies.
		// For simplicity, treat as empty value
		key := strings.TrimSpace(parts[0])
		if strings.HasSuffix(key, ";") {
			result.Headers[strings.TrimSuffix(key, ";")] = ""
		} else {
			result.Headers[key] = ""
		}
	}
}
