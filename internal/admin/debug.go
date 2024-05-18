package admin

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
For debug use
*/

type llmDebugReq struct {
	Model string `json:"model" binding:"required"`
	Input string `json:"input" binding:"required"`
}

type llmDebugResp struct {
	Output string `json:"output"`
}

func LLMDebug(c *gin.Context) {
	reqBody := &llmDebugReq{}
	err := c.ShouldBindJSON(reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	ret, err := adapter.SimpleLLMCall(reqBody.Model, reqBody.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	resp := llmDebugResp{Output: ret}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: resp})
}
