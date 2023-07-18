package http

import (
	"context"
	"fmt"
	"time"

	"leaks/pkg/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	userLoginPage   = "userLogin.html"
	invalidCredsErr = "Невірний мейл або пароль"
)

func (u *UserHandler) PostLogin(c *gin.Context) {
	user := models.User{
		Email:    c.PostForm("email"),
		NickName: c.PostForm("nickName"),
		Password: c.PostForm("password"),
	}
	if err := validateUserInput(u.userService, user); err != nil {
		logger.Error(fmt.Sprint("Error while validating user input:", err))
		c.JSON(400, gin.H{"error": err})
		return
	}
	context, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	dbUser, err := u.userService.UserByNick(context, user.NickName)
	if err != nil {
		logger.Error(fmt.Sprint("Error while getting user from db:", err))
		c.JSON(400, gin.H{"error": err})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		logger.Error(fmt.Sprint("Passwords are not similar:", gin.H{"error": err}))
		c.JSON(400, invalidCredsErr)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Email), []byte(user.Email)); err != nil {
		logger.Error(fmt.Sprint("Mails are not similar:", err))
		c.JSON(400, gin.H{"error": invalidCredsErr})
		return
	}

	logger.Info(fmt.Sprintf("User: %v logged in", dbUser.ID))

	middleware.AuthorizeUser(c, dbUser.ID)
	c.Status(200)
}

// LogOut logs out user by deleting cookies
func (u *UserHandler) LogOut(c *gin.Context) {
	middleware.LogOut(c)
	c.Redirect(303, "/user/login")
}
