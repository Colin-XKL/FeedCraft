package admin

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
	"time"
)

const tokenLifeSpanHour = 1

var (
	randomSignKey     string
	randomSignKeyOnce sync.Once
)

// getJwtSignKey returns API_SECRET. When it is unset, a random key is
// generated once per process instead of a public default, so tokens can't be
// forged; they stop being valid when the server restarts.
func getJwtSignKey() string {
	if secretKey := os.Getenv("API_SECRET"); secretKey != "" {
		return secretKey
	}
	randomSignKeyOnce.Do(func() {
		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err != nil {
			panic(fmt.Sprintf("generate JWT sign key: %v", err))
		}
		randomSignKey = hex.EncodeToString(buf)
		logrus.Warn("API_SECRET is not set; using a random JWT sign key. Admin sessions will end when the server restarts. Set API_SECRET to keep them.")
	})
	return randomSignKey
}

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user"] = username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifeSpanHour)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getJwtSignKey()))

}
func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getJwtSignKey()), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// ExtractToken extract token from header
func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return bearerToken
}

func ExtractUserFromToken(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getJwtSignKey()), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims["user"].(string), nil
	}
	return "", nil
}
