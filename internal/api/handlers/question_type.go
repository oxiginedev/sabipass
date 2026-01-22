package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oxiginedev/sabipass/internal/models"
)

type questionTypeHandler struct {
	questionTypeRepo models.QuestionTypeRepository
}

func NewQuestionTypeHandler(questionTypeRepo models.QuestionTypeRepository) *questionTypeHandler {
	return &questionTypeHandler{questionTypeRepo: questionTypeRepo}
}

func (q *questionTypeHandler) HandleGetAllQuestionTypes(c *gin.Context) {
	questionTypes, err := q.questionTypeRepo.FindAll(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			models.NewErrorResponse("failed to retrieve question types", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse("question types retrieved", questionTypes))
}
