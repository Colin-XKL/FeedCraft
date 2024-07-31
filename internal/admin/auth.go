// 前后端登录通信和鉴权流程方案
//
// 1. 前端发送登录请求时，将用户输入的密码使用 MD5 进行哈希处理，然后将用户名和 MD5 哈希后的密码发送到后端。
// 2. 后端接收到登录请求后，根据用户名从数据库中获取用户信息，包括密码哈希和盐值。
// 3. 后端使用相同的盐值和前端发送的 MD5 哈希密码再次进行SHA256哈希处理，得到最终的密码哈希。
// 4. 后端将计算得到的最终密码哈希与数据库中存储的密码哈希进行比较，如果一致，则验证通过。
// 5. 验证通过后，后端生成一个 JWT 令牌，并将其返回给前端。
// 6. 前端在后续的请求中，将 JWT 令牌放在请求头中，后端通过验证 JWT 令牌来进行鉴权。

package admin

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
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
	Username    string `json:"username" binding:"required"`
	Md5Password string `json:"md5Password" binding:"required"`
}

func LoginAuth(c *gin.Context) {
	var input UserAuth
	db := util.GetDatabase()

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[string]{Msg: err.Error()})
		return
	}
	if !loginValidator(input.Username, input.Md5Password, db) {
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
