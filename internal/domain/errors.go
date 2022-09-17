package domain

import "errors"

var (
	ErrFailedToCreate = errors.New("failed to create shortURL")
	ErrNotFound       = errors.New("shortURL not found")
)
