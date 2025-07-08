#!/bin/sh

# Get APP_ENV from environment, default to development if not set
APP_ENV=${APP_ENV:-development}

# Set DB_DSN based on APP_ENV
if [ "$APP_ENV" = "development" ]; then
    export DB_DSN="postgres://postgres@postgres-listing-service/listing_development?sslmode=disable"
else
    export DB_DSN="postgres://postgres@postgres-listing-service/listing_test?sslmode=disable"
fi

if [ "$1" = "create" ]; then
    migrate create -ext sql -dir ./db/migrations -seq "$2"
else
    migrate -path ./db/migrations -database "$DB_DSN" "$@"
fi