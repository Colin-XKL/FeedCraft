package dao

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// ErrRecipeAlreadyExists is returned when creating a recipe whose id is already taken.
var ErrRecipeAlreadyExists = errors.New("a recipe with this name already exists")

// IsUniqueConstraintError reports whether err is a SQL unique/primary-key conflict.
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "duplicated key")
}
