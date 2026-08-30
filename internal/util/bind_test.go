package util

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bindSample struct {
	ID    string `json:"id" validate:"required"`
	Craft string `json:"craft" validate:"required"`
}

func TestFriendlyBindError_RequiredFields(t *testing.T) {
	validate := validator.New()
	err := validate.Struct(bindSample{})
	require.Error(t, err)

	msg := FriendlyBindError(err)
	assert.Contains(t, msg, "name is required")
	assert.Contains(t, msg, "craft is required")
	assert.NotContains(t, msg, "Key:")
	assert.NotContains(t, msg, "CustomRecipeV2")
	assert.NotContains(t, msg, "failed on the")
}

func TestFriendlyBindError_NonValidation(t *testing.T) {
	err := json.Unmarshal([]byte("{"), &map[string]any{})
	require.Error(t, err)
	assert.Equal(t, "invalid request body", FriendlyBindError(err))
}

func TestFriendlyBindError_Nil(t *testing.T) {
	assert.Equal(t, "", FriendlyBindError(nil))
}
