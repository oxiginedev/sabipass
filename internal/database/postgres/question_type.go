package postgres

import (
	"context"

	"github.com/oxiginedev/sabipass/internal/models"
)

type questionTypeRepo struct {
	db *DB
}

func NewQuestionTypeRepository(db *DB) models.QuestionTypeRepository {
	return &questionTypeRepo{db: db}
}

func (q *questionTypeRepo) FindAll(ctx context.Context) ([]models.QuestionType, error) {
	ctx, cancel := q.db.WithContext(ctx)
	defer cancel()

	var questionTypes []models.QuestionType
	query := q.db.NewSelect().Model(&questionTypes)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return questionTypes, nil
}
