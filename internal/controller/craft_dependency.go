package controller

import (
	"FeedCraft/internal/craft"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AnalyzeCraftDependencies(c *gin.Context) {
	roots, err := craft.AnalyzeDependencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to analyze dependencies: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: roots})
}
