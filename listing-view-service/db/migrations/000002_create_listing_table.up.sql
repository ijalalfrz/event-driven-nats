-- denormalize user detail to avoid join
CREATE TABLE IF NOT EXISTS listings (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    listing_type VARCHAR NOT NULL,
    price BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    user_detail JSONB NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);

-- Create indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_listings_user_id ON listings(user_id);
CREATE INDEX IF NOT EXISTS idx_listings_listing_type ON listings(listing_type);
