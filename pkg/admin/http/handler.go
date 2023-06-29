package http

import (
	"fmt"
	"html/template"
	admin "leaks/pkg/admin/service"
	errHandler "leaks/pkg/err"
	"leaks/pkg/logger"

	"github.com/gin-gonic/gin"
)

var logg = logger.NewLogger()

type AdminHandler struct {
	tmpl *template.Template
	leak *admin.AdminLeakService
	user *admin.AdminUserService
}

func New(tmpl *template.Template) *AdminHandler {
	driveService := admin.NewAdminLeakService()
	userService := admin.NewAdminUserService()
	return &AdminHandler{tmpl, &driveService, &userService}
}

// MainPage displays the main page of admin panel
func (a *AdminHandler) MainPage(ctx *gin.Context) {
	err := a.tmpl.ExecuteTemplate(ctx.Writer, "admin.html", nil)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.ServerErr)
	}
}
