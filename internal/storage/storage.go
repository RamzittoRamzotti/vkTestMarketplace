package storage

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
	ErrUserExists   = errors.New("User exists")
)

type AdFilter struct {
	Page      int
	Limit     int
	SortBy    string // "created_at" или "price"
	SortOrder string // "asc" или "desc"
	MinPrice  float64
	MaxPrice  float64
}
