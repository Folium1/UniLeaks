package delivery

import (
	"fmt"
	"leaks/models"
	"log"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

// Define constants for the cookie names
const (
	RefreshString = "Refresh"
	AuthString    = "Authorization"
)

// Define global variables for token durations
var (
	RefreshTokenDuration = ((60 * 60) * 24) * 30 // 30 days in seconds
	AccesTokenDuration   = 15 * 60               // 15 minutes in seconds

)

// deleteCookies removes the authentication and refresh tokens from cookies.
func (h *Handler) deleteCookies(c *gin.Context) {
	c.SetCookie(AuthString, "", -1, "", "", false, false)
	c.SetCookie(RefreshString, "", -1, "", "", false, false)
}

// createAuthToken generates an authentication token for a given user ID.
func (h *Handler) createAuthToken(userId int) models.Token {
	authToken := models.Token{
		TokenType: AuthString,
		Exp:       AccesTokenDuration,
		UserId:    userId,
	}
	authToken.Value = h.generateToken(userId, authToken.Exp)

	return authToken
}

// createRefreshToken generates a refresh token for a given user ID.
func (h *Handler) createRefreshToken(userId int) models.Token {
	refreshToken := models.Token{
		TokenType: RefreshString,
		Exp:       RefreshTokenDuration,
		UserId:    userId,
	}
	refreshToken.Value = h.generateToken(userId, refreshToken.Exp)
	return refreshToken
}

// generateToken generates a JWT token for a given user ID and expiration time.
func (h *Handler) generateToken(userId int, expire int) string {
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

// AuthTokenFromCookies retrieves the authentication token from cookies in the context.
// If there is an error retrieving the token or the token is empty, it returns an error.
// Otherwise, it splits the token to get the value and returns a Token object with TokenType "AuthString".
func (h *Handler) authTokenFromCookies(c *gin.Context) (string, error) {
	cookieToken, err := c.Cookie(AuthString)
	if err != nil || cookieToken == "" {
		return "", err
	}
	tokenValue := strings.Split(cookieToken, " ")[1]
	return tokenValue, nil
}

// RefreshTokenFromCookies retrieves the refresh token from cookies in the context.
// If there is an error retrieving the token or the token is empty, it returns an error.
// Otherwise, it splits the token to get the value and returns a Token object with TokenType "RefreshString".
func (h *Handler) refreshTokenFromCookies(c *gin.Context) (string, error) {
	cookieToken, err := c.Cookie(RefreshString)
	if err != nil || cookieToken == "" {
		return "", err
	}
	tokenValue := strings.Split(cookieToken, " ")[1]
	return tokenValue, nil
}

// validateToken validates a given token string. Rerturns the user ID if the token is valid.
func (h *Handler) validateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SALT")), nil
	})
	if err != nil || !token.Valid {
		log.Println(err)
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))
	return userId, nil
}

// SetTokenToCookies sets given token to cookies in a Gin context
func (h *Handler) SetTokenToCookies(c *gin.Context, token models.Token) {
	c.SetCookie(token.TokenType, fmt.Sprintf("Bearer "+token.Value), int(token.Exp), "", "", true, true)
}

// handleInvalidToken handles middleware errors by deleting cookies and redirecting to login page
func (h *Handler) handleInvalidToken(code int, c *gin.Context) {
	h.deleteCookies(c)
	c.Redirect(code, "/user/login")
}
