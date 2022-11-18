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
	"github.com/valyala/bytebufferpool"
)

func NewStorage(cfg *config.Storage) *Storage {
	return &Storage{
		// adapter: newLocalAdapter(cfg.Options),
		adapter: newAwsAdapter(cfg.Options),
	}
}

type Storage struct {
	adapter storageAdapter
}

func (s *Storage) ListArtifacts(repo string) ([]*entity.Artifact, error) {
	result := []*entity.Artifact{}

	items, err := s.adapter.ListItems(repo)
	if err != nil {
		return nil, err
	}

	for _, info := range items {
		if info.IsDir {
			inner, err := s.ListArtifacts(path.Join(repo, info.Name))
			if err != nil {
				return nil, err
			}

			result = append(result, inner...)
			continue
		}

		if strings.HasSuffix(info.Name, ".pom") {
			pathname := path.Join(repo, info.Name)
			pathname = path.Clean(strings.Replace(pathname, "\\", "/", -1))
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
		return nil, "", err
	}
	
	if !isRegularFile{
		w := &bytebufferpool.ByteBuffer{}
		err := s.createDirIndex(w, pathname)
		if err != nil {
			return nil, "", err
		}
		
		return io.NopCloser(bytes.NewReader(w.Bytes())), "text/html", nil	
	}
	
	reader, err := s.adapter.Read(pathname)
	ext := path.Ext(pathname)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return reader, contentType, err
}

func (s *Storage) Write(pathname string, r io.Reader) error {
	return s.adapter.Write(pathname, r)
}

func (s *Storage) createDirIndex(w io.Writer, pathname string) error {
	fileinfos, err := s.adapter.ListItems(pathname)
	if err != nil {
		return err
	}

	basePathEscaped := html.EscapeString(string(pathname))
	_, _ = fmt.Fprintf(w, "<html><head><title>%s</title><style>.dir { font-weight: bold }</style></head><body>", basePathEscaped)
	_, _ = fmt.Fprintf(w, "<h1>%s</h1>", basePathEscaped)
	_, _ = fmt.Fprintf(w, "<ul>")

	if len(basePathEscaped) > 1 {
		parentPathEscaped := html.EscapeString(string(path.Dir(pathname)))
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
		pathEscaped := html.EscapeString(string(path.Join(pathname, name)))
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