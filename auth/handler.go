package delivery

import (
	"html/template"
	userAdmin "uniLeaks/admin/service"
)

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
