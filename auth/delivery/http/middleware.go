package delivery

import (
	"context"
	"log"
	"net/http"
	"time"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
)

// AuthAndRefreshMiddleware checks for the presence of an auth token in cookies.
// If the auth token is missing, it tries to refresh the token using the refresh token.
// If the refresh token is missing, it deletes cookies and redirects to the login page. If the refresh token is present,
// It generates a new auth token and saves it. The user ID is then set in the context and cookies
// And the request is passed to the next middleware.
func (h Handler) AuthAndRefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout of 5 seconds.
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
		defer cancel()

		// Get the auth token from cookies.
		existedAuth, err := h.AuthTokenFromCookies(c)
		if err != nil {
			// If the auth token is missing, try to refresh the token using the refresh token.
			existedRefresh, err := h.RefreshTokenFromCookies(c)
			if err != nil {
				// If the refresh token is missing, redirect to the login page.
				h.handleInvalidToken(http.StatusFound, c)
				return
			}

			// Get the user ID from the refresh token.
			userId, err := h.useCase.UserId(ctx, existedRefresh)
			if err != nil {
				log.Println("Middleware: ", err)
				h.handleInvalidToken(http.StatusFound, c)
				return
			}

			// Generate a new auth token and save it.
			newAuthToken := models.Token{TokenType: AuthString, Value: h.generateToken(userId, AuthTokenDuration), Exp: AuthTokenDuration, UserId: userId}
			err = h.useCase.SaveToken(ctx, newAuthToken)
			if err != nil {
				log.Println("Middleware: ", err)
				h.handleInvalidToken(http.StatusFound, c)
				return
			}

			// Set the user ID in the context and cookies.
			c.Set("userId", userId)
			h.SetTokenToCookies(c, newAuthToken)
			c.Next()
			return
		}

		// Get the user ID from the auth token.
		userId, err := h.useCase.UserId(ctx, existedAuth)
		if err != nil {
			log.Println("Middleware: ", err)
			h.handleInvalidToken(http.StatusFound, c)
			return
		}

		// Set the user ID in the context and proceed to the next middleware.
		c.Set("userId", userId)
		c.Next()
	}
}
