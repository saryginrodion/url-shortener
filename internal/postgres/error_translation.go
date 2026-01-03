package postgres

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx"
	"roadmap.restapi/internal/database"
	"roadmap.restapi/internal/errormapper"
)

var errMap = errormapper.NewErrorMapper(
	errormapper.NewMapping(sql.ErrNoRows, database.ErrNotFound),
	errormapper.NewMapping(sql.ErrConnDone, database.ErrConnectionDone),
	errormapper.NewMapping(sql.ErrTxDone, database.ErrTransactionDone),
	errormapper.NewPredicateMapping(
		func(err error) bool {
			var pgErr pgx.PgError
			return errors.As(err, &pgErr) && pgErr.Code == "23505"
		},
		database.ErrUniqueViolation,
	),
	errormapper.NewPredicateMapping(
		func(err error) bool {
			var pgErr pgx.PgError
			return errors.As(err, &pgErr) && pgErr.Code == "23514"
		},
		database.ErrCheckViolation,
	),
	errormapper.NewPredicateMapping(
		func(err error) bool {
			var pgErr pgx.PgError
			return errors.As(err, &pgErr) && pgErr.Code == "23502"
		},
		database.ErrNotNullViolation,
	),
	errormapper.NewPredicateMapping(
		func(err error) bool {
			var pgErr pgx.PgError
			return errors.As(err, &pgErr) && pgErr.Code == "23503"
		},
		database.ErrForeignKeyViolation,
	),
)

func TranslateError(err error, log *slog.Logger) error {
	if err == nil {
		return nil
	}

	if to, match := errMap.MapAndCheck(err); match {
		return to
	}

	log.Error("unknown database error", "err", err, "origin", "postgres/error_translation.go")
	return err
}
