package url

import (
	"context"

	"github.com/google/uuid"
)

type URLRepository interface {
	ByID(ctx context.Context, id string) (*URL, error)
	Update(ctx context.Context, url *URL) error
	Create(ctx context.Context, url *URL) error
	Delete(ctx context.Context, id string) error
	ByUser(ctx context.Context, userID uuid.UUID) ([]URL, error)
}
