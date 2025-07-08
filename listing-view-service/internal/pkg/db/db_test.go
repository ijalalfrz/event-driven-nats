//go:build unit

package db

import (
	"testing"

	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	cfg := config.Config{
		DB: config.DB{
			DSN: "test",
		},
	}
	db := InitDB(cfg)

	assert.NotNil(t, db)
}
