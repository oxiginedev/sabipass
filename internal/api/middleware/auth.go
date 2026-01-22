package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oxiginedev/sabipass/internal/models"
	"github.com/oxiginedev/sabipass/internal/pkg/jwt"
	"github.com/oxiginedev/sidekik"
)

func RequireAuth(tokenManager jwt.TokenManager, userRepo models.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if sidekik.IsStringEmpty(authHeader) {
			slog.Info("[middleware]: empty authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("unauthenticated", nil))
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 && strings.ToLower(parts[0]) != "bearer" {
			slog.Error("[middleware]: malformed authorization header", slog.String("authHeader", authHeader))
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("unauthenticated", nil))
			return
		}

		tokenString := parts[1]

		validatedToken, err := tokenManager.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(419, models.NewErrorResponse("session expired", nil))
				return
			}

			slog.Error("[middleware]: invalid token", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("unauthenticated", nil))
			return
		}

		user, err := userRepo.FindOne(c.Request.Context(), &models.FindUserOptions{
			ID: validatedToken.UserID,
		})
		if err != nil {
			slog.Error("[middleware]: could not find user", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("unauthenticated", nil))
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
