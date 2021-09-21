package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrDBAccess        = errors.New("db access error")
	ErrNotReaching     = errors.New("not reaching")
	ErrRequest         = errors.New("request error")
	ErrTokenWasExpired = errors.New("token was Expired")
	ErrInvalidVerify   = errors.New("invalid verify")
	ErrInvalidValue    = errors.New("invalid value")
)
