package delivery

import (
	"fmt"
	"net/http"

	"leaks/pkg/models"

	"github.com/gin-gonic/gin"
)

// AuthAndRefreshMiddleware checks for the presence of an auth token in cookies.Sets the user ID in the context and cookies.If the auth token is missing, try to refresh the token using refresh token, if it is missing, redirecting user to login page
func (h *Handler) AuthAndRefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		existedAccessToken, err := h.authTokenFromCookies(c)
		if err != nil {

			existedRefreshToken, err := h.refreshTokenFromCookies(c)
			if err != nil {
				logger.Error(fmt.Sprint("Middleware: ", err))
				h.handleInvalidToken(http.StatusUnauthorized, c)
				return
			}
			userId, err := h.validateToken(existedRefreshToken)
			if err != nil {
				logger.Error(fmt.Sprint("Middleware: ", err))
				h.handleInvalidToken(http.StatusUnauthorized, c)
				return
			}
			newAccessToken := models.Token{
				TokenType: AuthString,
				Value:     h.generateToken(userId, AccesTokenDuration),
				Exp:       AccesTokenDuration, UserId: userId,
			}

			c.Set("userId", userId)
			h.SetTokenToCookies(c, newAccessToken)
			c.Next()
			return
		}
		userId, err := h.validateToken(existedAccessToken)
		if err != nil {
			logger.Error(fmt.Sprint("Middleware: ", err))
			h.handleInvalidToken(http.StatusFound, c)
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}

// OnlyAdminMiddleware checks if the user is an admin. Sets the user ID,and user nick in the context.
func (h *Handler) OnlyAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("userId")
		if !ok {
			logger.Error(fmt.Sprint("Middleware: ", "No user id in context"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка перевірки доступу"})
			return
		}
		user, err := h.userService.IsAdmin(userId.(int))
		if err != nil {
			logger.Error(fmt.Sprint("Middleware: ", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка перевірки доступу"})
			return
		}
		if !user.IsAdmin {
			logger.Error(fmt.Sprint("Middleware: ", "User is not admin"))
			c.JSON(http.StatusForbidden, gin.H{"error": "Недостатньо прав"})
			return
		}
		logger.Info(fmt.Sprintf("Admin %s(%d) is logged in", user.NickName, user.ID))
		c.Set("userId", userId)
		c.Set("userName", user.NickName)
		c.Next()
	}
}
