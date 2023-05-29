package user

import (
	"context"
	"errors"
	"net/http"
	"time"
	auth "uniLeaks/auth/delivery/http"
	errHandler "uniLeaks/delivery/err"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
)

const registerPage = "userRegister.html"

var (
	registerErr = errors.New("Під час регістрації сталась помилка, спробуйте ще раз")
	middleware  = auth.New()
)

// getRegister handles GET request to /user/register
func (u UserHandler) Register(c *gin.Context) {
	err := u.tmpl.ExecuteTemplate(c.Writer, registerPage, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// postRegister handles POST request to /user/register
func (u UserHandler) PostRegister(c *gin.Context) {
	newUser := models.User{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}
	// Validate user input
	if err := validateUserInput(newUser); err != nil {
		errHandler.ResponseWithErr(c, registerPage, err)
	}
	// Hash password
	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		errHandler.ResponseWithErr(c, registerPage, registerErr)
		return
	}
	newUser.Password = string(hashedPassword)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Save user to database
	userId, err := u.userService.Create(ctx, newUser)
	if err != nil {
		errHandler.ResponseWithErr(c, registerPage, registerErr)
		return
	}
	// Authentication successful
	middleware.AuthorizeUser(c, userId)
	c.Redirect(http.StatusFound, "/leaks")
}
