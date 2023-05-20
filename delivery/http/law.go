package delivery

import "github.com/gin-gonic/gin"

var lawSubjectsMap = map[string]subjects{
	
}

func (h Handler) lawMainPage(c *gin.Context) {
	c.String(200, "Welcome to law faculty")
}

func (h Handler) lawHandleSubj(c *gin.Context) {

}
