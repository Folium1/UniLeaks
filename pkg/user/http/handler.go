package http

import (
	"fmt"
	"html/template"
	"leaks/pkg/logger"
	userRepository "leaks/pkg/user/repository/mysql"
	userService "leaks/pkg/user/service"

	"github.com/joho/godotenv"
)

// Init initializes the environment variables
func init() {
	err := godotenv.Load()
	if err != nil {
		logg.Fatal(fmt.Sprint("Couldn't load local variables, err:", err))
	}
}

var logg = logger.NewLogger()

type UserHandler struct {
	tmpl        *template.Template
	userService *userService.UserUseCase
}

func New(tmpl *template.Template) *UserHandler {
	userService := userService.New(userRepository.New())
	return &UserHandler{tmpl: tmpl, userService: &userService}
}
