package admin

import (
	"FeedCraft/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func loginValidator(username, passwd string) bool {
	// TODO replace the mock val
	return username == "admin" && passwd == "adminadmin"
}

type UserAuth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginAuth(c *gin.Context) {
	var input UserAuth

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[string]{Msg: err.Error()})
		return
	}
	if !loginValidator(input.Username, input.Password) {
		c.JSON(http.StatusForbidden, util.APIResponse[string]{Msg: "invalid username or password"})
		return
	}
	token, err := GenerateToken(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[string]{Msg: err.Error()})
		c.Abort()
	}
	ret := map[string]string{
		"token": token,
	}
	c.JSON(http.StatusOK, util.APIResponse[map[string]string]{Data: ret})
}
