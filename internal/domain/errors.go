package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrBlocked       = errors.New("user is blocked")
	ErrForbidden     = errors.New("forbidden")
	ErrUnknownState  = errors.New("unknown state")
)
