package middleware

import (
	"net/http"
	"strings"

	"github.com/ardipermana59/go-template/internal/auth"
	"github.com/ardipermana59/go-template/internal/common/apperror"
	"github.com/ardipermana59/go-template/internal/common/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Unauthorized",
				apperror.NewErrors(apperror.NewError("authorization", "Authorization header is required")))
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Unauthorized",
				apperror.NewErrors(apperror.NewError("authorization", "Invalid authorization format. Use: Bearer <token>")))
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenParts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Unauthorized",
				apperror.NewErrors(apperror.NewError("token", "Invalid or expired token")))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Unauthorized",
				apperror.NewErrors(apperror.NewError("role", "User role not found in token")))
			c.Abort()
			return
		}

		role := userRole.(string)
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Error(c, http.StatusForbidden, "Forbidden",
				apperror.NewErrors(apperror.NewError("permission", "You don't have permission to access this resource")))
			c.Abort()
			return
		}

		c.Next()
	}
}
