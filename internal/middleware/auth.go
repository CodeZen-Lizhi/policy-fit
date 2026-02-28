package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhenglizhi/policy-fit/pkg/response"
)

const userIDKey = "user_id"

// Auth JWT 鉴权中间件
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithStatus(c, http.StatusUnauthorized, "PFIT-1002", "请重新登录")
			c.Abort()
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			response.ErrorWithStatus(c, http.StatusUnauthorized, "PFIT-1002", "请重新登录")
			c.Abort()
			return
		}

		tokenRaw := strings.TrimPrefix(authHeader, prefix)
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenRaw, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			response.ErrorWithStatus(c, http.StatusUnauthorized, "PFIT-1002", "请重新登录")
			c.Abort()
			return
		}

		userIDValue, ok := claims["user_id"]
		if !ok {
			response.ErrorWithStatus(c, http.StatusUnauthorized, "PFIT-1002", "请重新登录")
			c.Abort()
			return
		}

		switch id := userIDValue.(type) {
		case float64:
			c.Set(userIDKey, int64(id))
		case int64:
			c.Set(userIDKey, id)
		default:
			response.ErrorWithStatus(c, http.StatusUnauthorized, "PFIT-1002", "请重新登录")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID 获取当前用户 ID
func GetUserID(c *gin.Context) int64 {
	if value, ok := c.Get(userIDKey); ok {
		if userID, ok := value.(int64); ok {
			return userID
		}
	}
	return 0
}
