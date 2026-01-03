package url

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"roadmap.restapi/internal/ctxlogging"
	"roadmap.restapi/internal/database"
	"roadmap.restapi/internal/errormapper"
	"roadmap.restapi/internal/postgres"
)

type PostgresURLRepository struct {
	db     *sqlx.DB
	errMap *errormapper.ErrorMapper
}

func NewPostgresURLRepository(db *sqlx.DB) *PostgresURLRepository {
	return &PostgresURLRepository{
		db: db,
		errMap: errormapper.NewErrorMapper(
			errormapper.NewMapping(database.ErrNotFound, ErrURLNotFound),
			errormapper.NewMapping(database.ErrUniqueViolation, ErrURLAlreadyExists),
		),
	}
}

func (r *PostgresURLRepository) ByID(ctx context.Context, id string) (*URL, error) {
	log := ctxlogging.Get(ctx)
	var url URL
	err := r.db.GetContext(ctx, &url, r.db.Rebind("SELECT * FROM urls WHERE id = ?"), id)
	if err != nil {
		return nil, r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	return &url, nil
}

func (r *PostgresURLRepository) Update(ctx context.Context, url *URL) error {
	log := ctxlogging.Get(ctx)
	rows, err := r.db.NamedQueryContext(ctx, `UPDATE urls
	SET author_id = :author_id,
		url = :url,
		name = :name
	WHERE id = :id
	RETURNING *
	`, url)

	if err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	if !rows.Next() {
		return ErrURLNotFound
	}

	if err = rows.Err(); err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	if err = rows.StructScan(url); err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	return nil
}

func (r *PostgresURLRepository) Create(ctx context.Context, url *URL) error {
	log := ctxlogging.Get(ctx)
	rows, err := r.db.NamedQueryContext(ctx, `INSERT INTO urls (id, author_id, url, name)
	VALUES (:id, :author_id, :url, :name)
	RETURNING *
	`, url)

	if err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	if err = rows.Err(); err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	rows.Next()
	if err = rows.StructScan(url); err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	return nil
}

func (r *PostgresURLRepository) Delete(ctx context.Context, id string) error {
	log := ctxlogging.Get(ctx)
	res, err := r.db.ExecContext(ctx, r.db.Rebind(`DELETE FROM urls WHERE id = ?`), id)
	if err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	if rows == 0 {
		return ErrURLNotFound
	}

	return nil
}
func (r *PostgresURLRepository) ByUser(ctx context.Context, userID uuid.UUID) ([]URL, error) {
	log := ctxlogging.Get(ctx)
	urls := []URL{}
	err := r.db.SelectContext(ctx, &urls, r.db.Rebind(`SELECT * FROM urls WHERE author_id = ?`), userID)

	if err != nil {
		return urls, r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	return urls, nil
}
