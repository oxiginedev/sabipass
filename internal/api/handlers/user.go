package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oxiginedev/sabipass/internal/api/middleware"
	"github.com/oxiginedev/sabipass/internal/models"
)

type userHandler struct {
	userRepo models.UserRepository
}

func NewUserHandler(userRepo models.UserRepository) *userHandler {
	return &userHandler{
		userRepo: userRepo,
	}
}

func (u *userHandler) HandleGetCurrentUser(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.NewErrorResponse("unauthenticated", nil))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("user profile retrieved", user))
}
