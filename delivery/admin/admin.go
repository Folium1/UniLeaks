package admin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	errHandler "uniLeaks/err"

	"github.com/gin-gonic/gin"
)

func (a *AdminHandler) MainPage(ctx *gin.Context) {
	a.tmpl.ExecuteTemplate(ctx.Writer, "admin.html", nil)
}

func (a *AdminHandler) FilesList(ctx *gin.Context) {
	files, err := a.leak.FilesList()
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(ctx, a.tmpl, "errPage.html", errors.New("Couldn't get list of files"))
	}
	err = a.tmpl.ExecuteTemplate(ctx.Writer, "adminFilesList.html", files)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(ctx, a.tmpl, "errPage.html", errHandler.ServerErr)
	}
}

func (a *AdminHandler) DeleteFile(ctx *gin.Context) {
	fileId := ctx.Param("fileId")
	err := a.leak.DeleteFile(fileId)
	if err != nil {
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	ctx.Redirect(200, "/admin/files")
}

func (a *AdminHandler) DownloadFile(ctx *gin.Context) {
	fileId := ctx.Param("fileId")
	fileData, err := a.leak.File(fileId)
	if err != nil {
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, err)
	}
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", fileData.File.Name)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.FileReceivedErr)
		return
	}
	// Write the file content
	if _, err := tempFile.Write(fileData.File.Content); err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(ctx, a.tmpl, errHandler.ErrPage, errHandler.FileReceivedErr)
		return
	}
	fileData.File.Content = nil
	// Close the file
	defer func() {
		if err := tempFile.Close(); err != nil {
			log.Println(err)
		}
		if err := os.Remove(tempFile.Name()); err != nil {
			log.Println(err)
		}
	}()
	// Set the appropriate headers
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileData.File.Name))
	ctx.File(tempFile.Name())
}
