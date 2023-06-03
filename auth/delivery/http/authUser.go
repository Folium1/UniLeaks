package delivery

import (
	"github.com/gin-gonic/gin"
)

// AuthorizeUser authorizes user by creating and saving tokens to cookies
func (h Handler) AuthorizeUser(c *gin.Context, userId int) {
	// Create and save tokens to cookies
	refreshToken := h.createRefreshToken(userId)
	authToken := h.createAuthToken(userId)

	// Set tokens to cookies
	h.SetTokenToCookies(c, refreshToken)
	h.SetTokenToCookies(c, authToken)
}

// LogOut logs out user by deleting cookies
func (h Handler) LogOut(c *gin.Context) {
	h.deleteCookies(c)
}
