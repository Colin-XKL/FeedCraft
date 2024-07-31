package admin

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func loginValidator(username, md5Password string, db *gorm.DB) bool {
	user, err := dao.GetUserByUsername(db, username)
	if err != nil {
		return false
	}
	hashedPassword := dao.HashPassword(md5Password, user.Salt)
	return hashedPassword == user.PasswordHash
}

type UserAuth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginAuth(c *gin.Context) {
	var input UserAuth
	db := util.GetDatabase()

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[string]{Msg: err.Error()})
		return
	}
	md5Password := md5.Sum([]byte(input.Password))
	if !loginValidator(input.Username, hex.EncodeToString(md5Password[:]), db) {
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
