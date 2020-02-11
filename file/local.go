package file

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	mkdirPerm = 0755
)

type localImpl struct{}

func newLocal() storageOp {
	return &localImpl{}
}

func (l *localImpl) uploadFileTo(files []*multipart.FileHeader, uploadFolder string) error {
	// TODO: to use worker with goroutine
	for _, v := range files {
		file, err := v.Open()
		if err != nil {
			return errors.WithStack(err)
		}
		dir := filepath.Join(".", uploadFolder)
		if err := os.MkdirAll(dir, mkdirPerm); err != nil {
			return errors.WithStack(err)
		}
		f, err := os.Create(filepath.Join(dir, v.Filename))
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()

		if _, err = io.Copy(f, file); err != nil {
			return errors.WithStack(err)
		}
	}

	log.Println("uploadFileTo to local")
	return nil
}

func (l *localImpl) getType() string {
	return Local
}
