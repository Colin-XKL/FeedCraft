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

	var cfg config.SearchProviderConfig
	err := dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[config.SearchProviderConfig]{Data: cfg})
}

func SaveSearchProviderConfig(c *gin.Context) {
	var cfg config.SearchProviderConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	db := util.GetDatabase()

	if err := dao.SetJsonSetting(db, constant.KeySearchProviderConfig, cfg); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Msg: "success"})
}
