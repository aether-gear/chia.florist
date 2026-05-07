package storage

import (
	"io"
)

type Provider interface {
	Upload(input UploadInput) (Object, error)
	Delete(key string) error
	Exists(key string) (bool, error)
}

type UploadInput struct {
	Key           string
	File          io.Reader
	ContentType   string
	ContentLength int64
}
