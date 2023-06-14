package delivery

import (
	"errors"
	"fmt"
	"html/template"
	"leaks/logger"

	"github.com/gin-gonic/gin"
)

const ErrPage = "errPage.html"

var logg = logger.NewLogger()

// Define errors that can occur during the work of the server
var (
	UserIsBannedErr    = errors.New("Мейл заблокований")
	InvalidActionErr   = errors.New("Невірна дія")
	VirusDetectedErr   = errors.New("Знайдений вірус в файлі")
	FileNotFoundErr    = errors.New("Файл не знайдено")
	FileCheckErr       = errors.New("Сталась помилка, на стадії перевірки файлу на віруси")
	ServerErr          = errors.New("Помилка на сервері, спробуйте ще раз")
	FileReceiveErr     = errors.New("Помилка отримання файлу, спробуйте ще раз")
	FileListReceiveErr = errors.New("Помилка отримання списку файлів, спробуйте ще раз")
	BanUserErr         = errors.New("Помилка бану користувача, спробуйте ще раз")
	UserListReceiveErr = errors.New("Помилка отримання списку користувачів, спробуйте ще раз")
	FileSaveErr        = errors.New("Помилка збереження файлу, спробуйте ще раз")
	LikeDislikeErr     = errors.New("Помилка лайку/дизлайку, спробуйте ще раз")
)

// ResponseWithErr gives a response to user with code 400, in the given page and error message
func ResponseWithErr(c *gin.Context, t *template.Template, template string, err error) {
	err = t.ExecuteTemplate(c.Writer, template, err)
	if err != nil {
		logg.Error(fmt.Sprintf("Couldn't execute template: %v", err))
		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte("Internal server error"))
	}
}
