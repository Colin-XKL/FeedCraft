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
	Md5Password string `json:"md5_password"` // 新增字段
}

func CreateUser(c *gin.Context) {
	var input UserInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	user := dao.User{
		Username: input.Username,
		NickName: input.NickName,
		Email:    input.Email,
	}
	if err := dao.CreateUser(db, &user, input.Md5Password); err != nil {
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
	var input UserInfo
	if err := c.ShouldBindJSON(&input); err != nil {
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
	existingUser.NickName = input.NickName
	existingUser.Email = input.Email
	if input.Md5Password != "" {
		if err := dao.UpdateUser(db, existingUser, input.Md5Password); err != nil {
			c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
			return
		}
	} else {
		if err := dao.UpdateUser(db, existingUser, ""); err != nil {
			c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
			return
		}
	}

	if err := dao.UpdateUser(db, existingUser, input.Md5Password); err != nil {
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
func ChangePassword(c *gin.Context) {
	var input struct {
		Username        string `json:"username" binding:"required"`
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	db := util.GetDatabase()
	user, err := dao.GetUserByUsername(db, input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if !loginValidator(input.Username, input.CurrentPassword, db) {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Invalid current password"})
		return
	}

	if len(input.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Password must be at least 6 characters long"})
		return
	}
	if isNumeric(input.NewPassword) {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Password cannot be purely numeric"})
		return
	}
	if err := dao.UpdateUser(db, user, input.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Msg: "Password updated successfully"})
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
