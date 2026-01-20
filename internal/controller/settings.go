package controller

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
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
	_ = dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &existingCfg)

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
