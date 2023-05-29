package user

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
	errHandler "uniLeaks/delivery/err"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	userLoginPage = "userLogin.html"
)

// login handles GET request to /user/login
func (u UserHandler) Login(c *gin.Context) {
	err := u.tmpl.ExecuteTemplate(c.Writer, userLoginPage, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// postLogin handles POST request to /user/login
func (u UserHandler) PostLogin(c *gin.Context) {
	user := models.User{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}
	// Validate user input
	if err := validateUserInput(user); err != nil {
		errHandler.ResponseWithErr(c, userLoginPage, err)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	// Retrieve user from database
	dbUser, err := u.userService.GetByMail(ctx, user.Email)
	if err != nil {
		log.Println("Error while getting user from db:", err)
		errHandler.ResponseWithErr(c, userLoginPage, errors.New("Невірний мейл або пароль"))
		return
	}

	// Compare password's hash from db with the given password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		log.Println("Passwords are not similar:", err)
		errHandler.ResponseWithErr(c, userLoginPage, errors.New("Невірний мейл або пароль"))
		return
	}

	// Authentication successful
	middleware.AuthorizeUser(c, dbUser.ID)
	c.Redirect(http.StatusSeeOther, "/leaks/")
}

// logOut logs out user by deleting cookies
func (u UserHandler) LogOut(c *gin.Context) {
	middleware.LogOut(c)
	c.Redirect(http.StatusSeeOther, "/user/login")
}
