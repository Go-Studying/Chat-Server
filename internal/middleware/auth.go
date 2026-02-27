package middleware

import (
	"chat-server/internal/config"
	"chat-server/internal/tools/security"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const AuthCookieName = "auth_token"
const contextKey = "userID"

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Request.Cookie(AuthCookieName)
		tokenString := tokenCookie.Value
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, err := security.ParseJWT(tokenString, cfg.JWTSecret)
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
		c.Set(contextKey, userID)

		// Go Context에도 저장 (서비스용)
		ctx := SetUserToContext(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func CurrentUserID(c *gin.Context) (uuid.UUID, error) {
	userID, ok := c.Get(contextKey)
	if !ok {
		return uuid.Nil, errors.New("userID not found")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return uuid.Nil, errors.New("userID in context is not a string")
	}
	return uuid.Parse(userIDStr)
}
