package rest

import (
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"login_and_upload/model"
)

const (
	siginHTML = "users/login.plush.html"
	loginPath = "/users/login"
)

func (ri *restImpl) setCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get(cookieCurrentUserID); uid != nil {
			user, err := ri.dataOp.User.GetByEmail(uid.(string))
			if err != nil {
				return errors.WithStack(err)
			}
			if user == nil {
				return errors.WithStack(errors.New("Not Found"))
			}
			c.Set(cookieCurrentUser, user)
		}
		return next(c)
	}
}

func (ri *restImpl) login(c buffalo.Context) error {
	u := &model.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	user, err := ri.dataOp.User.GetByEmail(strings.TrimSpace(u.Email))
	if err != nil {
		return errors.WithStack(err)
	}
	if user == nil {
		return ri.renderError(c, u)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password)); err != nil {
		return ri.renderError(c, u)
	}

	c.Session().Set(cookieCurrentUserID, user.Email)
	c.Session().Set(cookieCurrentUserRole, user.Permission)
	c.Flash().Add(flashSucess, ri.t.Translate(c, "welcome_greeting"))

	return c.Redirect(http.StatusFound, indexURL)
}

func (ri *restImpl) loginPage(c buffalo.Context) error {
	if uid := c.Session().Get(cookieCurrentUserID); uid != nil {
		return c.Redirect(http.StatusFound, indexURL)
	}
	c.Set(cookieUser, model.User{})
	return c.Render(http.StatusOK, ri.renderEngine.HTML(siginHTML))
}

func (ri *restImpl) logout(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add(flashSucess, ri.t.Translate(c, "logged_out"))
	return c.Redirect(http.StatusFound, loginPath)
}

func (ri *restImpl) authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get(cookieCurrentUserID); uid != nil &&
			c.Request().URL.String() == indexURL {
			return next(c)
		}
		c.Flash().Add(flashDanger, ri.t.Translate(c, "unauthorized"))
		return c.Redirect(http.StatusFound, loginPath)
	}
}

func (ri *restImpl) renderError(c buffalo.Context, u *model.User) error {
	c.Set(cookieUser, u)
	c.Flash().Add(flashDanger, ri.t.Translate(c, "unauthorized"))
	return c.Render(http.StatusUnauthorized, ri.renderEngine.HTML(siginHTML))
}
