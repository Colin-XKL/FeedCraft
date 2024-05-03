package admin

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !loginValidator(input.Username, input.Password) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid username or password"})
		return
	}
	token, err := GenerateToken(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{"message": "validated!", "token": token, "username": input.Username})
}
