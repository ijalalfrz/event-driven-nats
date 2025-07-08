package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/model"
)

type ListingRepository struct {
	db *sql.DB
	errorMapper
	transactable
}

func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{
		db: db,
		transactable: transactable{
			db: db,
		},
	}
}

func (r *ListingRepository) GetAll(ctx context.Context, limit,
	offset int, userID *int64) ([]model.Listing, error) {
	var args []interface{}
	var paramIndex int = 1
	var user model.User
	var userBytes []byte

	query := `
		SELECT id, user_id, listing_type, price, user_detail,created_at, updated_at
		FROM listings
	`

	if userID != nil {
		query += ` WHERE user_id = $` + fmt.Sprintf("%d", paramIndex)
		args = append(args, userID)
		paramIndex++
	}

	query += fmt.Sprintf(`
			ORDER BY created_at DESC
			LIMIT $%d
			OFFSET $%d`, paramIndex, paramIndex+1)

	args = append(args, limit, offset)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, r.errorMapper.mapError(err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, r.errorMapper.mapError(err)
	}

	defer rows.Close()

	listings := []model.Listing{}
	for rows.Next() {
		var listing model.Listing

		err := rows.Scan(&listing.ID, &listing.UserID, &listing.ListingType,
			&listing.Price, &userBytes, &listing.CreatedAt, &listing.UpdatedAt)
		if err != nil {
			return nil, r.errorMapper.mapError(err)
		}

		err = json.Unmarshal(userBytes, &user)
		if err != nil {
			return nil, r.errorMapper.mapError(err)
		}

		listing.User = user

		listings = append(listings, listing)
	}

	return listings, nil
}

func (r *ListingRepository) CreateTx(ctx context.Context, tx *sql.Tx, listing *model.Listing) error {
	query := `
		INSERT INTO listings (id, user_id, listing_type, price, user_detail, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	userDetail, err := json.Marshal(listing.User)
	if err != nil {
		return fmt.Errorf("marshalling user detail: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return r.errorMapper.mapError(err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, listing.ID, listing.UserID,
		listing.ListingType, listing.Price, userDetail, listing.CreatedAt, listing.UpdatedAt)
	if err != nil {
		return r.errorMapper.mapError(err)
	}

	return nil
}
