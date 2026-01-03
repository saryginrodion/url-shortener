package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"roadmap.restapi/internal/ctxlogging"
	"roadmap.restapi/internal/database"
	"roadmap.restapi/internal/errormapper"
	"roadmap.restapi/internal/postgres"
)

type PostgresUserRepository struct {
	db     *sqlx.DB
	errMap *errormapper.ErrorMapper
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
		errMap: errormapper.NewErrorMapper(
			errormapper.NewMapping(database.ErrUniqueViolation, ErrUserAlreadyExists),
			errormapper.NewMapping(database.ErrNotFound, ErrUserNotFound),
		),
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	log := ctxlogging.Get(ctx)
	rows, err := r.db.NamedQueryContext(ctx, `
		INSERT INTO users (email, password_hash)
		VALUES (:email, :password_hash)
		RETURNING *
	`, user)
	if err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	rows.Next()
	if err = rows.Err(); err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	if err = rows.StructScan(user); err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	log := ctxlogging.Get(ctx)
	rows, err := r.db.NamedQueryContext(ctx, `
		UPDATE users 
		SET email = :email,
			password_hash = :password_hash
		WHERE id = :id
		RETURNING *
	`, user)

	if err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	if !rows.Next() {
		return ErrUserNotFound
	}

	if err = rows.Err(); err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	if err = rows.StructScan(user); err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := ctxlogging.Get(ctx)
	res, err := r.db.ExecContext(ctx, r.db.Rebind(`DELETE FROM users WHERE id = ?`), id)
	if err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *PostgresUserRepository) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	log := ctxlogging.Get(ctx)
	user := &User{}
	err := r.db.GetContext(ctx, user, r.db.Rebind(`SELECT * FROM users WHERE id = ?`), id)
	if err != nil {
		return nil, r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	return user, nil
}

func (r *PostgresUserRepository) ByEmail(ctx context.Context, email string) (*User, error) {
	log := ctxlogging.Get(ctx)
	user := &User{}
	err := r.db.GetContext(ctx, user, r.db.Rebind(`SELECT * FROM users WHERE email = ?`), email)
	if err != nil {
		return nil, r.errMap.MapAndLogUnmatched(postgres.TranslateError(err, log), log)
	}

	return user, nil
}
