package delivery

import (
	"errors"
	"log"
	"net/http"
	errHandler "uniLeaks/err"
	"uniLeaks/models"

	"github.com/gin-gonic/gin"
)

// AuthAndRefreshMiddleware checks for the presence of an auth token in cookies.
// If the auth token is missing, it tries to refresh the token using the refresh token.
// If the refresh token is missing, it deletes cookies and redirects to the login page. If the refresh token is present,
// It generates a new auth token and saves it. The user ID is then set in the context
// And the request is passed to the next middleware.
func (h *Handler) AuthAndRefreshMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the auth token from cookies.
		existedAuth, err := h.authTokenFromCookies(ctx)
		if err != nil {
			// If the auth token is missing, try to refresh the token using the refresh token.
			existedRefresh, err := h.refreshTokenFromCookies(ctx)
			if err != nil {
				log.Println("Middleware: ", err)
				// If the refresh token is missing, redirect to the login page.
				h.handleInvalidToken(http.StatusFound, ctx)
				return
			}

			userId, err := h.validateToken(existedRefresh)
			if err != nil {
				log.Println("Middleware: ", err)
				// If the refresh token is invalid, redirect to the login page.
				h.handleInvalidToken(http.StatusFound, ctx)
				return
			}

			// Generate a new auth token and save it.
			newAuthToken := models.Token{TokenType: AuthString, Value: h.generateToken(userId, AuthTokenDuration), Exp: AuthTokenDuration, UserId: userId}

			// Set the user ID in the context and cookies.
			ctx.Set("userId", userId)
			h.SetTokenToCookies(ctx, newAuthToken)
			ctx.Next()
			return
		}
		// If the auth token is present, validate it. If it is invalid, redirect to the login page. Otherwise, set the user ID in the context and proceed to the next middleware.
		userId, err := h.validateToken(existedAuth)
		if err != nil {
			log.Println("Middleware: ", err)
			h.handleInvalidToken(http.StatusFound, ctx)
			return
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}

// IsAdmin checks if the user is an admin.
func (h *Handler) IsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := ctx.Get("userId")
		if !ok {
			errHandler.ResponseWithErr(ctx, h.tmpl, errHandler.ErrPage, errors.New("Помилка отримання данних"))
			return
		}
		isAdmin, err := h.userService.IsAdmin(userId.(int))
		if err != nil {
			errHandler.ResponseWithErr(ctx, h.tmpl, errHandler.ErrPage, errors.New("Помилка отримання данних"))
			return
		}
		if !isAdmin {
			errHandler.ResponseWithErr(ctx, h.tmpl, errHandler.ErrPage, errors.New("Недостатньо прав"))
			return
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}
