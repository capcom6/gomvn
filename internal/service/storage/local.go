package storage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
)

const (
	OptionRoot = "root"
)

type localAdapter struct {
	root string
}

func newLocalAdapter(options map[string]string) *localAdapter {
	return &localAdapter{
		root: options[OptionRoot],
	}
}

func (a *localAdapter) IsRegularFile(pathname string) (bool, error) {
	fullpath := a.fullpath(pathname)

	info, err := os.Stat(fullpath)
	if err != nil {
		return false, fmt.Errorf("failed to stat file at %s: %w", pathname, err)
	}

	return !info.IsDir(), nil
}

func (a *localAdapter) ListItems(pathname string) ([]fileInfo, error) {
	fullpath := a.fullpath(pathname)

	isFile, err := a.IsRegularFile(pathname)
	if err != nil {
		return nil, err
	}
	if isFile {
		return nil, fmt.Errorf("%w: %s", ErrNotADirectory, pathname)
	}

	f, err := os.Open(fullpath)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory at %s: %w", pathname, err)
	}

	fileinfos, err := f.Readdir(0)
	_ = f.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read directory at %s: %w", pathname, err)
	}

	result := make([]fileInfo, len(fileinfos))
	for i, v := range fileinfos {
		result[i] = fileInfo{
			IsDir:   v.IsDir(),
			Name:    v.Name(),
			Size:    v.Size(),
			ModTime: v.ModTime(),
		}
	}
	return result, nil
}

func (a *localAdapter) Read(pathname string) (io.ReadCloser, error) {
	fullpath := a.fullpath(pathname)

	isFile, err := a.IsRegularFile(pathname)
	if err != nil {
		return nil, err
	}
	if !isFile {
		return nil, fmt.Errorf("%w: %s", ErrNotAFile, pathname)
	}

	file, err := os.Open(fullpath)
	if err != nil {
		return nil, fmt.Errorf("can't open file at %s: %w", pathname, err)
	}

	return file, nil
}

func (a *localAdapter) Write(pathname string, r io.Reader) error {
	file := a.fullpath(pathname)
	fdir := path.Dir(file)
	if err := os.MkdirAll(fdir, 0750); err != nil {
		return fmt.Errorf("failed to create directory at %s: %w", fdir, err)
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create file at %s: %w", file, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to write file at %s: %w", file, err)
	}

	if flErr := w.Flush(); flErr != nil {
		return fmt.Errorf("failed to flush file at %s: %w", file, flErr)
	}

	return nil
}

func (a *localAdapter) fullpath(name string) string {
	return path.Clean(path.Join(a.root, name))
}
