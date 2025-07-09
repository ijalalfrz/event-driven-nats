package repository

import (
	"context"
	"database/sql"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/model"
)

type UserRepository struct {
	db *sql.DB
	errorMapper
	transactable
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
		transactable: transactable{
			db: db,
		},
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (model.User, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return model.User{}, r.errorMapper.mapError(err)
	}

	defer stmt.Close()

	var user model.User
	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, r.errorMapper.mapError(err)
	}

	return user, nil
}

func (r *UserRepository) CreateTx(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
		INSERT INTO users (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			name = $2,
			updated_at = $4
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return r.errorMapper.mapError(err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.ID, user.Name, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return r.errorMapper.mapError(err)
	}

	return nil
}
