package user

import (
	"html/template"
	"uniLeaks/user"
	userRepository "uniLeaks/user/repository/mysql"
	userService "uniLeaks/user/service"
)

type UserHandler struct {
	tmpl        *template.Template
	userService user.Repository
}

func New(tmpl *template.Template) UserHandler {
	userService := userService.New(userRepository.New())
	return UserHandler{tmpl: tmpl, userService: userService}
}
