package dbpattern

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsConstraintViolation(err error, constraint string) bool {
	pgError := new(pgconn.PgError)

	ok := errors.As(err, &pgError)
	if ok && pgError.Code == pgerrcode.UniqueViolation && pgError.ConstraintName == constraint {
		return true
	}

	return false
}
