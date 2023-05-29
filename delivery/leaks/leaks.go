package leaks

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	errHandler "uniLeaks/delivery/err"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
)

const (
	uploadPage = "uploadFilePage.html"
)

type Faculties struct {
	Name  string
	Value string
}

// mainPage displays the main page with the list of faculties
func (l *LeaksHandler) MainPage(c *gin.Context) {
	err := l.tmpl.ExecuteTemplate(c.Writer, "leaks.html", nil)
	if err != nil {
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// uploadFilePage displays the upload file page
func (l *LeaksHandler) UploadFilePage(c *gin.Context) {
	err := l.tmpl.ExecuteTemplate(c.Writer, uploadPage, nil)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// uploadFile handles the file upload request
func (l *LeaksHandler) UploadFile(c *gin.Context) {
	data := models.LeakData{File: &models.File{}, Subject: &models.SubjectData{}, UserData: &models.UserFileData{}}
	err := parseFileData(data.File, c)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, uploadPage, errors.New("Помилка отримання данних з файлу"))
		return
	}
	err = parseSubjectData(data.Subject, c)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, uploadPage, errors.New("Помилка отримання заданих данних"))
		return
	}
	err = l.leakService.SaveFile(&data)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, uploadPage, err)
		return
	}
	runtime.GC()
	data.File.Content = nil
	c.Redirect(http.StatusSeeOther, "/leaks/get-files/all/")
}

// filesPage displays the get files page
func (l *LeaksHandler) FilesPage(c *gin.Context) {
	err := l.tmpl.ExecuteTemplate(c.Writer, "getFiles.html", nil)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// filesList retrieves and displays the list of files based on subject data
func (l *LeaksHandler) FilesList(c *gin.Context) {
	var data models.SubjectData
	err := parseSubjectData(&data, c)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errors.New("Помилка при отриманні данних, спробуйте ще раз"))
		return
	}

	files, err := l.leakService.FilesList(data)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errors.New("Помилка отримання списку файлів"))
		return
	}
	err = l.tmpl.ExecuteTemplate(c.Writer, "fileList.html", files)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errHandler.ServerErr)
		return
	}
}

// file retrieves and serves the requested file
func (l *LeaksHandler) DownloadFile(c *gin.Context) {
	fileId := c.Param("id")
	leakData, err := l.leakService.File(fileId)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errors.New("Помилка отримання файлу, спробуйте ще раз"))
		return
	}
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", leakData.File.Name)
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errors.New("Помилка отримання файлу, спробуйте ще раз"))
		return
	}
	// Write the file content
	if _, err := tempFile.Write(leakData.File.Content); err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errors.New("Помилка отримання файлу, спробуйте ще раз"))
		return
	}
	leakData.File.Content = nil
	// Close the file
	defer func() {
		if err := tempFile.Close(); err != nil {
			log.Println(err)
		}
		if err := os.Remove(tempFile.Name()); err != nil {
			log.Println(err)
		}
	}()
	runtime.GC()
	// Set the appropriate headers
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", leakData.File.Name))
	c.File(tempFile.Name())
}

// allFiles retrieves and displays all files
func (l *LeaksHandler) AllFiles(c *gin.Context) {
	files, err := l.leakService.AllFiles()
	if err != nil {
		log.Println(err)
		errHandler.ResponseWithErr(c, errHandler.ErrPage, errors.New("Помилка отримання списку файлів, спробуйте ще раз"))
		return
	}
	err = l.tmpl.ExecuteTemplate(c.Writer, "fileList.html", files)
	if err != nil {
		log.Fatal(err)
	}
}
