package admin

import (
	"fmt"
	"html/template"
	admin "leaks/admin/service"
	errHandler "leaks/err"
	"leaks/logger"

	"github.com/gin-gonic/gin"
)

var logg = logger.NewLogger()

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

// MainPage displays the main page of admin panel
func (a *AdminHandler) MainPage(ctx *gin.Context) {
	err := a.tmpl.ExecuteTemplate(ctx.Writer, "admin.html", nil)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.ServerErr)
	}
}
