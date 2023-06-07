package delivery

import (
	"errors"
	"html/template"

	"github.com/gin-gonic/gin"
)

const ErrPage = "errPage.html"

var (
	ErrVirusDetected = errors.New("Знайдений вірус в файлі")
	ErrFileNotFound     = errors.New("Файл не знайдено")
	ErrFileCheck     = errors.New("Сталась помилка, на стадії перевірки файлу на віруси")
	ServerErr       = errors.New("Помилка на сервері, спробуйте ще раз")
	FileReceivedErr = errors.New("Помилка отримання файлу, спробуйте ще раз")
)

// ResponseWithErr gives a response to user with code 400, in the given page and error message
func ResponseWithErr(c *gin.Context, t *template.Template, template string, err error) {
	t.ExecuteTemplate(c.Writer, template, err)
}
