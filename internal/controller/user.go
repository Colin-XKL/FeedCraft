package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

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
	user.PasswordHash = nil
	user.PasswordHash = nil
	c.JSON(http.StatusCreated, util.APIResponse[any]{Data: user})
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

	user.PasswordHash = nil
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: user})
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
		if err := dao.UpdateUser(db, existingUser); err != nil {
			c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
			return
		}
	}

	// 不返回密码哈希
	existingUser.PasswordHash = nil
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: existingUser})
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
	for _, user := range userList {
		user.PasswordHash = nil
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: userList})
}
