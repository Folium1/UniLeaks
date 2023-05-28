package delivery

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthorizeUser authorizes user by creating and saving tokens to cookies
func (h Handler) AuthorizeUser(c *gin.Context, userId int) {
	// Create and save tokens to cookies
	refreshToken := h.createRefreshToken(userId)
	authToken := h.createAuthToken(userId)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// Save tokens to database
	err := h.useCase.SaveToken(ctx, refreshToken)
	if err != nil {
		log.Println(err)
	}
	err = h.useCase.SaveToken(ctx, authToken)
	if err != nil {
		log.Println(err)
	}
	// Set tokens to cookies
	h.SetTokenToCookies(c, refreshToken)
	h.SetTokenToCookies(c, authToken)
}

// LogOut logs out user by deleting cookies
func (h Handler) LogOut(c *gin.Context) {
	h.deleteCookies(c)
}
