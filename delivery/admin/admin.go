package admin

import (
	"errors"
	"fmt"
	"io/ioutil"
	errHandler "leaks/err"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MainPage displays the main page of admin panel
func (a *AdminHandler) MainPage(ctx *gin.Context) {
	err := a.tmpl.ExecuteTemplate(ctx.Writer, "admin.html", nil)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.ServerErr)
	}
}

// UsersList handles the users list request
func (a *AdminHandler) FilesList(ctx *gin.Context) {
	// Get list of files
	files, err := a.leak.FilesList()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get list of files: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, "errPage.html", errors.New("Couldn't get list of files"))
	}
	// Execute template with list of files
	err = a.tmpl.ExecuteTemplate(ctx.Writer, "adminFilesList.html", files)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, "errPage.html", errHandler.ServerErr)
	}
}

// DeleteFile handles the file delete request
func (a *AdminHandler) DeleteFile(ctx *gin.Context) {
	// Get file id
	fileId := ctx.Param("fileId")
	// Delete file
	err := a.leak.DeleteFile(fileId)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't delete file: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	logg.Info(fmt.Sprintf("File %v was deleted by %v", fileId, ctx.MustGet("userId")))
	ctx.Redirect(200, "/admin/files")
}

// DownloadFile handles the file download request
func (a *AdminHandler) DownloadFile(ctx *gin.Context) {
	fileId := ctx.Param("fileId")
	fileData, err := a.leak.File(fileId)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get file: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", fileData.File.Name)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't create temporary file: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.FileReceiveErr)
		return
	}
	// Write the file content
	if _, err := tempFile.Write(fileData.File.Content); err != nil {
		logg.Error(fmt.Sprint("Couldn't write to temporary file: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.FileReceiveErr)
		return
	}
	fileData.File.Content = nil
	// Close the file
	defer func() {
		if err := tempFile.Close(); err != nil {
			logg.Error(fmt.Sprint("Couldn't close temporary file: ", err))
		}
		if err := os.Remove(tempFile.Name()); err != nil {
			logg.Error(fmt.Sprint("Couldn't remove temporary file: ", err))
		}
	}()
	logg.Info(fmt.Sprintf("File %v was downloaded by admin:%v", fileId, ctx.MustGet("userId")))
	// Set the appropriate headers
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileData.File.Name))
	ctx.File(tempFile.Name())
}

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
