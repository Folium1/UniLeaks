package delivery

import (
	userAdmin "leaks/pkg/admin/service"
	logg "leaks/pkg/logger"
)

var logger = logg.NewLogger()

type Handler struct {
	userService *userAdmin.AdminUserService
}

func New() *Handler {
	userService := userAdmin.NewAdminUserService()
	return &Handler{userService: userService}
}
