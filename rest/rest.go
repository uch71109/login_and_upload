package rest

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/logger"
	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/packr/v2"
	"github.com/golang/groupcache/lru"

	"login_and_upload/file"
	"login_and_upload/model"
)

const (
	htmlLayout  = "application.plush.html"
	sessionName = "_login_and_upload_session"

	defaultLocale = "en-US"

	flashSucess = "success"
	flashDanger = "danger"

	cookieUser            = "user"
	cookieFile            = "file"
	cookieCurrentUser     = "current_user"
	cookieCurrentUserID   = "current_user_id"
	cookieCurrentUserRole = "cookie_current_user_role"

	indexURL = "/"
)

var assetsBox = packr.New("app:assets", "../public")

type restImpl struct {
	renderEngine *render.Engine
	t            *i18n.Translator
	dataOp       model.DataOp
	fileOp       file.File
	cache        *lru.Cache
	storagePerm  map[string][]string
}

// RegisterRest register rest server
func RegisterRest(dataOp model.DataOp, fileOp file.File) *buffalo.App {
	renderEngine := render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: htmlLayout,

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,
	})

	// Setup and use translations
	t, _ := translations()
	ri := &restImpl{
		renderEngine: renderEngine,
		t:            t,
		dataOp:       dataOp,
		fileOp:       fileOp,
		storagePerm: map[string][]string{
			model.RoleAdmin:  []string{file.Local, file.S3, file.Blob, file.GCS},
			model.RoleNormal: []string{file.Local},
		},
	}

	app := buffalo.New(buffalo.Options{
		SessionName: sessionName,
		LogLvl:      logger.ErrorLevel,
	})

	app.Use(paramlogger.ParameterLogger)
	app.Use(csrf.New)
	app.Use(t.Middleware())

	// Auth Middlewares
	app.Use(ri.setCurrentUser)
	app.Use(ri.authorize)

	// Routes for index
	app.GET("/", ri.index)

	// Routes for users
	auth := app.Group("/users")
	auth.GET("/login", ri.loginPage)
	auth.POST("/login", ri.login)
	auth.DELETE("/logout", ri.logout)
	auth.Middleware.Skip(ri.authorize, ri.loginPage, ri.login, ri.logout)

	// Routes for files
	upload := app.Group("/files")
	upload.POST("/upload", ri.upload)
	upload.Middleware.Remove(ri.authorize)

	// serve files from the public directory
	app.ServeFiles("/", assetsBox)

	return app
}

func translations() (*i18n.Translator, error) {
	t, err := i18n.New(packr.New("app:locales", "../locales"), defaultLocale)
	if err != nil {
		return nil, err
	}
	return t, nil
}
