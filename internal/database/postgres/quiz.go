package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/oxiginedev/sabipass/internal/database"
	"github.com/oxiginedev/sabipass/internal/models"
	"github.com/oxiginedev/sidekik"
	"github.com/uptrace/bun"
)

type quizRepo struct {
	db *DB
}

func NewQuizRepository(db *DB) models.QuizRepository {
	return &quizRepo{db: db}
}

func (q *quizRepo) Create(ctx context.Context, quiz *models.Quiz) error {
	ctx, cancel := q.db.WithContext(ctx)
	defer cancel()

	return q.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(quiz).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (q *quizRepo) Update(ctx context.Context, quiz *models.Quiz) error {
	ctx, cancel := q.db.WithContext(ctx)
	defer cancel()

	return q.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewUpdate().
			Model(quiz).
			Where("id = ?", quiz.ID).
			Exec(ctx)
		return err
	})
}

func (q *quizRepo) FindOne(ctx context.Context, opts *models.FindQuizOptions) (*models.Quiz, error) {
	ctx, cancel := q.db.WithContext(ctx)
	defer cancel()

	var quiz = models.Quiz{
		Questions: []models.Question{},
	}
	query := q.db.NewSelect().Model(&quiz)

	if !sidekik.IsStringEmpty(opts.ID) {
		query.Where("id = ?", opts.ID)
	}

	if !sidekik.IsStringEmpty(opts.OwnerID) {
		query.Where("owner_id = ?", opts.OwnerID)
	}

	if err := query.
		Relation("Questions").
		Relation("Questions.QuestionOptions").
		Relation("Questions.QuestionType").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = database.ErrQuizNotFound
		}

		return nil, err
	}

	return &quiz, nil
}

func (q *quizRepo) FindAll(ctx context.Context, opts *models.ListQuizOptions) ([]models.Quiz, int64, error) {
	ctx, cancel := q.db.WithContext(ctx)
	defer cancel()

	var quizzes []models.Quiz
	query := q.db.NewSelect().Model(&quizzes)

	if !sidekik.IsStringEmpty(opts.OwnerID) {
		query.Where("owner_id = ?", opts.OwnerID)
	}

	if !sidekik.IsStringEmpty(opts.Search) {
		query.Where("title ILIKE ?", "%"+opts.Search+"%")
	}

	if opts.Visibility.IsValid() {
		query.Where("visibility = ?", opts.Visibility.String())
	}

	quizCount, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Limit(int(opts.Paginator.PerPage)).
		Offset(int(opts.Paginator.Offset())).
		Scan(ctx); err != nil {
		return nil, 0, err
	}

	return quizzes, int64(quizCount), nil
}
