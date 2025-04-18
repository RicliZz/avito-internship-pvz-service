package middleware

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func CheckRoleMiddleware(secret string, expectedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, models.Error{Message: "Invalid Header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, models.Error{Message: "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, models.Error{Message: "Invalid token"})
			return
		}
		role, ok := claims["role"].(string)
		if !ok || role != expectedRole {
			c.AbortWithStatusJSON(403, models.Error{Message: "Forbidden"})
			return
		}

		c.Next()
	}
}
