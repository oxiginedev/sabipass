package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/oxiginedev/sabipass/internal/database"
	"github.com/oxiginedev/sabipass/internal/models"
	"github.com/oxiginedev/sidekik"
	"github.com/uptrace/bun"
)

type userRepo struct {
	db *DB
}

func NewUserRepository(db *DB) models.UserRepository {
	return &userRepo{db}
}

func (u *userRepo) Create(ctx context.Context, user *models.User) error {
	ctx, cancel := u.db.WithContext(ctx)
	defer cancel()

	return u.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return database.ErrUserAlreadyExists
			}
			return err
		}
		return nil
	})
}

func (u *userRepo) FindOne(ctx context.Context, opts *models.FindUserOptions) (*models.User, error) {
	ctx, cancel := u.db.WithContext(ctx)
	defer cancel()

	var user models.User
	query := u.db.NewSelect().Model(&user)

	if !sidekik.IsStringEmpty(opts.ID) {
		query = query.Where("id = ?", opts.ID)
	}

	if !sidekik.IsStringEmpty(opts.Email) {
		query = query.Where("email = ?", opts.Email)
	}

	err := query.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
