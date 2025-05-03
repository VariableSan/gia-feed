package storage

import "errors"

var (
	ErrFeedExists   = errors.New("feed already exists")
	ErrFeedNotFound = errors.New("feed not found")
)
