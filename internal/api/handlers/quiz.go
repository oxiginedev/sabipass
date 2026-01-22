package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oxiginedev/sabipass/internal/api/middleware"
	"github.com/oxiginedev/sabipass/internal/database"
	"github.com/oxiginedev/sabipass/internal/models"
	"github.com/oxiginedev/sabipass/utils"
)

type quizHandler struct {
	quizRepo models.QuizRepository
}

func NewQuizHandler(quizRepo models.QuizRepository) *quizHandler {
	return &quizHandler{
		quizRepo: quizRepo,
	}
}

func (q *quizHandler) HandleCreateQuiz(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c)
	if !ok {
		slog.Error("[quiz handler]: could not get user from context")
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("unauthorized", nil))
		return
	}

	var req models.CreateOrEditQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("invalid request body", nil))
		return
	}

	err := utils.Validate(req)
	if err != nil {
		verr, _ := err.(*utils.ValidatorErrorBag)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity,
			models.NewErrorResponse(verr.Error(), verr.Errors))
		return
	}

	quiz := &models.Quiz{
		ID:          utils.Uuid(),
		OwnerID:     user.ID,
		Title:       req.Title,
		Description: utils.Ptr(req.Description),
		Visibility:  models.QuizVisibility(req.Visibility),
		CoverImage:  utils.Ptr(req.CoverImage),
	}

	if err := q.quizRepo.Create(c.Request.Context(), quiz); err != nil {
		slog.Error("[quiz handler]: could not create quiz", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("failed to create quiz", nil))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse("quiz created successfully", quiz))
}

func (q *quizHandler) HandleGetQuiz(c *gin.Context) {
	quizid := c.Param("quizid")

	quiz, err := q.quizRepo.FindOne(c.Request.Context(), &models.FindQuizOptions{
		ID: quizid,
	})
	if err != nil {
		if errors.Is(err, database.ErrQuizNotFound) {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("quiz not found", nil))
			return
		}

		slog.Error("[quiz handler]: could not get quiz", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("failed to get quiz", nil))
		return
	}

	user, ok := middleware.GetUserFromContext(c)
	if !ok {
		slog.Error("[quiz handler]: could not get user from context")
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("unauthorized", nil))
		return
	}

	if quiz.OwnerID != user.ID {
		c.JSON(http.StatusNotFound, models.NewErrorResponse("quiz not found", nil))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("quiz retrieved successfully", quiz))
}

func (q *quizHandler) HandleGetAllQuizzes(c *gin.Context) {
	search := c.Query("search")
	visibility := c.Query("visibility")

	user, ok := middleware.GetUserFromContext(c)
	if !ok {
		slog.Error("[quiz handler]: could not get user from context")
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse("unauthorized", nil))
		return
	}

	paginator := models.PaginatorFromContext(c)

	quizzes, totalCount, err := q.quizRepo.FindAll(c.Request.Context(), &models.ListQuizOptions{
		OwnerID:    user.ID,
		Search:     search,
		Visibility: models.QuizVisibility(visibility),
		Paginator:  paginator,
	})
	if err != nil {
		slog.Error("[quiz handler]: could not get quizzes", slog.Any("error", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.NewErrorResponse("failed to get quizzes", nil))
		return
	}

	c.JSON(http.StatusOK,
		models.NewPaginatedResponse("quizzes retrieved successfully", quizzes, totalCount, paginator))
}

func (q *quizHandler) HandleEditQuiz(c *gin.Context) {}
