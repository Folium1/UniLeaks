package user

import (
	"errors"
	"log"
	"os"
	"strings"
	"uniLeaks/models"

	"golang.org/x/crypto/bcrypt"
)

func isPasswordValid(pass string) bool {
	return len(pass) > 8
}

func isMailValid(mail string) bool {
	return strings.HasSuffix(mail, os.Getenv("MAIL_DOMAIN"))
}

// hashPassword hashes the given password
func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("Couldn't hash user's password", err)
		return []byte(""), err
	}
	return hashedPassword, nil
}

// validateUserInput validates the email and password fields of a user
func validateUserInput(user models.User) error {
	if !isPasswordValid(user.Password) {
		return errors.New("Your password is too short, at least 8 characters")
	}
	if !isMailValid(user.Email) {
		return errors.New("Your email is incorrect")
	}
	return nil
}
