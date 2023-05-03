package delivery

import (
	"fmt"
	"os"
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
	JwtTokenDuration     = 15 * 60               // 15 minutes in seconds
)

// Delete cookies by setting them to expired
func (h Handler) deleteCookies(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, false)
	c.SetCookie("Refresh", "", -1, "", "", false, false)
}

// Create an authentication token for a user
func (h Handler) createAuthToken(userId int) models.Token {
	var authToken models.Token
	authToken.TokenType = AuthtString
	authToken.Exp = JwtTokenDuration
	authToken.Tk = h.generateToken(userId, authToken.Exp)
	return authToken
}

// Create a refresh token for a user
func (h Handler) createRefreshToken(userId int) models.Token {
	var refreshToken models.Token
	refreshToken.TokenType = RefreshString
	refreshToken.Exp = RefreshTokenDuration
	refreshToken.Tk = h.generateToken(userId, refreshToken.Exp)
	return refreshToken
}

// Generate a JWT token using user ID and expiration time
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

// Get token from cookies in a Gin context
func (h Handler) getTokenFromCookies(c *gin.Context) (models.Token, error) {
	token, err := c.Cookie(AuthtString)
	if err != nil {
		refreshToken, err := c.Cookie(RefreshString)
		if err != nil {
			return models.Token{}, err
		}
		token := models.Token{TokenType: RefreshString, Tk: refreshToken}

		return token, nil
	}
	tk := models.Token{TokenType: "Authorization", Tk: token}
	return tk, nil
}

// Set token to cookies in a Gin context
func (h Handler) SetTokenToCookies(c *gin.Context, token models.Token) {
	c.SetCookie(token.TokenType, fmt.Sprintf("Bearer "+token.Tk), int(token.Exp), "", "", false, true)
}
