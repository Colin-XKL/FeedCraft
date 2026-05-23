package middleware

import (
	"errors"
	"net/http"
	"strings"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SystemAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Authorization format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		db := util.GetDatabase()

		_, err := dao.GetSystemAuthTokenByToken(db, tokenString)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid token"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Database error"})
			c.Abort()
			return
		}

		c.Next()
	}
}
