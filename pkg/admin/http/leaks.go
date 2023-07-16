package http

import (
	"fmt"
	"io/ioutil"
	errHandler "leaks/pkg/err"
	"os"

	"github.com/gin-gonic/gin"
)

func (a *AdminHandler) FilesOrderedByDislikes(c *gin.Context) {
	a.leak.FilesOrderedByDislikes()
	files, err := a.leak.FilesOrderedByDislikes()
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get list of files: ", err))
		c.JSON(400, gin.H{"error": "Couldn't get list of files"})
	}
	c.JSON(200, gin.H{"files": files})
}

func (a *AdminHandler) DeleteFile(c *gin.Context) {
	fileId := c.Param("fileId")
	err := a.leak.DeleteFile(fileId)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't delete file: ", err))
		c.JSON(400, gin.H{"error": "Couldn't delete file"})
	}
	l.Info(fmt.Sprintf("File %v was deleted by %v", fileId, c.MustGet("adminId")))
	c.Redirect(200, "/admin/files")
}

func (a *AdminHandler) DownloadFile(c *gin.Context) {
	fileData, err := a.leak.File(c.Param("fileId"))
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get file: ", err))
		c.JSON(400, gin.H{"error": "Couldn't get file"})
	}

	tempFile, err := ioutil.TempFile("", fileData.File.Name)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't create temporary file: ", err))
		c.JSON(400, gin.H{"error": "Error while receiving file"})
		return
	}
	if _, err := tempFile.Write(fileData.File.Content); err != nil {
		l.Error(fmt.Sprint("Couldn't write to temporary file: ", err))
		c.JSON(400, gin.H{"error": "Error while receiving file"})
		return
	}
	defer func() {
		if err := tempFile.Close(); err != nil {
			l.Error(fmt.Sprint("Couldn't close temporary file: ", err))
		}
		if err := os.Remove(tempFile.Name()); err != nil {
			l.Error(fmt.Sprint("Couldn't remove temporary file: ", err))
		}
	}()

	l.Info(fmt.Sprintf("File %v was downloaded by admin:%v", c.Param("fileId"), c.MustGet("adminId")))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileData.File.Name))
	c.File(tempFile.Name())
	c.Status(200)
}

func (a *AdminHandler) GetAllUserFiles(c *gin.Context) {
	userId := c.PostForm("userId")
	if userId == "" {
		l.Error(errHandler.NoUserIdErr.Error())
		c.JSON(400, gin.H{"error": errHandler.NoUserIdErr.Error()})
	}
	userFiles, err := a.leak.GetFilesListUploadedByUser(userId)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get user files list: ", err))
		c.JSON(400, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"files": userFiles})
}

func (a *AdminHandler) DeleteAllUserFiles(c *gin.Context) {
	userId := c.PostForm("userId")
	if userId == "" {
		l.Error(errHandler.NoUserIdErr.Error())
		c.JSON(400, gin.H{"error": errHandler.NoUserIdErr.Error()})
	}
	err := a.leak.DeleteAllFilesUploadedByUser(userId)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't delete user files: ", err))
		c.JSON(400, gin.H{"error": err.Error()})
	}

	l.Info(fmt.Sprintf("All files of user %v were deleted by %v(%v)", userId, c.MustGet("adminName"), c.MustGet("adminId")))
	c.Redirect(200, "/admin/")
}
