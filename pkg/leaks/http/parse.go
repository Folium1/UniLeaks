package http

import (
	"bytes"
	"fmt"
	"io"
	"leaks/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseData(c *gin.Context) (models.Leak, error) {
	data := models.Leak{}
	data.User.UserId = c.MustGet("userId").(string)
	err := parseFile(data.File, c)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't parse file data : ", err))
		return data, err
	}
	err = parseSubject(data.Subject, c)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't parse subject data : ", err))
		return data, err
	}
	return data, nil
}

func parseFile(data *models.File, c *gin.Context) error {
	file, err := c.FormFile("file-upload")
	if err != nil {
		return err
	}
	data.Description = c.PostForm("file-description")
	data.Name = file.Filename
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	buffer := bytes.Buffer{}
	_, err = io.Copy(&buffer, f)
	if err != nil {
		return err
	}
	data.Content = buffer.Bytes()
	return nil
}

func parseSubject(data *models.Subject, c *gin.Context) error {
	var err error
	data.Faculty = c.PostForm("faculty")
	data.Subject = c.PostForm("subject")

	semesterNum := c.PostForm("semester")
	data.Semester, err = strconv.ParseUint(semesterNum, 10, 64)
	if err != nil {
		return err
	}

	moduleNumStr := c.PostForm("module_num")
	data.ModuleNum, err = strconv.ParseUint(moduleNumStr, 10, 64)
	if err != nil {
		return err
	}

	data.YearOfEducation = c.PostForm("edu_year")
	data.IsExam = c.PostForm("is_exam") == "on"
	data.IsModuleTask = c.PostForm("is_module") == "on"
	return nil
}
