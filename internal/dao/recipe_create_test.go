package dao

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateCustomRecipeV2_DuplicateID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&CustomRecipeV2{}))

	recipe := &CustomRecipeV2{
		ID:           "e2e-ai-recipe-887",
		Craft:        "fulltext",
		SourceType:   "rss",
		SourceConfig: `{"http_fetcher":{"url":"https://example.com/rss.xml"}}`,
	}
	require.NoError(t, CreateCustomRecipeV2(db, recipe))

	duplicate := *recipe
	err = CreateCustomRecipeV2(db, &duplicate)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrRecipeAlreadyExists))
	assert.False(t, IsUniqueConstraintError(err), "typed sentinel should not leak SQL unique-constraint text")
}

func TestIsUniqueConstraintError(t *testing.T) {
	assert.False(t, IsUniqueConstraintError(nil))
	assert.False(t, IsUniqueConstraintError(errors.New("connection refused")))
	assert.True(t, IsUniqueConstraintError(gorm.ErrDuplicatedKey))
	assert.True(t, IsUniqueConstraintError(errors.New("constraint failed: UNIQUE constraint failed: custom_recipes_v2.id (1555)")))
}
