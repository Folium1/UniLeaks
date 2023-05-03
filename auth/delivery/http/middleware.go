package delivery

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthAndRefreshMiddleware is a Gin middleware function that handles authentication and token refreshing.
func (h Handler) AuthAndRefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from cookies, and redirect to the login page if it doesn't exist or is invalid.
		existedToken, err := h.getTokenFromCookies(c)
		if err != nil {
			h.deleteCookies(c)
			c.Redirect(http.StatusFound, "/user/login")
			return
		}

		// If the token is a refresh token, validate it and create a new access token.
		if existedToken.TokenType == RefreshString {
			res := strings.Split(existedToken.Tk, " ")
			if len(res) != 2 || res[0] != "Bearer" {
				h.deleteCookies(c)
				c.Redirect(http.StatusFound, "/user/login")
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			userId, err := h.useCase.GetUserId(ctx, existedToken)
			if err != nil {
				h.deleteCookies(c)
				c.Redirect(http.StatusFound, "/user/login")
				return
			}
			authToken := h.createAuthToken(userId)
			h.SetTokenToCookies(c, authToken)
		}

		// Proceed to the next middleware.
		c.Next()
	}
}
