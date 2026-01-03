package url

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID       string    `db:"id"`
	URL       string    `db:"url"`
	Name      string    `db:"name"`
	AuthorID  uuid.UUID `db:"author_id"`
	CreatedAt time.Time `db:"created_at"`
}
