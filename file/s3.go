package file

import (
	"log"
	"mime/multipart"
)

type s3Impl struct{}

func newS3() storageOp {
	return &s3Impl{}
}

func (s *s3Impl) uploadFileTo(files []*multipart.FileHeader, uploadFolder string) error {
	log.Println("uploadFileTo to s3")
	return nil
}

func (s *s3Impl) getType() string {
	return S3
}
