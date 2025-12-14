package controller

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
)

type DependencyStatus struct {
	Name       string `json:"name"`
	Configured bool   `json:"configured"`
	Healthy    bool   `json:"healthy"`
	Message    string `json:"message"`
}

func GetDependencyStatus(c *gin.Context) {
	statuses := []DependencyStatus{}
	envClient := util.GetEnvClient()

	// 1. Database (SQLite)
	dbStatus := DependencyStatus{Name: "Database (SQLite)"}
	dbPath := envClient.GetString("DB_SQLITE_PATH")
	if dbPath != "" {
		dbStatus.Configured = true
		// Safe check: relying on app startup for DB init, but avoiding util fatal call if possible.
		// Since util.GetDatabase is singleton with fatal, we trust app state if running.
		db := util.GetDatabase()
		sqlDB, err := db.DB()
		if err != nil {
			dbStatus.Message = fmt.Sprintf("Failed to get underlying DB: %v", err)
		} else {
			if err := sqlDB.Ping(); err != nil {
				dbStatus.Message = fmt.Sprintf("Ping failed: %v", err)
			} else {
				dbStatus.Healthy = true
				dbStatus.Message = "Connected"
			}
		}
	} else {
		dbStatus.Message = "FC_DB_SQLITE_PATH not set"
	}
	statuses = append(statuses, dbStatus)

	// 2. Redis
	redisStatus := DependencyStatus{Name: "Redis"}
	redisURI := envClient.GetString("REDIS_URI")
	if redisURI != "" {
		redisStatus.Configured = true
		// Manual connection to avoid log.Fatal in util.GetRedisClient
		opts, err := redis.ParseURL(redisURI)
		if err != nil {
			redisStatus.Message = fmt.Sprintf("Invalid Redis URI: %v", err)
		} else {
			rdb := redis.NewClient(opts)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if _, err := rdb.Ping(ctx).Result(); err != nil {
				redisStatus.Message = fmt.Sprintf("Ping failed: %v", err)
			} else {
				redisStatus.Healthy = true
				redisStatus.Message = "Connected"
				rdb.Close()
			}
		}
	} else {
		redisStatus.Message = "FC_REDIS_URI not set"
	}
	statuses = append(statuses, redisStatus)

	// 3. LLM Service
	llmStatus := DependencyStatus{Name: "LLM Service"}
	llmType := envClient.GetString("LLM_API_TYPE")
	if llmType == "" {
		llmType = "openai"
	}
	llmBase := envClient.GetString("LLM_API_BASE")
	llmKey := envClient.GetString("LLM_API_KEY")
	llmModel := envClient.GetString("LLM_API_MODEL")

	// Legacy fallback
	if llmBase == "" {
		llmBase = envClient.GetString("OPENAI_ENDPOINT")
	}
	if llmKey == "" {
		llmKey = envClient.GetString("OPENAI_AUTH_KEY")
	}
	if llmModel == "" {
		llmModel = envClient.GetString("OPENAI_DEFAULT_MODEL")
	}

	if llmType == "ollama" {
		if llmBase != "" {
			llmStatus.Configured = true
		} else {
			llmStatus.Message = "FC_LLM_API_BASE required for Ollama"
		}
	} else {
		if llmKey != "" {
			llmStatus.Configured = true
		} else {
			llmStatus.Message = "FC_LLM_API_KEY not set"
		}
	}

	if llmStatus.Configured {
		conf := openai.DefaultConfig(llmKey)
		if llmBase != "" {
			conf.BaseURL = llmBase
		}
		// In Ollama mode, key might be empty, which is fine for local instances,
		// but go-openai might strictly require a non-empty string for some calls.
		// We set a dummy key if empty to satisfy client init, as Ollama ignores it.
		if llmType == "ollama" && llmKey == "" {
			// Do not attempt to set APIToken field directly as it doesn't exist on Config struct in this version
			// Instead, create a config with a dummy key initially if needed,
			// or just let DefaultConfig handle it (which requires a key).
			conf = openai.DefaultConfig("ollama")
			conf.BaseURL = llmBase
		}

		client := openai.NewClientWithConfig(conf)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		models, err := client.ListModels(ctx)
		if err != nil {
			llmStatus.Message = fmt.Sprintf("Health check failed: %v", err)
		} else {
			if len(models.Models) > 0 {
				llmStatus.Healthy = true
				llmStatus.Message = fmt.Sprintf("Connected (Available models: %d)", len(models.Models))
			} else {
				llmStatus.Healthy = true
				llmStatus.Message = "Connected (No models found)"
			}
		}
	}
	statuses = append(statuses, llmStatus)

	// 4. Browserless / Puppeteer
	pupStatus := DependencyStatus{Name: "Browserless (Puppeteer)"}
	pupEndpoint := envClient.GetString("PUPPETEER_HTTP_ENDPOINT")
	if pupEndpoint != "" {
		pupStatus.Configured = true
		client := resty.New()
		client.SetTimeout(2 * time.Second)
		resp, err := client.R().Get(pupEndpoint)
		if err != nil {
			pupStatus.Message = fmt.Sprintf("Connection failed: %v", err)
		} else {
			if resp.StatusCode() >= 200 && resp.StatusCode() < 500 {
				pupStatus.Healthy = true
				pupStatus.Message = "Connected"
			} else {
				pupStatus.Message = fmt.Sprintf("Unhealthy status: %d", resp.StatusCode())
			}
		}
	} else {
		pupStatus.Message = "FC_PUPPETEER_HTTP_ENDPOINT not set"
	}
	statuses = append(statuses, pupStatus)

	c.JSON(http.StatusOK, util.APIResponse[[]DependencyStatus]{
		Data: statuses,
	})
}
