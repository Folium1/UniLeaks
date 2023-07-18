package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"leaks/pkg/models"

	"github.com/gin-gonic/gin"
)

var (
	registerErr = errors.New("Під час регістрації сталась помилка, спробуйте ще раз")
)

func (u *UserHandler) PostRegister(c *gin.Context) {
	newUser := models.User{
		Email:    c.PostForm("email"),
		NickName: c.PostForm("nickName"),
		Password: c.PostForm("password"),
		IsBanned: false,
		IsAdmin:  false,
	}
	if err := validateUserInput(u.userService, newUser); err != nil {
		logger.Error(fmt.Sprint("Error while validating user input:", err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := hashString(newUser.Password)
	if err != nil {
		logger.Error(fmt.Sprint("Error while hashing password:", err))
		c.JSON(500, gin.H{"error": registerErr})
		return
	}
	hashedMail, err := hashString(newUser.Email)
	if err != nil {
		logger.Error(fmt.Sprint("Error while hashing mail:", err))
		c.JSON(500, gin.H{"error": registerErr})
		return
	}

	newUser.Password = string(hashedPassword)
	newUser.Email = string(hashedMail)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userId, err := u.userService.CreateUser(ctx, newUser)
	if err != nil {
		logger.Error(fmt.Sprint("Error while saving user to db:", err))
		c.JSON(500, gin.H{"error": registerErr})
		return
	}

	logger.Info(fmt.Sprintf("User with id = %d registered", userId))

	middleware.AuthorizeUser(c, userId)
	c.Status(200)
}
