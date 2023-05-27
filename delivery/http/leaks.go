package delivery

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
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
func (h *Handler) mainPage(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, "leaks.html", nil)
	if err != nil {
		responseWithErr(c, ErrPage, serverErr)
		return
	}
}

// uploadFilePage displays the upload file page
func (h *Handler) uploadFilePage(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, uploadPage, nil)
	if err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, serverErr)
		return
	}
}

// uploadFile handles the file upload request
func (h *Handler) uploadFile(c *gin.Context) {
	data := models.LeakData{File: &models.File{}, Subject: &models.SubjectData{}, UserData: &models.UserFileData{}}
	err := parseFileData(data.File, c)
	if err != nil {
		log.Println(err)
		responseWithErr(c, uploadPage, errors.New("Помилка отримання данних з файлу"))
		return
	}
	err = parseSubjectData(data.Subject, c)
	if err != nil {
		log.Println(err)
		responseWithErr(c, uploadPage, errors.New("Помилка отримання заданих данних"))
		return
	}
	err = h.leakService.SaveFile(&data)
	if err != nil {
		log.Println(err)
		responseWithErr(c, uploadPage, err)
		return
	}
	runtime.GC()
	data.File.Content = nil
	c.Redirect(StatusSeeOther, "/leaks/get-files/all/")
}

// filesPage displays the get files page
func (h *Handler) filesPage(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, "getFiles.html", nil)
	if err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, serverErr)
		return
	}
}

// filesList retrieves and displays the list of files based on subject data
func (h *Handler) filesList(c *gin.Context) {
	var data models.SubjectData
	err := parseSubjectData(&data, c)
	if err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, errors.New("Помилка при отриманні данних, спробуйте ще раз"))
		return
	}

	files, err := h.leakService.FilesList(data)
	if err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, errors.New("Помилка отримання списку файлів"))
		return
	}
	err = h.tmpl.ExecuteTemplate(c.Writer, "fileList.html", files)
	if err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, serverErr)
		return
	}
}

// file retrieves and serves the requested file
func (h *Handler) file(c *gin.Context) {
	fileId := c.Param("id")
	leakData, err := h.leakService.File(fileId)
	if err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, errors.New("Помилка отримання файлу, спробуйте ще раз"))
		return
	}
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", leakData.File.Name)
	if err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, errors.New("Помилка отримання файлу, спробуйте ще раз"))
		return
	}
	// Write the file content
	if _, err := tempFile.Write(leakData.File.Content); err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, errors.New("Помилка отримання файлу, спробуйте ще раз"))
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

	// Set the appropriate headers
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", leakData.File.Name))
	c.File(tempFile.Name())
}

// allFiles retrieves and displays all files
func (h *Handler) allFiles(c *gin.Context) {
	files, err := h.leakService.AllFiles()
	if err != nil {
		log.Println(err)
		responseWithErr(c, ErrPage, errors.New("Помилка отримання списку файлів, спробуйте ще раз"))
		return
	}
	err = h.tmpl.ExecuteTemplate(c.Writer, "fileList.html", files)
	if err != nil {
		log.Fatal(err)
	}
}
