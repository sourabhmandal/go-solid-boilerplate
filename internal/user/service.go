package user

import (
	"authosaurous/internal/repository"

	"golang.org/x/net/context"
)

type UserService interface {
	RegisterUser(ctx context.Context, name, email string) error
	GetUserByID(ctx context.Context, userID string) (*repository.User, error)
	// Other methods...
}
