package controller

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
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

	// 1. Database
	dbStatus := DependencyStatus{Name: "Database (SQLite)"}
	dbPath := envClient.GetString("DB_SQLITE_PATH")
	if dbPath != "" {
		dbStatus.Configured = true
		db := util.GetDatabase()
		sqlDB, err := db.DB()
		if err != nil {
			dbStatus.Message = fmt.Sprintf("Failed to get sql.DB: %v", err)
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
		// Note: GetRedisClient calls log.Fatal if parse fails, so we trust it parses if env is set,
		// but to be safe and avoid crashing, we might want to manually parse or try/catch if it was possible in Go.
		// Given the util implementation, we assume if URI is present it's valid enough to try connecting.
		// Wait, GetRedisClient does log.Fatalf if parse fails. This is risky.
		// But I cannot change util in this step. I'll rely on env being correct or app would have crashed on startup anyway.
		rdb := util.GetRedisClient()
		if rdb != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if _, err := rdb.Ping(ctx).Result(); err != nil {
				redisStatus.Message = fmt.Sprintf("Ping failed: %v", err)
			} else {
				redisStatus.Healthy = true
				redisStatus.Message = "Connected"
			}
		}
	} else {
		redisStatus.Message = "FC_REDIS_URI not set"
	}
	statuses = append(statuses, redisStatus)

	// 3. LLM Service
	llmStatus := DependencyStatus{Name: "LLM Service"}
	// Check for various keys
	llmType := envClient.GetString("LLM_API_TYPE")
	llmKey := envClient.GetString("LLM_API_KEY")
	// Backwards compatibility
	if llmKey == "" {
		llmKey = envClient.GetString("OPENAI_AUTH_KEY")
	}

	if llmType == "ollama" {
		llmStatus.Configured = true // Ollama might not need key
	} else if llmKey != "" {
		llmStatus.Configured = true
	}

	if llmStatus.Configured {
		// We don't make a live call to save tokens/time, but we could add a specific test button later.
		// For now, assume healthy if configured.
		// Or maybe check if BaseURL is reachable?
		llmStatus.Healthy = true // Placeholder for now
		llmStatus.Message = "Configured (Health check skipped)"
	} else {
		llmStatus.Message = "FC_LLM_API_KEY / FC_OPENAI_AUTH_KEY not set"
	}
	statuses = append(statuses, llmStatus)

	// 4. Browserless / Puppeteer
	pupStatus := DependencyStatus{Name: "Browserless (Puppeteer)"}
	pupEndpoint := envClient.GetString("PUPPETEER_HTTP_ENDPOINT")
	if pupEndpoint != "" {
		pupStatus.Configured = true
		// Try to connect to the endpoint
		client := resty.New()
		client.SetTimeout(2 * time.Second)
		// Usually /version is a safe endpoint for browserless, or just GET /
		// The craft code uses it for /content.
		resp, err := client.R().Get(pupEndpoint)
		if err != nil {
			pupStatus.Message = fmt.Sprintf("Connection failed: %v", err)
		} else {
			// even 404 means it's reachable
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
