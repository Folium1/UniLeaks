package delivery

import "github.com/gin-gonic/gin"

func (h Handler) economMainPage(c *gin.Context) {
	c.String(200, "Welcome to econom faculty")
}
