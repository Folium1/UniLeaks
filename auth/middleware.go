package delivery

import (
	"errors"
	"fmt"
	errHandler "leaks/err"
	"leaks/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthAndRefreshMiddleware checks for the presence of an auth token in cookies.
// If the auth token is missing, it tries to refresh the token using the refresh token.
// If the refresh token is missing, it deletes cookies and redirects to the login page.
// If the refresh token is present,it generates a new auth token and saves it. The user ID is then set in the context
// And the request is passed to the next middleware.
func (h *Handler) AuthAndRefreshMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the auth token from cookies.
		existedAccesToken, err := h.authTokenFromCookies(ctx)
		if err != nil {
			// If the auth token is missing, try to refresh the token using the refresh token.
			existedRefresh, err := h.refreshTokenFromCookies(ctx)
			if err != nil {
				logg.Error(fmt.Sprint("Middleware: ", err))
				// If the refresh token is missing, redirect to the login page.
				h.handleInvalidToken(http.StatusFound, ctx)
				return
			}

			userId, err := h.validateToken(existedRefresh)
			if err != nil {
				logg.Error(fmt.Sprint("Middleware: ", err))
				// If the refresh token is invalid, redirect to the login page.
				h.handleInvalidToken(http.StatusFound, ctx)
				return
			}

			// Generate a new auth token and save it.
			newAccesToken := models.Token{TokenType: AuthString, Value: h.generateToken(userId, AccesTokenDuration), Exp: AccesTokenDuration, UserId: userId}

			// Set the user ID in the context and cookies.
			ctx.Set("userId", userId)
			h.SetTokenToCookies(ctx, newAccesToken)
			ctx.Next()
			return
		}
		// If the auth token is present, validate it. If it is invalid, redirect to the login page. Otherwise, set the user ID in the context and proceed to the next middleware.
		userId, err := h.validateToken(existedAccesToken)
		if err != nil {
			logg.Error(fmt.Sprint("Middleware: ", err))
			h.handleInvalidToken(http.StatusFound, ctx)
			return
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}

// OnlyAdminMiddleware checks if the user is an admin.
func (h *Handler) OnlyAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// retreive user's id
		userId, ok := ctx.Get("userId")
		if !ok {
			logg.Error(fmt.Sprint("Middleware: ", "No user id in context"))
			errHandler.ResponseWithErr(ctx, h.tmpl, errHandler.ErrPage, errors.New("Помилка перевірки доступу"))
			return
		}
		// check if user is an admin
		isAdmin, err := h.userService.IsAdmin(userId.(int))
		if err != nil {
			logg.Error(fmt.Sprint("Middleware: ", err))
			errHandler.ResponseWithErr(ctx, h.tmpl, errHandler.ErrPage, errors.New("Помилка перевірки доступу"))
			return
		}
		if !isAdmin {
			logg.Error(fmt.Sprint("Middleware: ", "User is not admin"))
			errHandler.ResponseWithErr(ctx, h.tmpl, errHandler.ErrPage, errors.New("Недостатньо прав"))
			return
		}
		ctx.Set("userId", userId)
		ctx.Next()
	}
}
