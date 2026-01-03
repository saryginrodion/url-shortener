package database

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("record not found")
	ErrUniqueViolation     = errors.New("unique constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrNotNullViolation    = errors.New("not null violation")
	ErrCheckViolation      = errors.New("check constraint violation")
	ErrConnectionDone      = errors.New("connection already returned to pool")
	ErrTransactionDone     = errors.New("transaction already closed")

	ErrConnectionProblem   = errors.New("database connection problem")
)
