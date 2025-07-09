package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/model"
	"github.com/nats-io/nats.go/jetstream"
)

// Mock errors
var (
	ErrMockDB      = errors.New("mock db error")
	ErrMockPublish = errors.New("mock publish error")
)

// MockUserRepository implements UserRepository interface
type MockUserRepository struct {
	users []model.User
	err   error
}

func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (model.User, error) {
	if m.err != nil {
		return model.User{}, m.err
	}
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return model.User{}, sql.ErrNoRows
}

func (m *MockUserRepository) CreateTx(ctx context.Context, tx *sql.Tx, user *model.User) error {
	if m.err != nil {
		return m.err
	}
	user.ID = int64(len(m.users) + 1)
	m.users = append(m.users, *user)
	return nil
}

func (m *MockUserRepository) WithTransaction(ctx context.Context, txFunc func(context.Context, *sql.Tx) error) error {
	if m.err != nil {
		return m.err
	}
	return txFunc(ctx, &sql.Tx{})
}

// MockPublisher implements Publisher interface
type MockPublisher struct {
	err error
}

func (m *MockPublisher) Publish(ctx context.Context, subject string, request interface{}) (*jetstream.PubAck, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &jetstream.PubAck{}, nil
}

// Test data
var mockUsers = []model.User{
	{
		ID:        1,
		Name:      "John Doe",
		CreatedAt: time.Now().UnixMicro(),
		UpdatedAt: time.Now().UnixMicro(),
	},
	{
		ID:        2,
		Name:      "Jane Doe",
		CreatedAt: time.Now().UnixMicro(),
		UpdatedAt: time.Now().UnixMicro(),
	},
}
