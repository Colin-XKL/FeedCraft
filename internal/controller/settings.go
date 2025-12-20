package controller

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSearchProviderConfig(c *gin.Context) {
	db := util.GetDatabase()
	var providerConfig config.SearchProviderConfig
	if err := dao.GetJsonValue(db, constant.KeySearchProviderConfig, &providerConfig); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	providerConfig.Mask()
	c.JSON(http.StatusOK, providerConfig)
}

func SaveSearchProviderConfig(c *gin.Context) {
	var input config.SearchProviderConfig
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	db := util.GetDatabase()
	var currentConfig config.SearchProviderConfig
	// Ignore error if not found, it will be empty
	_ = dao.GetJsonValue(db, constant.KeySearchProviderConfig, &currentConfig)

	currentConfig.Merge(input)

	if err := dao.SetJsonValue(db, constant.KeySearchProviderConfig, currentConfig); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	// Mask before returning
	currentConfig.Mask()
	c.JSON(http.StatusOK, currentConfig)
}
