package storage

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/samber/lo"
)

const (
	OptionLogin    = "login"
	OptionPassword = "password"
	OptionEndpoint = "endpoint"
	OptionRegion   = "region"
	OptionBucket   = "bucket"
	OptionPrefix   = "prefix"
)

type s3Adapter struct {
	session *session.Session
	s3      *s3.S3

	bucket string
	prefix string
}

func newS3Adapter(options map[string]string) *s3Adapter {
	login := options[OptionLogin]
	password := options[OptionPassword]
	endpoint := options[OptionEndpoint]
	region := options[OptionRegion]

	cfg := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(login, password, "")).
		WithRegion(region).
		WithEndpoint(endpoint)

	s, _ := session.NewSession(cfg)

	return &s3Adapter{
		session: s,
		s3:      s3.New(s),

		bucket: options[OptionBucket],
		prefix: options[OptionPrefix],
	}
}

func (a *s3Adapter) IsRegularFile(pathname string) (bool, error) {
	items, err := a.ListItems(path.Dir(pathname))
	if err != nil {
		return false, err
	}

	_, filename := path.Split(pathname)

	for _, v := range items {
		if v.Name == filename || v.Name == filename+"/" {
			return !v.IsDir, nil
		}
	}

	return false, os.ErrNotExist
}

func (a *s3Adapter) ListItems(pathname string) ([]fileInfo, error) {
	prefix := a.fullname(pathname) + "/"

	out, err := a.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(a.bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(prefix),
	})
	if err != nil {
		if aerr, ok := lo.ErrorsAs[awserr.RequestFailure](err); ok {
			if aerr.StatusCode() == http.StatusNotFound {
				return nil, fs.ErrNotExist
			}
		}
		return nil, fmt.Errorf("failed to list items at %s: %w", pathname, err)
	}

	result := []fileInfo{}
	for _, v := range out.CommonPrefixes {
		result = append(result, fileInfo{
			IsDir:   true,
			Name:    strings.Replace(*v.Prefix, prefix, "", 1),
			Size:    0,
			ModTime: time.Now(),
		})
	}

	for _, v := range out.Contents {
		result = append(result, fileInfo{
			IsDir:   false,
			Name:    strings.Replace(*v.Key, prefix, "", 1),
			Size:    *v.Size,
			ModTime: *v.LastModified,
		})
	}

	return result, nil
}

func (a *s3Adapter) Read(pathname string) (io.ReadCloser, error) {
	resp, err := a.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(a.fullname(pathname)),
	})
	if err != nil {
		if aerr, ok := lo.ErrorsAs[awserr.Error](err); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket, s3.ErrCodeNoSuchKey:
				return nil, fs.ErrNotExist
			}
		}
		return nil, fmt.Errorf("failed to read file at %s: %w", pathname, err)
	}

	return resp.Body, nil
}

func (a *s3Adapter) Write(pathname string, r io.Reader) error {
	tmp, err := os.CreateTemp("", "aws-")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	_, err = io.Copy(tmp, r)
	if err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	_, err = tmp.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek temporary file: %w", err)
	}

	_, err = a.s3.PutObject(&s3.PutObjectInput{
		Body:   tmp,
		Bucket: aws.String(a.bucket),
		Key:    aws.String(a.fullname(pathname)),
	})

	if err != nil {
		return fmt.Errorf("failed to write file at %s: %w", pathname, err)
	}

	return nil
}

func (a *s3Adapter) fullname(pathname string) string {
	return path.Join(a.prefix, pathname)
}
