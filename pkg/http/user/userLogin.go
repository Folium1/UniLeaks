package user

import (
	"context"
	"errors"
	"fmt"
	errHandler "leaks/pkg/err"
	"leaks/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	userLoginPage = "userLogin.html"
)

// Login handles GET request to /user/login
func (u *UserHandler) Login(ctx *gin.Context) {
	err := u.tmpl.ExecuteTemplate(ctx.Writer, userLoginPage, nil)
	if err != nil {
		logg.Error(fmt.Sprint("Error while executing template:", err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// PostLogin handles POST request to /user/login
func (u *UserHandler) PostLogin(ctx *gin.Context) {
	user := models.User{
		Email:    ctx.PostForm("email"),
		NickName: ctx.PostForm("nickName"),
		Password: ctx.PostForm("password"),
	}
	// Validate user input
	if err := validateUserInput(user); err != nil {
		logg.Error(fmt.Sprint("Error while validating user input:", err))
		errHandler.ResponseWithErr(ctx, u.tmpl, userLoginPage, err)
		return
	}
	context, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()
	// Retrieve user from database
	dbUser, err := u.userService.GetByNick(context, user.NickName)
	if err != nil {
		logg.Error(fmt.Sprint("Error while getting user from db:", err))
		errHandler.ResponseWithErr(ctx, u.tmpl, userLoginPage, errors.New("Невірний мейл або пароль"))
		return
	}

	// Compare password's hash from db with the given password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		logg.Error(fmt.Sprint("Passwords are not similar:", err))
		errHandler.ResponseWithErr(ctx, u.tmpl, userLoginPage, errors.New("Невірний мейл або пароль"))
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Email), []byte(user.Email)); err != nil {
		logg.Error(fmt.Sprint("Mails are not similar:", err))
		errHandler.ResponseWithErr(ctx, u.tmpl, userLoginPage, errors.New("Невірний мейл або пароль"))
		return
	}
	// Log user login
	logg.Info(fmt.Sprintf("User: %v logged in", dbUser.ID))

	// Authentication successful
	middleware.AuthorizeUser(ctx, dbUser.ID)
	ctx.Redirect(http.StatusSeeOther, "/leaks/")
}

// LogOut logs out user by deleting cookies
func (u *UserHandler) LogOut(ctx *gin.Context) {
	middleware.LogOut(ctx)
	ctx.Redirect(http.StatusSeeOther, "/user/login")
}
