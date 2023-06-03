package delivery

import (
	"log"
	"net/http"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
)

// AuthAndRefreshMiddleware checks for the presence of an auth token in cookies.
// If the auth token is missing, it tries to refresh the token using the refresh token.
// If the refresh token is missing, it deletes cookies and redirects to the login page. If the refresh token is present,
// It generates a new auth token and saves it. The user ID is then set in the context
// And the request is passed to the next middleware.
func (h Handler) AuthAndRefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the auth token from cookies.
		existedAuth, err := h.authTokenFromCookies(c)
		if err != nil {
			// If the auth token is missing, try to refresh the token using the refresh token.
			existedRefresh, err := h.refreshTokenFromCookies(c)
			if err != nil {
				log.Println("Middleware: ", err)
				// If the refresh token is missing, redirect to the login page.
				h.handleInvalidToken(http.StatusFound, c)
				return
			}

			userId, err := h.validateToken(existedRefresh)
			if err != nil {
				log.Println("Middleware: ", err)
				// If the refresh token is invalid, redirect to the login page.
				h.handleInvalidToken(http.StatusFound, c)
				return
			}

			// Generate a new auth token and save it.
			newAuthToken := models.Token{TokenType: AuthString, Value: h.generateToken(userId, AuthTokenDuration), Exp: AuthTokenDuration, UserId: userId}

			// Set the user ID in the context and cookies.
			c.Set("userId", userId)
			h.SetTokenToCookies(c, newAuthToken)
			c.Next()
			return
		}
		// If the auth token is present, validate it. If it is invalid, redirect to the login page. Otherwise, set the user ID in the context and proceed to the next middleware.
		userId, err := h.validateToken(existedAuth)
		if err != nil {
			log.Println("Middleware: ", err)
			h.handleInvalidToken(http.StatusFound, c)
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
