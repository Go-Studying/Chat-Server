package middleware

import (
	"chat-server/internal/tools/security"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, err := security.ParseJWT(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unverified token"})
			return
		}

		userID, err := claims.GetSubject()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("userID", userID)

		// Go Context에도 저장 (서비스용)
		ctx := SetUserToContext(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetCurrentUser(c *gin.Context) (string, error) {
	userID, ok := c.Get("userID")
	if !ok {
		return "", errors.New("userID not found")
	}

	return userID.(string), nil
}

func extractToken(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")

	if len(token) < 7 || !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("invalid authorization header")
	}
	return token[7:], nil
}
