package admin

import (
	"fmt"
	errHandler "leaks/pkg/err"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BanUser bans user by id
func (a *AdminHandler) BanUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't convert userId to int: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	err = a.user.BanUser(userIdInt)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't ban user: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	logg.Info(fmt.Sprintf("User %v was banned by %v", userId, ctx.MustGet("userId")))
	ctx.Redirect(200, "/admin/users")
}

// AllUsers gives a list of all users
func (a *AdminHandler) AllUsers(ctx *gin.Context) {
	users, err := a.user.AllUsers()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get all users: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	err = a.tmpl.ExecuteTemplate(ctx.Writer, "allUsers.html", users)
	if err != nil {
		logg.Error(err.Error())
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.ServerErr)
	}
}

// GetBannedUsers gives a list of all banned users
func (a *AdminHandler) GetBannedUsers(ctx *gin.Context) {
	bannedUsers, err := a.user.GetBannedUsers()
	if err != nil {
		logg.Error(err.Error())
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	err = a.tmpl.ExecuteTemplate(ctx.Writer, "BannedUsers.html", bannedUsers)
	if err != nil {
		logg.Error(err.Error())
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.ServerErr)
	}
}

// UnbanUser cancel ban of user
func (a *AdminHandler) UnbanUser(ctx *gin.Context) {
	userId := ctx.PostForm("userId")
	if userId == "" {
		logg.Error(errHandler.NoUserIdErr.Error())
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.NoUserIdErr)
	}
	err := a.user.UnbanUser(userId)
	if err != nil {
		logg.Error(err.Error())
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
}
