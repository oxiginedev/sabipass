package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oxiginedev/sabipass/config"
	"github.com/oxiginedev/sabipass/internal/api/handlers"
	"github.com/oxiginedev/sabipass/internal/api/middleware"
	"github.com/oxiginedev/sabipass/internal/models"
	"github.com/oxiginedev/sabipass/internal/pkg/jwt"
)

type API struct {
	cfg              *config.Config
	tokenManager     jwt.TokenManager
	userRepo         models.UserRepository
	quizRepo         models.QuizRepository
	questionTypeRepo models.QuestionTypeRepository
}

func NewAPI(cfg *config.Config,
	tokenManager jwt.TokenManager,
	userRepo models.UserRepository,
	quizRepo models.QuizRepository,
	questionTypeRepo models.QuestionTypeRepository,
) *API {
	return &API{
		cfg:              cfg,
		tokenManager:     tokenManager,
		userRepo:         userRepo,
		quizRepo:         quizRepo,
		questionTypeRepo: questionTypeRepo,
	}
}

func (a *API) RegisterRoutes() http.Handler {
	router := gin.New()
	if a.cfg.Environment == config.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	oauthHandler := handlers.NewOauthHandler(a.cfg, a.tokenManager, a.userRepo)
	userHandler := handlers.NewUserHandler(a.userRepo)
	quizHandler := handlers.NewQuizHandler(a.quizRepo)
	questionTypeHandler := handlers.NewQuestionTypeHandler(a.questionTypeRepo)

	router.Use(gin.Recovery())
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, models.NewErrorResponse("the requested route was not found", nil))
	})

	router.GET("/oauth/google/redirect", oauthHandler.HandleGoogleLoginRedirect)
	router.GET("/oauth/google/callback", oauthHandler.HandleGoogleLoginCallback)

	authRouter := router.Group("/", middleware.RequireAuth(a.tokenManager, a.userRepo))
	{
		authRouter.GET("/users/me", userHandler.HandleGetCurrentUser)

		authRouter.GET("/quizzes", quizHandler.HandleGetAllQuizzes)
		authRouter.POST("/quizzes", quizHandler.HandleCreateQuiz)
		authRouter.GET("/quizzes/:quizid", quizHandler.HandleGetQuiz)
		authRouter.PATCH("/quizzes/:quizid", quizHandler.HandleEditQuiz)

		authRouter.GET("/question-types", questionTypeHandler.HandleGetAllQuestionTypes)
	}

	return router
}
