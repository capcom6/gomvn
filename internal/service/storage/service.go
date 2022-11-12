package storage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gomvn/gomvn/internal/entity"
)

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		root: "data/repository",
	}
}

type LocalStorage struct {
	root string
}

func (s *LocalStorage) file(name string) string {
	return path.Clean(path.Join(s.root, name))
}

func (s *LocalStorage) List(repo string) []*entity.Artifact {
	result := []*entity.Artifact{}
	repoPath := s.file(repo)
	_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".pom") {
			path = strings.Replace(path, "\\", "/", -1)
			path = strings.Replace(path, repoPath+"/", "", 1)
			artifact := entity.NewArtifact(path, info.ModTime())
			result = append(result, artifact)
		}
		return nil
	})
	return result
}

func (s *LocalStorage) Open(path string) (io.ReadCloser, error) {
	if !s.FileExists(path) {
		return nil, os.ErrNotExist
	}
	
	fullpath := s.file(path)
	file, err := os.Open(fullpath)
	if err != nil {
		return nil, fmt.Errorf("can't open file at %s: %w", path, err)
	}

	return file, nil
}

func (s *LocalStorage) FileExists(path string) bool {
	fullpath := s.file(path)
	if !strings.HasPrefix(fullpath, s.root) {
		return false
	}

	info, err := os.Stat(fullpath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (s *LocalStorage) WriteFromRequest(c *fiber.Ctx, pathname string) error {
	file := s.file(pathname)
	fdir := path.Dir(file)
	if err := os.MkdirAll(fdir, 0750); err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	if err := c.Request().BodyWriteTo(w); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
