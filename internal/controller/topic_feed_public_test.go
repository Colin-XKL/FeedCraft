package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/model"
	"FeedCraft/internal/source"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const topicFeedTestSourceType = constant.SourceType("topic_feed_test_source")

func init() {
	source.Register(topicFeedTestSourceType, func(cfg *config.SourceConfig) (source.Source, error) {
		return &topicFeedTestSource{cfg: cfg}, nil
	})
}

type topicFeedTestSource struct {
	cfg *config.SourceConfig
}

func (s *topicFeedTestSource) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	return &model.CraftFeed{
		Title:       "Test Feed",
		Link:        "https://example.com",
		Id:          "https://example.com/feed.xml",
		Description: "test feed",
		Articles: []*model.CraftArticle{
			{
				Title:       "Hello Topic",
				Link:        "https://example.com/articles/1",
				Id:          "article-1",
				Description: "test article",
				Created:     time.Unix(1700000000, 0),
			},
		},
	}, nil
}

func (s *topicFeedTestSource) BaseURL() string {
	if s.cfg != nil && s.cfg.HttpFetcher != nil && s.cfg.HttpFetcher.URL != "" {
		return s.cfg.HttpFetcher.URL
	}
	return "https://example.com/feed.xml"
}

func TestPublicTopicFeed(t *testing.T) {
	db := topicFeedTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.CustomRecipeV2{}, &dao.TopicFeed{}))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/topic/:id", PublicTopicFeed)

	t.Run("returns rss for topic", func(t *testing.T) {
		recipeID := uniqueTestID("recipe-success")
		topicID := uniqueTestID("topic-success")
		createTopicTestRecipe(t, db, recipeID)
		createTopicTestTopic(t, db, &dao.TopicFeed{
			ID:          topicID,
			Title:       "My Topic",
			Description: "Topic description",
			InputURIs:   []string{"feedcraft://recipe/" + recipeID},
		})

		req, err := http.NewRequest(http.MethodGet, "/topic/"+topicID, nil)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/rss+xml; charset=utf-8", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "<rss")
		assert.Contains(t, recorder.Body.String(), "Hello Topic")
	})

	t.Run("returns 404 when topic is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/topic/"+uniqueTestID("missing"), nil)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assertJSONMessage(t, recorder, "Topic feed not found")
	})

	t.Run("returns 500 when topic build fails", func(t *testing.T) {
		topicID := uniqueTestID("topic-invalid")
		createTopicTestTopic(t, db, &dao.TopicFeed{
			ID:        topicID,
			Title:     "Broken Topic",
			InputURIs: []string{"feedcraft://broken/abc"},
		})

		req, err := http.NewRequest(http.MethodGet, "/topic/"+topicID, nil)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assertJSONMessageContains(t, recorder, "unsupported internal resource type")
	})

	t.Run("returns rss when upstream partially fails", func(t *testing.T) {
		recipeID := uniqueTestID("recipe-partial")
		topicID := uniqueTestID("topic-partial")
		createTopicTestRecipe(t, db, recipeID)
		createTopicTestTopic(t, db, &dao.TopicFeed{
			ID:        topicID,
			Title:     "Partial Topic",
			InputURIs: []string{"feedcraft://recipe/" + recipeID, "http://127.0.0.1:1/unreachable.xml"},
		})

		req, err := http.NewRequest(http.MethodGet, "/topic/"+topicID, nil)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/rss+xml; charset=utf-8", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "Hello Topic")
	})

	t.Run("returns 500 when all upstreams fail", func(t *testing.T) {
		topicID := uniqueTestID("topic-all-failed")
		createTopicTestTopic(t, db, &dao.TopicFeed{
			ID:        topicID,
			Title:     "Failed Topic",
			InputURIs: []string{"http://127.0.0.1:1/a.xml", "http://127.0.0.1:1/b.xml"},
		})

		req, err := http.NewRequest(http.MethodGet, "/topic/"+topicID, nil)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assertJSONMessageContains(t, recorder, "all upstream providers failed")
	})
}

