package controller

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type DependencyStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "Healthy", "Unhealthy", "Configured", "Not Configured"
	Details string `json:"details"`
	Error   string `json:"error,omitempty"`
	Latency string `json:"latency,omitempty"`
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "******"
	}
	return s[:4] + "******" + s[len(s)-4:]
}

func maskURI(uri string) string {
	if strings.HasPrefix(uri, "redis://") {
		parts := strings.Split(uri, "@")
		if len(parts) > 1 {
			auth := strings.Split(parts[0], "//")
			if len(auth) > 1 {
				return auth[0] + "//******@" + parts[1]
			}
		}
	}
	return uri
}

func checkSQLite(env *viper.Viper, activeCheck bool) DependencyStatus {
	sqlitePath := env.GetString("DB_SQLITE_PATH")
	if sqlitePath == "" {
		return DependencyStatus{Name: "SQLite", Status: "Not Configured", Error: "FC_DB_SQLITE_PATH not set"}
	}
	path := filepath.Join(sqlitePath, "feed-craft.db")
	details := fmt.Sprintf("Path: %s", path)

	if !activeCheck {
		return DependencyStatus{Name: "SQLite", Status: "Configured", Details: details}
	}

	start := time.Now()
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return DependencyStatus{Name: "SQLite", Status: "Unhealthy", Details: details, Error: err.Error()}
	}
	sqlDB, err := db.DB()
	if err != nil {
		return DependencyStatus{Name: "SQLite", Status: "Unhealthy", Details: details, Error: err.Error()}
	}
	defer func() { _ = sqlDB.Close() }()

	if err := sqlDB.Ping(); err != nil {
		return DependencyStatus{Name: "SQLite", Status: "Unhealthy", Details: details, Error: err.Error()}
	}
	return DependencyStatus{Name: "SQLite", Status: "Healthy", Details: details, Latency: time.Since(start).String()}
}

func checkRedis(env *viper.Viper, activeCheck bool) DependencyStatus {
	redisURI := env.GetString("REDIS_URI")
	if redisURI == "" {
		return DependencyStatus{Name: "Redis", Status: "Not Configured", Error: "FC_REDIS_URI not set"}
	}
	details := fmt.Sprintf("URI: %s", maskURI(redisURI))

	if !activeCheck {
		return DependencyStatus{Name: "Redis", Status: "Configured", Details: details}
	}

	start := time.Now()
	opts, err := redis.ParseURL(redisURI)
	if err != nil {
		return DependencyStatus{Name: "Redis", Status: "Unhealthy", Details: details, Error: "Invalid URI: " + err.Error()}
	}
	rdb := redis.NewClient(opts)
	defer func() { _ = rdb.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return DependencyStatus{Name: "Redis", Status: "Unhealthy", Details: details, Error: err.Error()}
	}
	return DependencyStatus{Name: "Redis", Status: "Healthy", Details: details, Latency: time.Since(start).String()}
}

func checkBrowserless(env *viper.Viper, activeCheck bool) DependencyStatus {
	endpoint := env.GetString("PUPPETEER_HTTP_ENDPOINT")
	if endpoint == "" {
		return DependencyStatus{Name: "Browserless", Status: "Not Configured"}
	}
	details := fmt.Sprintf("Endpoint: %s", endpoint)

	if !activeCheck {
		return DependencyStatus{Name: "Browserless", Status: "Configured", Details: details}
	}

	start := time.Now()
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		return DependencyStatus{Name: "Browserless", Status: "Unhealthy", Details: details, Error: err.Error()}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return DependencyStatus{Name: "Browserless", Status: "Healthy", Details: details, Latency: time.Since(start).String()}
	}
	return DependencyStatus{Name: "Browserless", Status: "Unhealthy", Details: details, Error: resp.Status}
}

func checkLLM(env *viper.Viper, activeCheck bool) DependencyStatus {
	llmBase := env.GetString("LLM_API_BASE")
	// For Ollama, API Key might be empty
	llmKey := env.GetString("LLM_API_KEY")
	llmType := env.GetString("LLM_API_TYPE")
	if llmType == "" {
		llmType = "openai"
	}
	llmModel := env.GetString("LLM_API_MODEL")

	if llmBase == "" {
		return DependencyStatus{Name: "LLM Service", Status: "Not Configured", Error: "FC_LLM_API_BASE not set"}
	}

	details := fmt.Sprintf("Type: %s, Model: %s, Base: %s", llmType, llmModel, llmBase)
	if llmKey != "" {
		details += fmt.Sprintf(", Key: %s", maskString(llmKey))
	}

	if !activeCheck {
		return DependencyStatus{Name: "LLM Service", Status: "Configured", Details: details}
	}

	start := time.Now()
	// Reuse SimpleLLMCall logic
	// We use a simple math query to verify functionality
	_, err := adapter.SimpleLLMCall(adapter.UseDefaultModel, "1+1=?")
	if err != nil {
		return DependencyStatus{Name: "LLM Service", Status: "Unhealthy", Details: details, Error: err.Error()}
	}
	return DependencyStatus{Name: "LLM Service", Status: "Healthy", Details: details, Latency: time.Since(start).String()}
}

func getStatuses(activeCheck bool) []DependencyStatus {
	envClient := util.GetEnvClient()
	return []DependencyStatus{
		checkSQLite(envClient, activeCheck),
		checkRedis(envClient, activeCheck),
		checkBrowserless(envClient, activeCheck),
		checkLLM(envClient, activeCheck),
	}
}

func GetDependencyStatus(c *gin.Context) {
	statuses := getStatuses(false)
	c.JSON(http.StatusOK, util.APIResponse[[]DependencyStatus]{
		StatusCode: 0,
		Msg:        "Success",
		Data:       statuses,
	})
}

func CheckDependencyStatus(c *gin.Context) {
	statuses := getStatuses(true)
	c.JSON(http.StatusOK, util.APIResponse[[]DependencyStatus]{
		StatusCode: 0,
		Msg:        "Success",
		Data:       statuses,
	})
}
