package user

import (
	"context"
	"errors"
	"fmt"
	auth "leaks/auth"
	errHandler "leaks/err"
	"leaks/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const registerPage = "userRegister.html"

var (
	registerErr = errors.New("Під час регістрації сталась помилка, спробуйте ще раз")
	middleware  = auth.New()
)

// GetRegister handles GET request to /user/register
func (u *UserHandler) Register(ctx *gin.Context) {
	err := u.tmpl.ExecuteTemplate(ctx.Writer, registerPage, nil)
	if err != nil {
		logg.Error(fmt.Sprint("Error while executing template:", err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// PostRegister handles POST request to /user/register
func (u *UserHandler) PostRegister(ctx *gin.Context) {
	newUser := models.User{
		Email:    ctx.PostForm("email"),
		NickName: ctx.PostForm("nickName"),
		Password: ctx.PostForm("password"),
		IsBanned: false,
		IsAdmin:  false,
	}
	// Validate user input
	if err := validateUserInput(newUser); err != nil {
		logg.Error(fmt.Sprint("Error while validating user input:", err))
		errHandler.ResponseWithErr(ctx, u.tmpl, registerPage, err)
	}
	// Check if user with such nickname already banned
	if err := checkIfMailIsBanned(u.userService, newUser.Email); err != nil {
		// Check for a particular err
		if err == errHandler.UserIsBannedErr {
			u.tmpl.ExecuteTemplate(ctx.Writer, registerPage, err)
		} else {
			errHandler.ResponseWithErr(ctx, u.tmpl, errHandler.ErrPage, err)
		}
	}
	// Hash password and mail
	hashedPassword, err := hashString(newUser.Password)
	if err != nil {
		logg.Error(fmt.Sprint("Error while hashing password:", err))
		errHandler.ResponseWithErr(ctx, u.tmpl, registerPage, registerErr)
		return
	}
	hashedMail, err := hashString(newUser.Email)
	if err != nil {
		logg.Error(fmt.Sprint("Error while hashing mail:", err))
		errHandler.ResponseWithErr(ctx, u.tmpl, registerPage, registerErr)
		return
	}
	// Set hashed password and mail to user struct
	newUser.Password = string(hashedPassword)
	newUser.Email = string(hashedMail)
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Save user to database
	userId, err := u.userService.CreateUser(context, newUser)
	if err != nil {
		logg.Error(fmt.Sprint("Error while saving user to db:", err))
		errHandler.ResponseWithErr(ctx, u.tmpl, registerPage, err)
		return
	}
	logg.Info(fmt.Sprintf("User with id %d registered", userId))
	// Authentication successful
	middleware.AuthorizeUser(ctx, userId)
	ctx.Redirect(http.StatusFound, "/leaks")
}
