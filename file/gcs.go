package file

import (
	"log"
	"mime/multipart"
)

type gcsImpl struct{}

func newGCS() storageOp {
	return &gcsImpl{}
}

func (g *gcsImpl) uploadFileTo(files []*multipart.FileHeader, uploadFolder string) error {
	log.Println("uploadFileTo to gcs")
	return nil
}

func (g *gcsImpl) getType() string {
	return GCS
}
