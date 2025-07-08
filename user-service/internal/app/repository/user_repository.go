package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/model"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/exception"
)

type UserRepository struct {
	db *sql.DB
	errorMapper
	transactable
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:           db,
		transactable: transactable{db: db},
	}
}

func (r *UserRepository) GetAll(ctx context.Context, limit, offset int) ([]model.User, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, r.errorMapper.mapError(err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, limit, offset)
	if err != nil {
		return nil, r.errorMapper.mapError(err)
	}

	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, r.errorMapper.mapError(err)
		}
		users = append(users, user)
	}

	return users, nil
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
	if errors.Is(err, sql.ErrNoRows) {
		err := exception.ErrRecordNotFound
		err.MessageVars = map[string]interface{}{
			"name": "user",
		}

		return model.User{}, err
	}

	if err != nil {
		return model.User{}, r.errorMapper.mapError(err)
	}

	return user, nil
}

func (r *UserRepository) CreateTx(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
		INSERT INTO users (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return r.errorMapper.mapError(err)
	}

	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, user.Name, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		return r.errorMapper.mapError(err)
	}

	return nil
}
