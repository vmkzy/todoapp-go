package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForgeinKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown")
)
