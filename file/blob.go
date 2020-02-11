package file

import (
	"log"
	"mime/multipart"
)

type blobImpl struct{}

func newBlob() storageOp {
	return &blobImpl{}
}

func (b *blobImpl) uploadFileTo(files []*multipart.FileHeader, uploadFolder string) error {
	log.Println("uploadFileTo to blob")
	return nil
}

func (b *blobImpl) getType() string {
	return Blob
}
