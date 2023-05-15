package delivery

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"uniLeaks/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

// Define constants for the cookie names
const (
	RefreshString = "Refresh"
	AuthtString   = "Authorization"
)

// Define global variables for token durations
var (
	RefreshTokenDuration = ((60 * 60) * 24) * 30 // 30 days in seconds
	AuthTokenDuration    = 15 * 60               // 15 minutes in seconds
)

// deleteCookies removes the authentication and refresh tokens from cookies.
func (h Handler) deleteCookies(c *gin.Context) {
	c.SetCookie(AuthtString, "", -1, "", "", false, false)
	c.SetCookie(RefreshString, "", -1, "", "", false, false)
}

// createAuthToken generates an authentication token for a given user ID.
func (h Handler) createAuthToken(userId int) models.Token {
	authToken := models.Token{
		TokenType: AuthtString,
		Exp:       AuthTokenDuration,
		UserId:    userId,
	}
	authToken.Value = h.generateToken(userId, authToken.Exp)

	return authToken
}

// createRefreshToken generates a refresh token for a given user ID.
func (h Handler) createRefreshToken(userId int) models.Token {
	refreshToken := models.Token{
		TokenType: RefreshString,
		Exp:       RefreshTokenDuration,
		UserId:    userId,
	}
	refreshToken.Value = h.generateToken(userId, refreshToken.Exp)
	return refreshToken
}

// generateToken generates a JWT token for a given user ID and expiration time.
func (h Handler) generateToken(userId int, expire int) string {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["Exp"] = expire

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("SALT"))
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return ""
	}
	return signedToken
}

// getTokenFromCookies gets the authentication or refresh token from cookies.
func (h Handler) getTokenFromCookies(c *gin.Context) (models.Token, error) {
	token, err := c.Cookie(AuthtString)
	if err != nil {
		refreshToken, err := c.Cookie(RefreshString)
		if err != nil {
			return models.Token{}, err
		}
		token := models.Token{TokenType: RefreshString, Value: refreshToken}

		return token, nil
	}
	tk := models.Token{TokenType: AuthtString, Value: token}
	return tk, nil
}

// getAuthTokenFromCookies retrieves the authentication token from cookies in the context.
// If there is an error retrieving the token or the token is empty, it returns an error.
// Otherwise, it splits the token to get the value and returns a Token object with TokenType "AuthtString".
func (h Handler) getAuthTokenFromCookies(c *gin.Context) (models.Token, error) {
	cookieToken, err := c.Cookie(AuthtString)
	if err != nil || cookieToken == "" {
		return models.Token{}, err
	}
	tokenValue := strings.Split(cookieToken, " ")[1]
	authToken := models.Token{TokenType: AuthtString, Value: tokenValue}
	return authToken, nil
}

// getRefreshTokenFromCookies retrieves the refresh token from cookies in the context.
// If there is an error retrieving the token or the token is empty, it returns an error.
// Otherwise, it splits the token to get the value and returns a Token object with TokenType "RefreshString".
func (h Handler) getRefreshTokenFromCookies(c *gin.Context) (models.Token, error) {
	cookieToken, err := c.Cookie(RefreshString)
	if err != nil || cookieToken == "" {
		return models.Token{}, err
	}
	tokenValue := strings.Split(cookieToken, " ")[1]
	authToken := models.Token{TokenType: RefreshString, Value: tokenValue}
	return authToken, nil
}

// getUserId gets the user ID from the given token using the GetUserId function.
// It returns the user ID and any error encountered during the retrieval process.
func (h Handler) getUserId(tk models.Token) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userId, err := h.useCase.GetUserId(ctx, tk)
	if err != nil {
		return -1, err
	}
	return userId, nil
}

// SetTokenToCookies sets given token to cookies in a Gin context
func (h Handler) SetTokenToCookies(c *gin.Context, token models.Token) {
	c.SetCookie(token.TokenType, fmt.Sprintf("Bearer "+token.Value), int(token.Exp), "", "", true, true)
}

// handleInvalidToken handles middleware errors by deleting cookies and redirecting to login page
func (h Handler) handleInvalidToken(code int, c *gin.Context) {
	h.deleteCookies(c)
	c.Redirect(code, "/user/login")
}
