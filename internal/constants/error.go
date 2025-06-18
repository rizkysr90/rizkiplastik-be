package constants

import "errors"

var (
	ErrAlreadyExists              = errors.New("already exists")
	ErrCodePostgreUniqueViolation = "23505"
)
