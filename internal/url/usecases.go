package url

import (
	"context"

	"github.com/google/uuid"
)

type UseCases struct {
	repo URLRepository
}

func NewUseCases(repo URLRepository) *UseCases {
	return &UseCases{
		repo: repo,
	}
}

func (u *UseCases) Delete(ctx context.Context, authorID uuid.UUID, urlID string) error {
	url, err := u.repo.ByID(ctx, urlID)
	if err != nil {
		return err
	}

	if url.AuthorID != authorID {
		return ErrUserIsNotAuthor
	}

	if err = u.repo.Delete(ctx, urlID); err != nil {
		return err
	}

	return nil
}

func (u *UseCases) Update(ctx context.Context, authorID uuid.UUID, url *URL) error {
	url, err := u.repo.ByID(ctx, url.ID)
	if err != nil {
		return err
	}

	if url.AuthorID != authorID {
		return ErrUserIsNotAuthor
	}

	if err = u.repo.Update(ctx, url); err != nil {
		return err
	}

	return nil
}
