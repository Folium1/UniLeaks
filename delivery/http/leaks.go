package delivery

import (
	"errors"
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

var facultiesList = map[string]Faculties{
	"fim": {"Факультет іноземних мов", "fim"},
	"ff":  {"Філологічний факультет", "ff"},
	"eco": {"Економічний факультет", "eco"},
}

func (h Handler) mainPage(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, "leaks.html", facultiesList)
	if err != nil {
		responseWithErr(c, ErrPage, errors.New("There was an err while loading the page"))
		return
	}
}

func (h Handler) uploadFilePage(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, uploadPage, nil)
	if err != nil {
		responseWithErr(c, ErrPage, errors.New("There was an err while loading the page"))
		return
	}
}

func (h Handler) uploadFile(c *gin.Context) {
	data := models.LeakData{File: &models.File{}, Subject: &models.SubjectData{}, UserData: &models.UserFileData{}}
	err := parseFileData(&data, c)
	if err != nil {
		responseWithErr(c, uploadPage, errors.New("The error has been occured while getting data from the file"))
		return
	}
	err = h.leakService.SaveFile(data)
	if err != nil {
		responseWithErr(c, uploadPage, err)
		return
	}
	c.Redirect(200, "/leaks/")
}

func (h Handler) getFilesPage(c *gin.Context) {
	err := h.tmpl.ExecuteTemplate(c.Writer, "getFiles.html", nil)
	if err != nil {
		responseWithErr(c, ErrPage, errors.New("Internal server error"))
		return
	}
}

func (h Handler) getFilesList(c *gin.Context) {
	var data models.SubjectData
	err := parseSubjectData(&data, c)
	if err != nil {
		responseWithErr(c, ErrPage, errors.New("Помилка при отриманні заданих данних, спробуйте ще раз"))
		return
	}
	files, err := h.leakService.GetList(data)
	if err != nil {
		responseWithErr(c, ErrPage, errors.New("Internal server error"))
		return
	}
	err = h.tmpl.ExecuteTemplate(c.Writer, "fileList.html", files)
	if err != nil {
		responseWithErr(c, ErrPage, errors.New("Internal server error"))
		return
	}
}
