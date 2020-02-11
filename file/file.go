package file

import (
	"errors"
	"mime/multipart"
	"sync"
)

// NewFile returns a file manager
func NewFile(concurrentJobs int,
) File {
	return &fileImpl{
		conJobs: make(chan struct{}, concurrentJobs),
		storageOp: map[string]storageOp{
			Local: newLocal(),
			S3:    newS3(),
			Blob:  newBlob(),
			GCS:   newGCS(),
		},
	}
}

type fileImpl struct {
	conJobs   chan struct{}
	storageOp map[string]storageOp
}

func (fi *fileImpl) UploadFileTo(files []*multipart.FileHeader, uploadFolder string, types ...string) error {
	var wg sync.WaitGroup
	uploadErr := []error{}
	uploadErrMutex := &sync.Mutex{}

	// types: local/gcs/blob/s3
	for _, t := range types {
		op := fi.storageOp[t]
		wg.Add(1)
		fi.limitJob(func() {
			defer wg.Done()
			if err := op.uploadFileTo(files, uploadFolder); err != nil {
				fi.appendErrorWithMutex(err, &uploadErr, uploadErrMutex)
			}
		})
	}
	wg.Wait()

	if len(uploadErr) == 0 {
		return nil
	}
	return errors.New("upload file errors")
}

func (fi *fileImpl) limitJob(f func()) {
	fi.conJobs <- struct{}{}
	go func() {
		defer func() {
			<-fi.conJobs
		}()
		f()
	}()
}

func (fi *fileImpl) appendErrorWithMutex(err error, errs *[]error, mutex *sync.Mutex) {
	mutex.Lock()
	*errs = append(*errs, err)
	mutex.Unlock()
}
