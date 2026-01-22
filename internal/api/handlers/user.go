package handlers

import (
	"github.com/gin-gonic/gin"
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

}
