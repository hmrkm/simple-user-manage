package usecase

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrDBAccess = errors.New("db access error")
)
