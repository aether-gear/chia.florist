package storage

import (
	"io"
)

type UploadInput struct {
	Key           string
	File          io.Reader
	ContentType   string
	ContentLength int64
}

type ObjectResponse struct {
	Key         string
	ContentType string
}

type ObjectStorage interface {
	Upload(input UploadInput) (*ObjectResponse, error)
	Delete(key string) error
	Exists(key string) (bool, error)
}

type URLResolver interface {
	PublicURL(key string) string
	SignedURL(key string) (string, error)
}

type BucketManager interface {
	EnsureBucket() error
}

type Provider interface {
	ObjectStorage
	URLResolver
	BucketManager
}
