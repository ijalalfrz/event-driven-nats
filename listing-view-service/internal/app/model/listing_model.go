package model

// denormalized listing model with user as json to minimize join
type Listing struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	ListingType string `json:"listing_type"`
	Price       int64  `json:"price"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	User        User
}
