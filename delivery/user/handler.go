package user

import (
	"html/template"
	"uniLeaks/user"
	userRepository "uniLeaks/user/repository/mysql"
	userService "uniLeaks/user/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	router      *gin.Engine
	tmpl        *template.Template
	userService user.Repository
}

func New(router *gin.Engine, tmpl *template.Template) UserHandler {
	userService := userService.New(userRepository.New())
	return UserHandler{router: router, tmpl: tmpl, userService: userService}
}
