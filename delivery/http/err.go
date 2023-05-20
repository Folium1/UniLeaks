package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	NotFound            = 404
	BadRequest          = 400
	Unauthorized        = 401
	InternalServerError = 500
	StatusSeeOther      = 303
	StatusFound         = 302
)

var ErrPage = "errPage.html"

func responseWithErr(c *gin.Context, template string, err error) {
	c.HTML(http.StatusBadRequest, template,
		gin.H{"Error": err})
}
