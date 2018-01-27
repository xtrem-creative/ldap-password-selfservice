package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

// LDAP config
var ldapURL string
var ldapMethod string
var ldapBindDn string
var ldapPassword string
var ldapBase string
var ldapFilter string

var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		// Init app
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_ldap_password_selfservice_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
			return app
		}
		app.Use(T.Middleware())

		app.GET("/", ShowFormHandler)
		app.PUT("/", FormHandler)

		app.ServeFiles("/", assetsBox) // serve files from the public directory

		// Get mandatory env vars
		if ldapURL, err = envy.MustGet("SELFSERVICE_LDAP_URL"); err != nil {
			app.Stop(errors.Wrap(err, "SELFSERVICE_LDAP_URL env var is not set"))
			return app
		}

		if ldapMethod, err = envy.MustGet("SELFSERVICE_LDAP_METHOD"); err != nil {
			app.Stop(errors.Wrap(err, "SELFSERVICE_LDAP_METHOD env var is not set"))
			return app
		}

		if ldapBindDn, err = envy.MustGet("SELFSERVICE_LDAP_BIND_DN"); err != nil {
			app.Stop(errors.Wrap(err, "SELFSERVICE_LDAP_BIND_DN env var is not set"))
			return app
		}

		if ldapPassword, err = envy.MustGet("SELFSERVICE_LDAP_PASSWORD"); err != nil {
			app.Stop(errors.Wrap(err, "SELFSERVICE_LDAP_PASSWORD env var is not set"))
			return app
		}

		if ldapBase, err = envy.MustGet("SELFSERVICE_LDAP_BASE"); err != nil {
			app.Stop(errors.Wrap(err, "SELFSERVICE_LDAP_BASE env var is not set"))
			return app
		}

		if ldapFilter, err = envy.MustGet("SELFSERVICE_LDAP_FILTER"); err != nil {
			app.Stop(errors.Wrap(err, "SELFSERVICE_LDAP_FILTER env var is not set"))
			return app
		}
	}

	return app
}
