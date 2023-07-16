package delivery

import (
	"fmt"
	"log"
	"os"
	"strings"

	"leaks/pkg/models"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

// Define constants for the cookie names
const (
	RefreshString = "Refresh"
	AuthString    = "Authorization"
)

var (
	RefreshTokenDuration = ((60 * 60) * 24) * 30 // 30 days in seconds
	AccesTokenDuration   = 15 * 60               // 15 minutes in seconds

)

func (h *Handler) deleteCookies(c *gin.Context) {
	c.SetCookie(AuthString, "", -1, "", "", false, false)
	c.SetCookie(RefreshString, "", -1, "", "", false, false)
}

func (h *Handler) createAuthToken(userId int) models.Token {
	authToken := models.Token{
		TokenType: AuthString,
		Exp:       AccesTokenDuration,
		UserId:    userId,
	}
	authToken.Value = h.generateToken(userId, authToken.Exp)

	return authToken
}

func (h *Handler) createRefreshToken(userId int) models.Token {
	refreshToken := models.Token{
		TokenType: RefreshString,
		Exp:       RefreshTokenDuration,
		UserId:    userId,
	}
	refreshToken.Value = h.generateToken(userId, refreshToken.Exp)
	return refreshToken
}

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

func (h *Handler) authTokenFromCookies(c *gin.Context) (string, error) {
	cookieToken, err := c.Cookie(AuthString)
	if err != nil || cookieToken == "" {
		return "", err
	}
	tokenValue := strings.Split(cookieToken, " ")[1]
	return tokenValue, nil
}

func (h *Handler) refreshTokenFromCookies(c *gin.Context) (string, error) {
	cookieToken, err := c.Cookie(RefreshString)
	if err != nil || cookieToken == "" {
		return "", err
	}
	tokenValue := strings.Split(cookieToken, " ")[1]
	return tokenValue, nil
}

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

func (h *Handler) SetTokenToCookies(c *gin.Context, token models.Token) {
	c.SetCookie(token.TokenType, fmt.Sprintf("Bearer "+token.Value), int(token.Exp), "", "", true, true)
}

func (h *Handler) handleInvalidToken(code int, c *gin.Context) {
	h.deleteCookies(c)
	c.Redirect(code, "/user/login")
}
