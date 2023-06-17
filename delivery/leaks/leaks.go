package leaks

import (
	"errors"
	"fmt"
	"io/ioutil"

	errHandler "leaks/err"
	"leaks/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	uploadPage = "uploadFilePage.html"
)

// MainPage displays the main page with the list of faculties
func (l *LeaksHandler) MainPage(ctx *gin.Context) {
	err := l.tmpl.ExecuteTemplate(ctx.Writer, "leaks.html", nil)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// UploadFilePage displays the upload file page
func (l *LeaksHandler) UploadFilePage(ctx *gin.Context) {
	err := l.tmpl.ExecuteTemplate(ctx.Writer, uploadPage, nil)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// UploadFile handles the file upload request
func (l *LeaksHandler) UploadFile(ctx *gin.Context) {
	// Get user id
	userId, ok := ctx.Get("userId")
	if !ok {
		logg.Info("Couldn't get user id from context")
		errHandler.ResponseWithErr(ctx, l.tmpl, uploadPage, errors.New("Помилка отримання данних з файлу"))
		return
	}
	// Parse file and subject data
	data := models.LeakData{File: &models.File{}, Subject: &models.SubjectData{}, UserData: &models.UserFileData{UserId: fmt.Sprintf("%v", userId)}}
	err := parseFileData(data.File, ctx)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't parse file data: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, uploadPage, errors.New("Помилка отримання данних з файлу"))
		return
	}
	err = parseSubjectData(data.Subject, ctx)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't parse subject data: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, uploadPage, errors.New("Помилка отримання заданих данних"))
		return
	}
	// Save file
	err = l.leakService.SaveFile(data)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't save file: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, uploadPage, err)
		return
	}
	logg.Info(fmt.Sprintf("File %s uploaded by %v", data.File.Name, userId))
	ctx.Redirect(http.StatusSeeOther, "/leaks/upload-files/")
}

// FilesPage displays the get files page
func (l *LeaksHandler) FilesPage(ctx *gin.Context) {
	err := l.tmpl.ExecuteTemplate(ctx.Writer, "getFiles.html", nil)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// FilesList retrieves and displays the list of files based on subject data
func (l *LeaksHandler) FilesList(ctx *gin.Context) {
	// Parse subject data
	var data models.SubjectData
	err := parseSubjectData(&data, ctx)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't parse subject data: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errors.New("Помилка при отриманні данних, спробуйте ще раз"))
		return
	}
	// Get files list
	files, err := l.leakService.FilesList(data)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get files list: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errors.New("Помилка отримання списку файлів"))
		return
	}
	// Display files list
	err = l.tmpl.ExecuteTemplate(ctx.Writer, "files.html", files)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// DownloadFile handles the file download request
func (l *LeaksHandler) DownloadFile(ctx *gin.Context) {
	fileId := ctx.Param("id")
	leakData, err := l.leakService.File(fileId)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get file: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.FileReceiveErr)
		return
	}
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", leakData.File.Id)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't create temporary file: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.FileReceiveErr)
		return
	}
	// Write the file content
	if _, err := tempFile.Write(leakData.File.Content); err != nil {
		logg.Error(fmt.Sprint("Couldn't write file content: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.FileReceiveErr)
		return
	}
	leakData.File.Content = nil
	// Close the file
	defer func() {
		if err := tempFile.Close(); err != nil {
			logg.Error(fmt.Sprint("Couldn't close file: ", err))
		}
		if err := os.Remove(tempFile.Name()); err != nil {
			logg.Error(fmt.Sprint("Couldn't remove file: ", err))
		}
	}()
	logg.Info(fmt.Sprintf("File downloaded:%v/%v/%v by %v", leakData.Subject.Faculty, leakData.Subject.Subject, leakData.File.Name, ctx.MustGet("userId")))
	// Set the appropriate headers
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", leakData.File.Name))
	ctx.File(tempFile.Name())
}

// AllFiles retrieves and displays all files
func (l *LeaksHandler) AllFiles(ctx *gin.Context) {
	files, err := l.leakService.AllFiles()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get files list: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errors.New("Помилка отримання списку файлів, спробуйте ще раз"))
		return
	}
	err = l.tmpl.ExecuteTemplate(ctx.Writer, "files.html", files)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// MyFiles retrieves and displays files uploaded by the user
func (l *LeaksHandler) MyFiles(ctx *gin.Context) {
	userId, ok := ctx.Get("userId")
	if !ok {
		logg.Error("Couldn't get userId from context")
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
	// Get files list
	files, err := l.leakService.MyFiles(fmt.Sprintf("%v", userId))
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get files list: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errors.New("Помилка отримання списку файлів, спробуйте ще раз"))
		return
	}
	// Execute template
	err = l.tmpl.ExecuteTemplate(ctx.Writer, "myFiles.html", files)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't execute template: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// LikeDislike handles the like/dislike request
func (l *LeaksHandler) LikeDislikeFile(ctx *gin.Context) {
	var data models.LikeDislikeData
	// Retrieve data from request
	err := ctx.BindJSON(&data)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't bind json: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
	// Get userId from context
	userId, ok := ctx.Get("userId")
	if !ok {
		logg.Error("Couldn't get userId from context")
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
	data.UserId = fmt.Sprintf("%v", userId)
	// Like/dislike file
	err = l.leakService.LikeDislikeFile(data)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't like/dislike: ", err))
		errHandler.ResponseWithErr(ctx, l.tmpl, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
	ctx.Status(http.StatusOK)
}
