package user

import (
	"html/template"
	"leaks/logger"
	userRepository "leaks/user/repository/mysql"
	userService "leaks/user/service"
)

var logg = logger.NewLogger()

type UserHandler struct {
	tmpl        *template.Template
	userService *userService.UserUseCase
}

func New(tmpl *template.Template) *UserHandler {
	userService := userService.New(userRepository.New())
	return &UserHandler{tmpl: tmpl, userService: &userService}
}
