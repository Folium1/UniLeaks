package delivery

import (
	"html/template"
	userAdmin "leaks/admin/service"
	"leaks/logger"
)

var logg = logger.NewLogger()

type Handler struct {
	tmpl        *template.Template
	userService userAdmin.UserService
}

// New returns a new instance of the auth handler.
func New() Handler {
	tmpl := template.Must(template.ParseGlob("templates/*"))
	userService := userAdmin.NewUserService()
	return Handler{tmpl: tmpl, userService: userService}
}
