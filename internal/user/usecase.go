package user

import "golang.org/x/net/context"

type UserUseCase interface {
	RegisterUser(ctx context.Context, name, email string) error
	GetUserByID(ctx context.Context, userID string) (*User, error)
	// Other methods...
}
