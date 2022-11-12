package service

// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"os"
// 	"path"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// )

// func NewLocalStorage() *LocalStorage {
// 	return &LocalStorage{
// 		root: "data/repository",
// 	}
// }

// type LocalStorage struct {
// 	root string
// }

// func (s *LocalStorage) file(name string) string {
// 	return path.Clean(path.Join(s.root, name))
// }

// func (s *LocalStorage) Open(path string) (io.ReadCloser, error) {
// 	if !s.FileExists(path) {
// 		return nil, os.ErrNotExist
// 	}

// 	fullpath := s.file(path)
// 	file, err := os.Open(fullpath)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't open file at %s: %w", path, err)
// 	}

// 	return file, nil
// }

// func (s *LocalStorage) FileExists(path string) bool {
// 	fullpath := s.file(path)
// 	if !strings.HasPrefix(fullpath, s.root) {
// 		return false
// 	}

// 	info, err := os.Stat(fullpath)
// 	if os.IsNotExist(err) {
// 		return false
// 	}
// 	return !info.IsDir()
// }

// func (s *LocalStorage) WriteFromRequest(c *fiber.Ctx, path string) error {
// 	file := s.file(path)
// 	fdir := dir(file)
// 	if err := os.MkdirAll(fdir, 0750); err != nil {
// 		return err
// 	}

// 	f, err := os.Create(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	w := bufio.NewWriter(f)

// 	if err := c.Request().BodyWriteTo(w); err != nil {
// 		return err
// 	}
// 	if err := w.Flush(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func dir(path string) string {
// 	index := strings.LastIndex(path, "/")
// 	return path[:index]
// }
