package controller

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type DependencyStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "Healthy", "Unhealthy", "Not Configured"
	Details string `json:"details"`
	Error   string `json:"error,omitempty"`
	Latency string `json:"latency,omitempty"`
}

func maskString(s string) string {
	if len(s) < 8 {
		return "******"
	}
	return s[:4] + "******" + s[len(s)-4:]
}

func maskURI(uri string) string {
	// redis://user:password@host:port/db
	// we want to mask password
	if strings.HasPrefix(uri, "redis://") {
		parts := strings.Split(uri, "@")
		if len(parts) > 1 {
			// user:password part
			auth := strings.Split(parts[0], "//")
			if len(auth) > 1 {
				return auth[0] + "//******@" + parts[1]
			}
		}
	}
	return uri // Return as is if format not recognized or no password likely
}

func GetDependencyStatus(c *gin.Context) {
	var statuses []DependencyStatus
	envClient := util.GetEnvClient()

	// 1. SQLite Check
	sqlitePath := envClient.GetString("DB_SQLITE_PATH")
	if sqlitePath == "" {
		statuses = append(statuses, DependencyStatus{Name: "SQLite", Status: "Not Configured", Error: "FC_DB_SQLITE_PATH not set"})
	} else {
		start := time.Now()
		path := filepath.Join(sqlitePath, "feed-craft.db")
		details := fmt.Sprintf("Path: %s", path)
		// Try to open connection safely
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			statuses = append(statuses, DependencyStatus{Name: "SQLite", Status: "Unhealthy", Details: details, Error: err.Error()})
		} else {
			sqlDB, err := db.DB()
			if err != nil {
				statuses = append(statuses, DependencyStatus{Name: "SQLite", Status: "Unhealthy", Details: details, Error: err.Error()})
			} else {
				if err := sqlDB.Ping(); err != nil {
					statuses = append(statuses, DependencyStatus{Name: "SQLite", Status: "Unhealthy", Details: details, Error: err.Error()})
				} else {
					statuses = append(statuses, DependencyStatus{Name: "SQLite", Status: "Healthy", Details: details, Latency: time.Since(start).String()})
				}
				sqlDB.Close()
			}
		}
	}

	// 2. Redis Check
	redisURI := envClient.GetString("REDIS_URI")
	if redisURI == "" {
		statuses = append(statuses, DependencyStatus{Name: "Redis", Status: "Not Configured", Error: "FC_REDIS_URI not set"})
	} else {
		start := time.Now()
		details := fmt.Sprintf("URI: %s", maskURI(redisURI))
		opts, err := redis.ParseURL(redisURI)
		if err != nil {
			statuses = append(statuses, DependencyStatus{Name: "Redis", Status: "Unhealthy", Details: details, Error: "Invalid URI: " + err.Error()})
		} else {
			rdb := redis.NewClient(opts)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := rdb.Ping(ctx).Err(); err != nil {
				statuses = append(statuses, DependencyStatus{Name: "Redis", Status: "Unhealthy", Details: details, Error: err.Error()})
			} else {
				statuses = append(statuses, DependencyStatus{Name: "Redis", Status: "Healthy", Details: details, Latency: time.Since(start).String()})
			}
			rdb.Close()
		}
	}

	// 3. Browserless Check
	browserlessEndpoint := envClient.GetString("PUPPETEER_HTTP_ENDPOINT")
	if browserlessEndpoint == "" {
		statuses = append(statuses, DependencyStatus{Name: "Browserless", Status: "Not Configured"})
	} else {
		start := time.Now()
		details := fmt.Sprintf("Endpoint: %s", browserlessEndpoint)
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(browserlessEndpoint)
		if err != nil {
			statuses = append(statuses, DependencyStatus{Name: "Browserless", Status: "Unhealthy", Details: details, Error: err.Error()})
		} else {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 500 {
				statuses = append(statuses, DependencyStatus{Name: "Browserless", Status: "Healthy", Details: details, Latency: time.Since(start).String()})
			} else {
				statuses = append(statuses, DependencyStatus{Name: "Browserless", Status: "Unhealthy", Details: details, Error: resp.Status})
			}
		}
	}

	// 4. LLM Check
	llmBase := envClient.GetString("LLM_API_BASE")
	llmKey := envClient.GetString("LLM_API_KEY")
	if llmBase == "" {
		statuses = append(statuses, DependencyStatus{Name: "LLM Service", Status: "Not Configured", Error: "FC_LLM_API_BASE not set"})
	} else {
		start := time.Now()
		details := fmt.Sprintf("Base: %s, Key: %s", llmBase, maskString(llmKey))
		config := openai.DefaultConfig(llmKey)
		config.BaseURL = llmBase
		client := openai.NewClientWithConfig(config)
		_, err := client.ListModels(context.Background())
		if err != nil {
			statuses = append(statuses, DependencyStatus{Name: "LLM Service", Status: "Unhealthy", Details: details, Error: err.Error()})
		} else {
			statuses = append(statuses, DependencyStatus{Name: "LLM Service", Status: "Healthy", Details: details, Latency: time.Since(start).String()})
		}
	}

	c.JSON(http.StatusOK, util.APIResponse[[]DependencyStatus]{
		StatusCode: 0,
		Msg:        "Success",
		Data:       statuses,
	})
}
