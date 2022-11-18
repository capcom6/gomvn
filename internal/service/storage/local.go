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
		return false, err
	}

	return !info.IsDir(), nil
}

func (a *localAdapter) ListItems(pathname string) ([]fileInfo, error)  {
	fullpath := a.fullpath(pathname)

	isFile, err := a.IsRegularFile(pathname)
	if err != nil {
		return nil, err
	}
	if isFile {
		return nil, fmt.Errorf("is not directory at %s", fullpath)
	}

	f, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}

	fileinfos, err := f.Readdir(0)
	_ = f.Close()
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("is directory at %s", fullpath)
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
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	
	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func (a *localAdapter) fullpath(name string) string {
	return path.Clean(path.Join(a.root, name))
}