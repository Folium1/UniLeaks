package delivery

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthorizeUser(c *gin.Context, userId int) {
	refreshToken := h.createRefreshToken(userId)
	authToken := h.createAuthToken(userId)

	h.SetTokenToCookies(c, refreshToken)
	h.SetTokenToCookies(c, authToken)
}

func (h *Handler) LogOut(c *gin.Context) {
	h.deleteCookies(c)
}
