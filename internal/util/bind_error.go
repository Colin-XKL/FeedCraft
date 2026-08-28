package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FriendlyJSONBindError converts gin/validator bind failures into a short
// user-facing message, without leaking struct names or validator tag dumps.
func FriendlyJSONBindError(err error) string {
	if err == nil {
		return "invalid request body"
	}
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		msgs := make([]string, 0, len(verrs))
		for _, fe := range verrs {
			msgs = append(msgs, friendlyValidationMessage(fe))
		}
		if len(msgs) > 0 {
			return strings.Join(msgs, "; ")
		}
	}
	return "invalid request body"
}

func friendlyValidationMessage(fe validator.FieldError) string {
	field := friendlyFieldName(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func friendlyFieldName(field string) string {
	switch field {
	case "ID":
		return "name"
	case "Craft":
		return "craft"
	default:
		if field == "" {
			return "field"
		}
		return strings.ToLower(field)
	}
}
