package domain

import (
	"context"
	"time"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// UserUsecase represent the User usecases
type UserUsecase interface {
	Fetch(ctx context.Context, num int64) ([]User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *User) error
	Store(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
}

// UserRepository represent the User repository contract
type UserRepository interface {
	Fetch(ctx context.Context, num int64) ([]User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *User) error
	Store(ctx context.Context, a *User) error
	Delete(ctx context.Context, id int64) error
}
