package controller

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/favicon"
	"FeedCraft/internal/source/fetcher/provider"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchProviderConfigResponse struct {
	config.SearchProviderConfig
	HasAPIKey bool `json:"has_api_key"`
}

type SearchProviderConfigRequest struct {
	config.SearchProviderConfig
	UpdateAPIKey bool `json:"update_api_key"`
}

type FaviconProviderConfigResponse struct {
	config.FaviconSettings
	Providers []favicon.ProviderDescriptor `json:"providers"`
}

type FaviconProviderPreviewRequest struct {
	Settings   config.FaviconSettings `json:"settings"`
	ProviderID string                 `json:"provider_id"`
	PageURL    string                 `json:"page_url" binding:"required"`
	Size       int                    `json:"size"`
}

type FaviconProviderPreviewResponse struct {
	URL        string `json:"url"`
	ProviderID string `json:"provider_id"`
}

func GetSearchProviderConfig(c *gin.Context) {
	db := util.GetDatabase()

	var cfg config.SearchProviderConfig
	err := dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	resp := SearchProviderConfigResponse{
		SearchProviderConfig: cfg,
		HasAPIKey:            cfg.APIKey != "",
	}
	resp.APIKey = ""

	c.JSON(http.StatusOK, util.APIResponse[SearchProviderConfigResponse]{Data: resp})
}

func SaveSearchProviderConfig(c *gin.Context) {
	var req SearchProviderConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	db := util.GetDatabase()

	// Fetch existing config to handle empty APIKey
	var existingCfg config.SearchProviderConfig
	if err := dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &existingCfg); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if !req.UpdateAPIKey && req.APIKey == "" {
		req.APIKey = existingCfg.APIKey
	}

	if err := dao.SetJsonSetting(db, constant.KeySearchProviderConfig, req.SearchProviderConfig); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Msg: "success"})
}

func CheckSearchProviderConfig(c *gin.Context) {
	var req SearchProviderConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if !req.UpdateAPIKey && req.APIKey == "" {
		db := util.GetDatabase()
		var existingCfg config.SearchProviderConfig
		if err := dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &existingCfg); err == nil {
			req.APIKey = existingCfg.APIKey
		}
	}

	prv, err := provider.Get(req.Provider, &req.SearchProviderConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to create provider: " + err.Error()})
		return
	}

	_, err = prv.Fetch(c.Request.Context(), "FeedCraft")
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Connection check failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Msg: "success"})
}

func GetFaviconProviderConfig(c *gin.Context) {
	settings := favicon.Settings()
	c.JSON(http.StatusOK, util.APIResponse[FaviconProviderConfigResponse]{
		Data: FaviconProviderConfigResponse{
			FaviconSettings: settings,
			Providers:       favicon.Providers(),
		},
	})
}

func SaveFaviconProviderConfig(c *gin.Context) {
	var settings config.FaviconSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if err := favicon.Save(util.GetDatabase(), settings); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[FaviconProviderConfigResponse]{
		Data: FaviconProviderConfigResponse{
			FaviconSettings: favicon.Settings(),
			Providers:       favicon.Providers(),
		},
		Msg: "success",
	})
}

func PreviewFaviconProviderConfig(c *gin.Context) {
	var req FaviconProviderPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	if req.Size == 0 {
		req.Size = 64
	}
	if req.Size < 16 || req.Size > 256 {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "size must be between 16 and 256"})
		return
	}

	iconURL, providerID, err := favicon.BuildURLFromSettings(req.Settings, req.ProviderID, req.PageURL, req.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	if iconURL == "" {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "page_url must be an absolute HTTP or HTTPS URL"})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[FaviconProviderPreviewResponse]{
		Data: FaviconProviderPreviewResponse{
			URL:        iconURL,
			ProviderID: providerID,
		},
	})
}
