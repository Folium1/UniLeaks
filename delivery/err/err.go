package delivery

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ErrPage = "errPage.html"

var (
	ServerErr = errors.New("Помилка на сервері, спробуйте ще раз")
)

// responseWithErr gives a response to user with code 400, in the given page and error message
func ResponseWithErr(c *gin.Context, template string, err error) {
	c.HTML(http.StatusBadRequest, template,
		gin.H{"Error": err})
}
