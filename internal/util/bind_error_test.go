package util

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type bindSample struct {
	ID    string `validate:"required"`
	Craft string `validate:"required"`
}

func TestFriendlyJSONBindErrorRequiredFields(t *testing.T) {
	validate := validator.New()
	err := validate.Struct(bindSample{})
	require.Error(t, err)

	msg := FriendlyJSONBindError(err)
	require.Contains(t, msg, "name is required")
	require.Contains(t, msg, "craft is required")
	require.NotContains(t, msg, "CustomRecipeV2")
	require.NotContains(t, msg, "Field validation")
}

func TestFriendlyJSONBindErrorNonValidator(t *testing.T) {
	require.Equal(t, "invalid request body", FriendlyJSONBindError(nil))
	require.Equal(t, "invalid request body", FriendlyJSONBindError(assertErr("syntax")))
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
