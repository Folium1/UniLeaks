package delivery

import (
	"github.com/gin-gonic/gin"
)

func (h Handler) AuthorizeUser(c *gin.Context, userId int) {
	refrToken := h.createRefreshToken(userId)
	authToken := h.createAuthToken(userId)

	h.SetTokenToCookies(c, refrToken)
	h.SetTokenToCookies(c, authToken)
}
