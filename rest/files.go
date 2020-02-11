package rest

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"

	"login_and_upload/model"
)

const (
	uploadFile   = "uploadFile"
	uploadFolder = "uploads"

	defaultMaxMemory = 32 << 20 // 32 MB
)

func (ri *restImpl) upload(c buffalo.Context) error {
	file := &model.File{}
	if err := c.Bind(file); err != nil {
		return errors.WithStack(err)
	}

	role := c.Session().Get(cookieCurrentUserRole)
	storagePerms := ri.storagePerm[role.(string)]

	r := c.Request()
	if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
		return errors.WithStack(err)
	}
	files := r.MultipartForm.File[uploadFile]
	if len(files) == 0 {
		c.Flash().Add(flashDanger, ri.t.Translate(c, "file_not_found"))
		return c.Redirect(http.StatusFound, indexURL)
	}

	if err := ri.fileOp.UploadFileTo(files, uploadFolder, storagePerms...); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add(flashSucess, ri.t.Translate(c, "file_created"))
	return c.Redirect(http.StatusFound, indexURL)
}
