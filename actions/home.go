package actions

import (
	"crypto/tls"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	"gopkg.in/ldap.v2"
)

// ShowFormHandler shows the password change form
func ShowFormHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

// FormHandler handles the password change form
func FormHandler(c buffalo.Context) error {
	f := struct {
		Username        string `form:"Username"`
		CurrentPassword string `form:"CurrentPassword"`
		NewPassword     string `form:"NewPassword"`
		ConfirmPassword string `form:"ConfirmPassword"`
	}{}
	if err := c.Bind(&f); err != nil {
		return errors.WithStack(err)
	}

	if f.NewPassword != f.ConfirmPassword {
		c.Flash().Add("danger", T.Translate(c, "password-mismatch"))
		return c.Redirect(301, "/")
	}

	// Dial LDAP server
	l, err := ldap.Dial("tcp", ldapURL)
	if err != nil {
		c.Logger().Error(err)
		c.Flash().Add("danger", T.Translate(c, "an-error-occurred"))
		return c.Redirect(301, "/")
	}
	defer l.Close()

	if ldapMethod == "tls" {
		// Reconnect with TLS
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			c.Logger().Error(err)
			c.Flash().Add("danger", T.Translate(c, "an-error-occurred"))
			return c.Redirect(301, "/")
		}
	} else if ldapMethod != "plain" {
		c.Logger().Error(errors.Wrap(err, "SSL method not supported yet."))
		c.Flash().Add("danger", T.Translate(c, "an-error-occurred"))
		return c.Redirect(301, "/")
	}

	err = l.Bind(ldapBindDn, ldapPassword)
	if err != nil {
		c.Logger().Error(err)
		c.Flash().Add("danger", T.Translate(c, "an-error-occurred"))
		return c.Redirect(301, "/")
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		ldapBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(ldapFilter, f.Username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		c.Logger().Error(err)
		c.Flash().Add("danger", T.Translate(c, "an-error-occurred"))
		return c.Redirect(301, "/")
	}

	if len(sr.Entries) != 1 {
		c.Flash().Add("danger", T.Translate(c, "invalid-credentials"))
		return c.Redirect(301, "/")
	}

	userDn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userDn, f.CurrentPassword)
	if err != nil {
		c.Flash().Add("danger", T.Translate(c, "invalid-credentials"))
		return c.Redirect(301, "/")
	}

	// Issue a password modify request
	passwordModifyRequest := ldap.NewPasswordModifyRequest(userDn, f.CurrentPassword, f.NewPassword)
	_, err = l.PasswordModify(passwordModifyRequest)

	if err != nil {
		c.Logger().Error(err)
		c.Flash().Add("danger", T.Translate(c, "password-could-not-be-changed"))
		return c.Redirect(301, "/")
	}

	c.Flash().Add("success", T.Translate(c, "password-changed"))

	return c.Redirect(301, "/")
}
