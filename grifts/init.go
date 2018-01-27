package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/xtrem-creative/ldap_password_selfservice/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
