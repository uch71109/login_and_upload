package rest

import (
	"net/http"

	"github.com/gobuffalo/buffalo"

	"login_and_upload/model"
)

const (
	indexHTML = "index.html"
)

func (ri *restImpl) index(c buffalo.Context) error {
	c.Set(cookieUser, &model.User{})
	c.Set(cookieFile, &model.File{})
	return c.Render(http.StatusOK, ri.renderEngine.HTML(indexHTML))
}
