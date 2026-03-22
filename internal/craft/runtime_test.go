package craft

import (
	"context"
	"testing"
	"time"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/engine"
	"FeedCraft/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestResolveCraftAtoms_FlowAndCustomAtom(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	require.NoError(t, dao.CreateCraftAtom(db, &dao.CraftAtom{
		Name:         "limit-five",
		TemplateName: "limit",
		Params:       map[string]string{"num": "5"},
	}))
	require.NoError(t, dao.CreateCraftFlow(db, &dao.CraftFlow{
		Name: "flow-a",
		CraftFlowConfig: []dao.CraftFlowItem{
			{CraftName: "limit-five"},
			{CraftName: "guid-fix"},
		},
	}))

	atoms, err := ResolveCraftAtoms(db, "flow-a")
	require.NoError(t, err)
	require.Len(t, atoms, 2)
	assert.Equal(t, "limit", atoms[0].TemplateName)
	assert.Equal(t, "guid-fix", atoms[1].TemplateName)
	assert.Equal(t, "5", atoms[0].Params["num"])
}

func TestBuildProcessor_UsesNativeProcessors(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	processor, err := BuildProcessor(db, "limit,time-limit,guid-fix,relative-link-fix", "https://example.com/feed.xml")
	require.NoError(t, err)
	require.NotNil(t, processor)

	flow, ok := processor.(*engine.FlowCraftProcessor)
	require.True(t, ok)
	require.Len(t, flow.Processors, 4)
	assert.IsType(t, &LimitProcessor{}, flow.Processors[0])
	assert.IsType(t, &TimeLimitProcessor{}, flow.Processors[1])
	assert.IsType(t, &GUIDFixProcessor{}, flow.Processors[2])
	assert.IsType(t, &RelativeLinkFixProcessor{}, flow.Processors[3])
}

func TestBuildProcessor_FallsBackToLegacyForUnsupportedAtom(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	processor, err := BuildProcessor(db, "cleanup", "https://example.com/feed.xml")
	require.NoError(t, err)
	require.NotNil(t, processor)

	flow, ok := processor.(*engine.FlowCraftProcessor)
	require.True(t, ok)
	require.Len(t, flow.Processors, 1)
	_, isNestedFlow := flow.Processors[0].(*engine.FlowCraftProcessor)
	assert.True(t, isNestedFlow)
}

func TestNativeProcessors_EndToEnd(t *testing.T) {
	now := time.Now()
	processor := &engine.FlowCraftProcessor{
		Processors: []engine.FeedProcessor{
			&KeywordProcessor{
				Mode:       KeywordIncludeMode,
				MatchScope: KeywordMatchAll,
				Keywords:   []string{"keep"},
			},
			&GUIDFixProcessor{},
			&RelativeLinkFixProcessor{OriginalFeedURL: "https://example.com/feed.xml"},
			&LimitProcessor{MaxItems: 1},
			&TimeLimitProcessor{
				Days: 7,
				Now:  func() time.Time { return now },
			},
		},
	}

	feed := &model.CraftFeed{
		Link: "https://example.com",
		Articles: []*model.CraftArticle{
			{
				Title:       "keep article",
				Content:     "keep this",
				Description: "keep this",
				Link:        "/article-1",
				Created:     now,
			},
			{
				Title:       "drop article",
				Content:     "drop this",
				Description: "drop this",
				Link:        "/article-2",
				Created:     now,
			},
		},
	}

	result, err := processor.Process(context.Background(), feed)
	require.NoError(t, err)
	require.Len(t, result.Articles, 1)
	assert.Equal(t, "https://example.com/article-1", result.Articles[0].Link)
	assert.NotEmpty(t, result.Articles[0].Id)
}

func newCraftRuntimeTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&dao.CraftAtom{}, &dao.CraftFlow{}))
	return db
}