func TestTopicFeedAdminEndpoints(t *testing.T) {
	db := topicFeedTestDatabase(t)
	require.NoError(t, db.AutoMigrate(
		&dao.CustomRecipeV2{},
		&dao.TopicFeed{},
		&dao.ResourceHealth{},
		&dao.ExecutionLog{},
		&dao.SystemNotification{},
	))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/admin/topics/validate", ValidateTopicFeed)
	router.GET("/api/admin/topics/:id/detail", GetTopicFeedDetail)

	t.Run("validate returns success for valid topic config", func(t *testing.T) {
		recipeID := uniqueTestID("recipe-validate-ok")
		createTopicTestRecipe(t, db, recipeID)

		body := `{
			"id":"` + uniqueTestID("topic-validate-ok") + `",
			"title":"Tech Topic",
			"input_uris":["feedcraft://recipe/` + recipeID + `"],
			"aggregator_config":[{"type":"limit","option":{"max":"10"}}]
		}`

		req, err := http.NewRequest(http.MethodPost, "/api/admin/topics/validate", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		var response util.APIResponse[TopicValidationResult]
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		assert.True(t, response.Data.Valid)
		assert.Empty(t, response.Data.Errors)
	})

	t.Run("validate reports invalid aggregator", func(t *testing.T) {
		recipeID := uniqueTestID("recipe-invalid-step")
		createTopicTestRecipe(t, db, recipeID)
		body := `{
			"id":"` + uniqueTestID("topic-invalid-step") + `",
			"title":"Broken Step Topic",
			"input_uris":["feedcraft://recipe/` + recipeID + `"],
			"aggregator_config":[{"type":"limit","option":{"max":"0"}}]
		}`

		req, err := http.NewRequest(http.MethodPost, "/api/admin/topics/validate", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		var response util.APIResponse[TopicValidationResult]
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		assert.False(t, response.Data.Valid)
		require.NotEmpty(t, response.Data.Errors)
		assert.Contains(t, response.Data.Errors[0].Message, "invalid max")
	})

	t.Run("validate reports topic cycle", func(t *testing.T) {
		topicID := uniqueTestID("topic-cycle")
		body := `{
			"id":"` + topicID + `",
			"title":"Cycle Topic",
			"input_uris":["feedcraft://topic/` + topicID + `"],
			"aggregator_config":[{"type":"limit","option":{"max":"10"}}]
		}`

		req, err := http.NewRequest(http.MethodPost, "/api/admin/topics/validate", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		var response util.APIResponse[TopicValidationResult]
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		assert.False(t, response.Data.Valid)
		require.NotEmpty(t, response.Data.Errors)
		assert.Contains(t, response.Data.Errors[0].Message, "topic dependency cycle detected")
	})

	t.Run("detail returns topic config and related runtime records", func(t *testing.T) {
		topicID := uniqueTestID("topic-detail")
		createTopicTestTopic(t, db, &dao.TopicFeed{
			ID:          topicID,
			Title:       "Detail Topic",
			Description: "detail description",
			InputURIs:   []string{"https://example.com/feed.xml"},
			AggregatorConfig: []dao.AggregatorStep{
				{Type: "limit", Option: map[string]string{"max": "20"}},
			},
		})
		now := time.Now()
		require.NoError(t, dao.UpsertResourceHealth(db, &dao.ResourceHealth{
			ResourceType:        dao.ResourceTypeTopic,
			ResourceID:          topicID,
			ResourceName:        "Detail Topic",
			CurrentStatus:       dao.ResourceStatusDegraded,
			ConsecutiveFailures: 2,
			LastFailureAt:       &now,
			LastErrorKind:       "upstream_partial_failure",
			LastErrorMessage:    "partial failure",
		}))
		require.NoError(t, dao.CreateExecutionLog(db, &dao.ExecutionLog{
			ResourceType: dao.ResourceTypeTopic,
			ResourceID:   topicID,
			ResourceName: "Detail Topic",
			Trigger:      "topic_aggregation",
			Status:       dao.ExecutionStatusPartialSuccess,
			ErrorKind:    "upstream_partial_failure",
			Message:      "topic completed with partial upstream failures",
		}))
		require.NoError(t, dao.CreateSystemNotification(db, &dao.SystemNotification{
			ResourceType: dao.ResourceTypeTopic,
			ResourceID:   topicID,
			EventType:    "topic.degraded",
			Title:        "Topic degraded",
			Content:      "detail topic has partial failures",
			StatusAfter:  dao.ResourceStatusDegraded,
			DedupeKey:    uniqueTestID("notification"),
		}))

		req, err := http.NewRequest(http.MethodGet, "/api/admin/topics/"+topicID+"/detail", nil)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
		var response util.APIResponse[TopicDetailResponse]
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
		assert.Equal(t, topicID, response.Data.Topic.ID)
		assert.Equal(t, "/topic/"+topicID, response.Data.PublicURL)
		assert.Equal(t, dao.ResourceStatusDegraded, response.Data.Health.CurrentStatus)
		assert.Len(t, response.Data.RecentExecutions, 1)
		assert.Len(t, response.Data.RelatedNotifications, 1)
	})

	t.Run("detail returns 404 when topic is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/admin/topics/"+uniqueTestID("topic-missing")+"/detail", nil)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assertJSONMessage(t, recorder, "Topic feed not found")
	})
}

func topicFeedTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	util.SetDatabaseForTest(db)
	return db
}

func createTopicTestRecipe(t *testing.T, db *gorm.DB, recipeID string) {
	t.Helper()

	sourceConfig := `{"type":"` + string(topicFeedTestSourceType) + `","http_fetcher":{"url":"https://example.com/feed.xml"}}`
	require.NoError(t, dao.CreateCustomRecipeV2(db, &dao.CustomRecipeV2{
		ID:           recipeID,
		Craft:        "proxy",
		SourceType:   string(topicFeedTestSourceType),
		SourceConfig: sourceConfig,
	}))
}

func createTopicTestTopic(t *testing.T, db *gorm.DB, topic *dao.TopicFeed) {
	t.Helper()
	require.NoError(t, dao.CreateTopicFeed(db, topic))
}

func assertJSONMessage(t *testing.T, recorder *httptest.ResponseRecorder, expected string) {
	t.Helper()
	var response util.APIResponse[any]
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, expected, response.Msg)
}

func assertJSONMessageContains(t *testing.T, recorder *httptest.ResponseRecorder, expected string) {
	t.Helper()
	var response util.APIResponse[any]
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Contains(t, response.Msg, expected)
}

func uniqueTestID(prefix string) string {
	return prefix + "-" + time.Now().Format("20060102150405.000000000")
}
