package storage

import "errors"

var (
	ErrNotADirectory = errors.New("not a directory")
	ErrNotAFile      = errors.New("not a file")
)
