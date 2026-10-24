package admin

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func requestWithToken(token string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/api/admin/user/info", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)
	return c
}

func forgedToken(t *testing.T, key string) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"user":       "admin",
		"exp":        time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte(key))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return token
}

// TestTokenValid_RejectsOldDefaultKeyWithoutSecret is the regression test for
// issue #950: without API_SECRET, a token signed with the old public default
// key "FeedCraft" was accepted.
func TestTokenValid_RejectsOldDefaultKeyWithoutSecret(t *testing.T) {
	t.Setenv("API_SECRET", "")

	if err := TokenValid(requestWithToken(forgedToken(t, "FeedCraft"))); err == nil {
		t.Fatal("token signed with the old default key was accepted")
	}

	token, err := GenerateToken("admin")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	if err := TokenValid(requestWithToken(token)); err != nil {
		t.Fatalf("token from GenerateToken rejected: %v", err)
	}
	if user, err := ExtractUserFromToken(requestWithToken(token)); err != nil || user != "admin" {
		t.Fatalf("ExtractUserFromToken = %q, %v; want admin", user, err)
	}
}

func TestTokenValid_UsesAPISecret(t *testing.T) {
	t.Setenv("API_SECRET", "test-secret")

	if err := TokenValid(requestWithToken(forgedToken(t, "test-secret"))); err != nil {
		t.Fatalf("token signed with API_SECRET rejected: %v", err)
	}
	if err := TokenValid(requestWithToken(forgedToken(t, "FeedCraft"))); err == nil {
		t.Fatal("token signed with the old default key was accepted")
	}
}
