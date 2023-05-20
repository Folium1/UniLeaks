package delivery

import (
	"context"
	"errors"
	"log"
	"time"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) getLogin(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, "userLogin.html", nil)
	if err != nil {
		c.AbortWithStatus(InternalServerError)
		return
	}
}

func (h Handler) postLogin(c *gin.Context) {
	user := models.User{
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	// Retrieve user from database
	dbUser, err := h.userService.GetByMail(ctx, user.Email)
	if err != nil {
		log.Println("Error while getting user from db:", err)
		responseWithErr(c, "userLogin.html", errors.New("Invalid email or password"))
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		log.Println("Passwords are not similar:", err)
		responseWithErr(c, "userLogin.html", errors.New("Invalid email or password"))
		return
	}

	// Authentication successful
	Middleware.AuthorizeUser(c, dbUser.ID)
	c.Redirect(StatusSeeOther, "/leaks/")
}

func (h Handler) logOut(c *gin.Context) {
	Middleware.LogOut(c)
	c.Redirect(StatusSeeOther, "/user/login")
}
