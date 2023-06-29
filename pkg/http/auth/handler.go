package delivery

import (
	"html/template"
	userAdmin "leaks/pkg/admin/service"
	"leaks/pkg/logger"
)

var logg = logger.NewLogger()

type Handler struct {
	tmpl        *template.Template
	userService userAdmin.AdminUserService
}

// New returns a new instance of the auth handler.
func New() *Handler {
	tmpl := template.Must(template.ParseGlob("templates/*"))
	userService := userAdmin.NewAdminUserService()
	return &Handler{tmpl: tmpl, userService: userService}
}
