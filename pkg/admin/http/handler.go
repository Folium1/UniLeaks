package http

import (
	admin "leaks/pkg/admin/service"
	"leaks/pkg/logger"

	"github.com/gin-gonic/gin"
)

var l = logger.NewLogger()

type AdminHandler struct {
	leak *admin.AdminLeakService
	user *admin.AdminUserService
}

func New() *AdminHandler {
	driveService := admin.NewAdminLeakService()
	userService := admin.NewAdminUserService()
	return &AdminHandler{driveService, userService}
}

func (a *AdminHandler) MainPage(c *gin.Context) {
	// if user got here, he is an admin
	c.Status(200)
}
