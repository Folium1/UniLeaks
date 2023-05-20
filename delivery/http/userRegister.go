package delivery

import (
	"context"
	"errors"
	"time"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
)

var registerPage = "userRegister.html"

func (h Handler) getRegister(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, registerPage, nil)
	if err != nil {
		c.AbortWithStatus(InternalServerError)
		return
	}
}

func (h Handler) postRegister(c *gin.Context) {
	newUser := models.User{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}
	if err := validateUserInput(newUser); err != nil {
		responseWithErr(c, registerPage, err)
	}
	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		responseWithErr(c, registerPage, errors.New("Під час регістрації сталась помилка, спробуйте ще раз"))
		return
	}
	newUser.Password = string(hashedPassword)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId, err := h.userService.Create(ctx, newUser)
	if err != nil {
		responseWithErr(c, registerPage, errors.New("Під час регістрації сталась помилка, спробуйте ще раз"))
		return
	}
	Middleware.AuthorizeUser(c, userId)
	c.Redirect(StatusFound, "/leaks")
}
