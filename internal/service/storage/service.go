package storage

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"mime"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gomvn/gomvn/internal/config"
	"github.com/gomvn/gomvn/internal/entity"
)

const (
	DriverLocal = "local"
	DriverAWS   = "s3"
)

type Storage struct {
	adapter storageAdapter
}

func NewStorage(cfg *config.Storage) *Storage {
	var adapter storageAdapter
	switch cfg.Driver {
	case DriverLocal:
		adapter = newLocalAdapter(cfg.Options)
	case DriverAWS:
		adapter = newS3Adapter(cfg.Options)
	default:
		adapter = newLocalAdapter(map[string]string{"root": "data/repository"})
	}

	return &Storage{
		adapter: adapter,
	}
}

func (s *Storage) ListArtifacts(repo string) ([]*entity.Artifact, error) {
	result := []*entity.Artifact{}

	items, err := s.adapter.ListItems(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to list items at %s: %w", repo, err)
	}

	for _, info := range items {
		if info.IsDir {
			inner, recErr := s.ListArtifacts(path.Join(repo, info.Name))
			if recErr != nil {
				return nil, recErr
			}

			result = append(result, inner...)
			continue
		}

		if strings.HasSuffix(info.Name, ".pom") {
			pathname := path.Join(repo, info.Name)
			pathname = path.Clean(strings.ReplaceAll(pathname, "\\", "/"))
			pathname = pathname[strings.Index(pathname, "/")+1:]

			artifact := entity.NewArtifact(pathname, info.ModTime)
			result = append(result, artifact)
		}
	}

	return result, nil
}

func (s *Storage) Open(pathname string) (io.ReadCloser, string, error) {
	isRegularFile, err := s.adapter.IsRegularFile(pathname)
	if err != nil {
		return nil, "", fmt.Errorf("failed to stat file at %s: %w", pathname, err)
	}

	if !isRegularFile {
		var w bytes.Buffer

		if idxErr := s.createDirIndex(&w, pathname); idxErr != nil {
			return nil, "", idxErr
		}

		return io.NopCloser(bytes.NewReader(w.Bytes())), "text/html", nil
	}

	reader, err := s.adapter.Read(pathname)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file at %s: %w", pathname, err)
	}

	ext := path.Ext(pathname)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return reader, contentType, nil
}

func (s *Storage) Write(pathname string, r io.Reader) error {
	err := s.adapter.Write(pathname, r)
	if err != nil {
		return fmt.Errorf("failed to write file at %s: %w", pathname, err)
	}

	return nil
}

func (s *Storage) createDirIndex(w io.Writer, pathname string) error {
	fileinfos, err := s.adapter.ListItems(pathname)
	if err != nil {
		return fmt.Errorf("failed to list items at %s: %w", pathname, err)
	}

	basePathEscaped := html.EscapeString(pathname)
	_, _ = fmt.Fprintf(
		w,
		"<html><head><title>%s</title><style>.dir { font-weight: bold }</style></head><body>",
		basePathEscaped,
	)
	_, _ = fmt.Fprintf(w, "<h1>%s</h1>", basePathEscaped)
	_, _ = fmt.Fprintf(w, "<ul>")

	if len(basePathEscaped) > 1 {
		parentPathEscaped := html.EscapeString(path.Dir(pathname))
		_, _ = fmt.Fprintf(w, `<li><a href="/%s" class="dir">..</a></li>`, parentPathEscaped)
	}

	fm := make(map[string]fileInfo, len(fileinfos))
	filenames := make([]string, 0, len(fileinfos))

	for _, fi := range fileinfos {
		name := fi.Name
		fm[name] = fi
		filenames = append(filenames, name)
	}

	sort.Strings(filenames)
	for _, name := range filenames {
		pathEscaped := html.EscapeString(path.Join(pathname, name))
		fi := fm[name]
		auxStr := "dir"
		className := "dir"
		if !fi.IsDir {
			auxStr = fmt.Sprintf("file, %d bytes", fi.Size)
			className = "file"
		}
		_, _ = fmt.Fprintf(w, `<li><a href="/%s" class="%s">%s</a>, %s, last modified %s</li>`,
			pathEscaped, className, html.EscapeString(name), auxStr, fsModTime(fi.ModTime))
	}

	_, _ = fmt.Fprintf(w, "</ul></body></html>")

	return nil
}

func fsModTime(t time.Time) time.Time {
	return t.In(time.UTC).Truncate(time.Second)
}
