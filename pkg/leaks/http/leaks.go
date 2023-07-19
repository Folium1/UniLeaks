package http

import (
	"errors"
	"fmt"

	"leaks/pkg/models"

	"github.com/gin-gonic/gin"
)

func (l *LeaksHandler) MainPage(c *gin.Context) {
	// If the control reaches this point, it means the user is authenticated.
	c.Status(200)
}

func (l *LeaksHandler) UploadFilePage(c *gin.Context) {
	// If the control reaches this point, it means the user is authenticated.
	c.Status(200)
}

func (l *LeaksHandler) FilesPage(c *gin.Context) {
	// If the control reaches this point, it means the user is authenticated.
	c.Status(200)
}

func (l *LeaksHandler) UploadFile(c *gin.Context) {
	leaksData, err := parseData(c)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't parse data : ", err))
		c.JSON(400, gin.H{"error": errors.New("Помилка отримання заданих данних")})
		return
	}
	err = l.leakService.SaveFile(leaksData)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't save file: ", err))
		c.JSON(500, gin.H{"error": errors.New("Помилка збереження файлу")})
		return
	}
	logger.Info(fmt.Sprintf("File %s uploaded by %v", leaksData.File.Name, c.MustGet("userId")))
	c.Redirect(303, "/leaks/upload-files/")
}

func (l *LeaksHandler) FilesList(c *gin.Context) {
	var data models.Subject
	err := parseSubject(&data, c)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't parse subject : ", err))
		c.JSON(400, gin.H{"error": errors.New("Помилка отримання заданих данних")})
		return
	}
	files, err := l.leakService.FilesList(data)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get files list: ", err))
		c.JSON(500, gin.H{"error": errors.New("Помилка отримання файлів")})
		return
	}
	c.JSON(200, files)
}

func (l *LeaksHandler) DownloadFile(c *gin.Context) {
	leak, err := l.leakService.File(c.Param("id"))
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get file: ", err))
		c.JSON(500, gin.H{"error": errors.New("Помилка отримання файлу")})
		return
	}
	tempFile, err := createTempFile(leak)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't create temporary file: ", err))
		c.JSON(500, gin.H{"error": errors.New("Помилка отримання файлу")})
	}
	defer deleteTempFile(tempFile)

	logger.Info(fmt.Sprintf("File downloaded:%v/%v/%v by %v", leak.Subject.Faculty, leak.Subject.Subject, leak.File.Name, c.MustGet("userId")))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", leak.File.Name))
	c.File(tempFile.Name())
}

func (l *LeaksHandler) AllFiles(c *gin.Context) {
	files, err := l.leakService.AllFiles()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get files list: ", err))
		c.JSON(500, gin.H{"error": errors.New("Помилка отримання файлів")})
		return
	}
	c.JSON(200, files)
}

func (l *LeaksHandler) MyFiles(c *gin.Context) {
	files, err := l.leakService.UserFiles(c.MustGet("userId").(string))
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get files list: ", err))
		c.JSON(500, gin.H{"error": errors.New("Помилка отримання файлів")})
		return
	}
	c.JSON(200, files)
}

func (l *LeaksHandler) LikeDislikeFile(c *gin.Context) {
	data := models.LikeDislike{UserId: c.MustGet("userId").(string)}
	err := c.BindJSON(&data)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't bind json: ", err))
		c.JSON(400, gin.H{"error": errors.New("Помилка отримання заданих данних")})
		return
	}
	err = l.leakService.LikeDislikeFile(data)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't like/dislike: ", err))
		c.JSON(500, gin.H{"error": errors.New("Помилка при відправці оцінки")})
		return
	}
	c.Status(200)
}
