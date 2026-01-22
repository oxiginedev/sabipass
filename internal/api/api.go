package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oxiginedev/sabipass/config"
	"github.com/oxiginedev/sabipass/internal/api/handlers"
	"github.com/oxiginedev/sabipass/internal/models"
)

type API struct {
	cfg      *config.Config
	userRepo models.UserRepository
}

func NewAPI(cfg *config.Config, userRepo models.UserRepository) *API {
	return &API{cfg: cfg, userRepo: userRepo}
}

func (a *API) RegisterRoutes() http.Handler {
	router := gin.New()
	if a.cfg.Environment == config.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	oauthHandler := handlers.NewOauthHandler(a.cfg, a.userRepo)

	router.Use(gin.Recovery())
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	router.GET("/oauth/google/redirect", oauthHandler.HandleGoogleLoginRedirect)
	router.GET("/oauth/google/callback", oauthHandler.HandleGoogleLoginCallback)

	return router
}
