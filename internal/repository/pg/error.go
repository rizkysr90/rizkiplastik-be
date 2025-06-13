package pg

import "errors"

// Sentinel errors - define once, use everywhere
var (
	ErrAlreadyExists = errors.New("repository : already exists")
)
