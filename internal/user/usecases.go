package user

import (
	"context"

	"github.com/google/uuid"
)

type UseCases struct {
	repo UserRepository
	hasher PasswordHasher
}

func NewUseCases(repo UserRepository, hahser PasswordHasher) *UseCases {
	return &UseCases{
		repo: repo,
		hasher: hahser,
	}
}

func (u *UseCases) NewUser(ctx context.Context, email string, password string) (*User, error) {
	passwordHashed := u.hasher.Hash(password)
	user := &User{
		ID: uuid.New(),
		Email: email,
		PasswordHash: passwordHashed,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UseCases) ByEmailAndPassword(ctx context.Context, email string, password string) (*User, error) {
	user, err := u.repo.ByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !u.hasher.Check(password, user.PasswordHash) {
		return nil, ErrPasswordIncorrect
	}

	return user, nil
}
