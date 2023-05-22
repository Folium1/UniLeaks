package delivery

import (
	"log"
	"os"
	"strconv"

	// authHttp "uniLeaks/auth/delivery/http"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
)

func parseFileData(data *models.LeakData, c *gin.Context) error {
	// userId := c.MustGet("userId")
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	data.File.Description = c.PostForm("file-description")
	data.File.OpenedFile, err = os.Open(file.Filename)
	if err != nil {
		return err
	}
	return nil
}

func parseSubjectData(data *models.SubjectData, c *gin.Context) error {
	var err error
	data.Faculty = c.PostForm("faculty")
	data.Subject = c.PostForm("subject")
	data.ModuleNum, err = strconv.ParseUint(c.PostForm("module_num"), 0, 0)
	if err != nil {
		log.Println("Err to parse to uint, err:", err)
		return err
	}
	data.YearOfEducation = c.PostForm("edu_year")
	data.IsExam = c.PostForm("is_exam") == "on"
	data.IsModuleTask = c.PostForm("is_module") == "on"
	return nil
}
