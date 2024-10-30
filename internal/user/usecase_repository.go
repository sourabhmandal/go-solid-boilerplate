package user

import (
	"errors"
	"golang.org/x/net/context"
	"strconv"
)

type UserUseCaseSqlc struct {
	userRepository Querier
}

func NewUserUseCase(userRepo Querier) *UserUseCaseSqlc {
	return &UserUseCaseSqlc{userRepository: userRepo}
}

// RegisterUser registers a new user in the system.
func (u *UserUseCaseSqlc) RegisterUser(ctx context.Context, name, email string) error {
	// Check if the user already exists based on email
	existingUser, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Create a new user in the schema
	_, err = u.userRepository.CreateUser(ctx, &CreateUserParams{
		Email: email,
		Name:  name,
		Bio:   nil, // Assuming empty bio for new user
	})
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by their ID.
func (u *UserUseCaseSqlc) GetUserByID(ctx context.Context, userID string) (*User, error) {
	// Convert userID to int64 if necessary
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get the user by ID from the schema
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
