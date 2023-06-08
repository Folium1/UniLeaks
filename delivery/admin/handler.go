package admin

import (
	"html/template"
	admin "uniLeaks/admin/service"
)

type AdminHandler struct {
	tmpl *template.Template
	leak *admin.LeakService
	user *admin.UserService
}

func New(tmpl *template.Template) *AdminHandler {
	driveService := admin.NewLeakService()
	userService := admin.NewUserService()
	return &AdminHandler{tmpl, &driveService, &userService}
}
