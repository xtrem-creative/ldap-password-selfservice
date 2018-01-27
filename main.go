package main

import (
	"log"

	"github.com/xtrem-creative/ldap_password_selfservice/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
