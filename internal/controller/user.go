package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// UserInfo 对外通信使用的结构
type UserInfo struct {
	NickName    string `json:"nickname"`
	Email       string `json:"email"`
	Username    string `json:"username" binding:"required"`
	Md5Password string `json:"md5Password" binding:"required"` // 前端传递过来的只有md5哈希过的密码
}

func CreateUser(c *gin.Context) {
	var user dao.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := dao.CreateUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	// 不返回密码哈希
	userInfo := UserInfo{
		Username: user.Username,
		NickName: user.NickName,
		Email:    user.Email,
	}
	c.JSON(http.StatusCreated, util.APIResponse[any]{Data: userInfo})
}

func GetUser(c *gin.Context) {
	username := c.Param("username")
	db := util.GetDatabase()

	user, err := dao.GetUserByUsername(db, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	userInfo := UserInfo{
		Username: user.Username,
		NickName: user.NickName,
		Email:    user.Email,
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: userInfo})
}

func UpdateUser(c *gin.Context) {
	username := c.Param("username")
	var user dao.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	existingUser, err := dao.GetUserByUsername(db, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	// 仅更新非密码字段
	existingUser.NickName = user.NickName
	existingUser.Email = user.Email
	if user.Password != "" {
		// No need to set the password field here
	}

	if err := dao.UpdateUser(db, existingUser); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	// 不返回密码哈希
	userInfo := UserInfo{
		Username: existingUser.Username,
		NickName: existingUser.NickName,
		Email:    existingUser.Email,
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: userInfo})
}

func DeleteUser(c *gin.Context) {
	username := c.Param("username")
	db := util.GetDatabase()

	if err := dao.DeleteUser(db, username); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{})
}

func ListUsers(c *gin.Context) {
	db := util.GetDatabase()
	userList, err := dao.ListUsers(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
	}
	var userInfoList []UserInfo
	for _, user := range userList {
		userInfo := UserInfo{
			Username: user.Username,
			NickName: user.NickName,
			Email:    user.Email,
		}
		userInfoList = append(userInfoList, userInfo)
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: userInfoList})
}
