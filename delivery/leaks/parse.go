package leaks

import (
	"bytes"
	"io"
	"leaks/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseFileData parses the file data from the request
func parseFileData(data *models.File, c *gin.Context) error {
	file, err := c.FormFile("file-upload")
	if err != nil {
		return err
	}
	fileDescription := c.PostForm("file-description")
	data.Description = fileDescription
	data.Name = file.Filename
	// Open the file
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a buffer to store the file content
	buffer := &bytes.Buffer{}

	// Copy the file content to the buffer
	_, err = io.Copy(buffer, f)
	if err != nil {
		return err
	}
	// Set the buffer bytes as the file content
	data.Content = buffer.Bytes()
	return nil
}

// parseSubjectData parses the subject data from the request
func parseSubjectData(data *models.SubjectData, c *gin.Context) error {
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
