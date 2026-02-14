package storage

import (
	"io"
	"time"
)

type storageAdapter interface {
	IsRegularFile(pathname string) (bool, error)
	ListItems(pathname string) ([]fileInfo, error)
	Read(pathname string) (io.ReadCloser, error)
	Write(pathname string, r io.Reader) error
}

type fileInfo struct {
	IsDir   bool
	Name    string
	Size    int64
	ModTime time.Time
}
