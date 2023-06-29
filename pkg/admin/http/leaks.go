package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	errHandler "leaks/pkg/err"
	"os"

	"github.com/gin-gonic/gin"
)

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

// GetUserFilesList returns the list of files uploaded by user
func (a *AdminHandler) GetUserFilesList(ctx *gin.Context) {
	userId := ctx.PostForm("userId")
	if userId == "" {
		logg.Error(errHandler.NoUserIdErr.Error())
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.NoUserIdErr)
	}
	userFiles, err := a.leak.GetUserFilesList(userId)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get user files list: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	err = a.tmpl.ExecuteTemplate(ctx.Writer, "adminUserFilesList.html", userFiles)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.ServerErr)
	}
}

// DeleteUserFile deletes the file uploaded by user
func (a *AdminHandler) DeleteAllUserFiles(ctx *gin.Context) {
	userId := ctx.PostForm("userId")
	if userId == "" {
		logg.Error(errHandler.NoUserIdErr.Error())
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.NoUserIdErr)
	}
	err := a.leak.DeleteAllUserFiles(userId)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't delete user files: ", err))
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	logg.Info(fmt.Sprintf("All files of user %v were deleted by %v", userId, ctx.MustGet("userId")))
	ctx.Redirect(200, "/admin/")
}
