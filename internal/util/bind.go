package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FriendlyBindError turns gin/validator bind failures into a short, user-facing
// message instead of dumping struct tags like `Key: 'CustomRecipeV2.ID'`.
func FriendlyBindError(err error) string {
	if err == nil {
		return ""
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) && len(ve) > 0 {
		parts := make([]string, 0, len(ve))
		for _, fieldErr := range ve {
			parts = append(parts, formatFieldValidation(fieldErr))
		}
		return strings.Join(parts, "; ")
	}

	return "invalid request body"
}

func formatFieldValidation(fieldErr validator.FieldError) string {
	name := friendlyFieldName(fieldErr.Field())
	if fieldErr.Tag() == "required" {
		return name + " is required"
	}
	return fmt.Sprintf("%s is invalid", name)
}

func friendlyFieldName(field string) string {
	switch field {
	case "ID":
		return "name"
	case "Craft":
		return "craft"
	default:
		return strings.ToLower(field)
	}
}
