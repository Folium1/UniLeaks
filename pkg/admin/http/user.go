package http

import (
	"fmt"

	errHandler "leaks/pkg/err"

	"github.com/gin-gonic/gin"
)

func (a *AdminHandler) BanUser(c *gin.Context) {
	userId := c.Param("userId")
	err := a.user.BanUser(userId)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't ban user: ", err))
		c.JSON(400, gin.H{"error": "Couldn't ban user"})
	}
	l.Info(fmt.Sprintf("User %v was banned by %v", userId, c.MustGet("userId")))
	c.Redirect(200, "/admin/users")
}

func (a *AdminHandler) AllUsers(c *gin.Context) {
	users, err := a.user.AllUsers()
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get all users: ", err))
		c.JSON(400, gin.H{"error": "Couldn't get all users"})
	}
	c.JSON(200, gin.H{"users": users})
}

func (a *AdminHandler) GetBannedUsers(c *gin.Context) {
	bannedUsers, err := a.user.BannedUsers()
	if err != nil {
		l.Error(err.Error())
		c.JSON(400, gin.H{"error": "Couldn't get banned users"})
	}
	c.JSON(200, gin.H{"bannedUsers": bannedUsers})
}

func (a *AdminHandler) UnbanUser(c *gin.Context) {
	userId := c.PostForm("userId")
	if userId == "" {
		l.Error(errHandler.NoUserIdErr.Error())
		c.JSON(400, gin.H{"error": errHandler.NoUserIdErr.Error()})
	}
	err := a.user.UnbanUser(userId)
	if err != nil {
		l.Error(err.Error())
		c.JSON(400, gin.H{"error": "Couldn't unban user"})
	}
}
