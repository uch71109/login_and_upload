package file

import (
	"mime/multipart"
)

const (
	// Local defines local type
	Local = "local"
	// S3 defines S3 type
	S3 = "s3"
	// Blob defines Blob type
	Blob = "blob"
	// GCS defines GCS type
	GCS = "gcs"
)

// File defines file operation
type File interface {
	UploadFileTo(files []*multipart.FileHeader, uploadFolder string, types ...string) error
}

// storageOp defines cloud storage operation
type storageOp interface {
	uploadFileTo(files []*multipart.FileHeader, uploadFolder string) error
	getType() string
}
