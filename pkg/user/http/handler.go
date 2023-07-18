package http

import (
	"fmt"

	auth "leaks/pkg/auth"
	logg "leaks/pkg/logger"
	userService "leaks/pkg/user/service"

	"github.com/joho/godotenv"
)

var middleware = auth.New()

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(fmt.Sprint("Couldn't load local variables, err:", err))
	}
}

var logger = logg.NewLogger()

type UserHandler struct {
	userService *userService.UserUseCase
}

func New() *UserHandler {
	userService := userService.New()
	return &UserHandler{userService: &userService}
}
