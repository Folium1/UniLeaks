package user

import (
	"context"
	"errors"
	"leaks/pkg/models"
	userService "leaks/pkg/user/service"
	"os"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// isPasswordValid checks if the given password is valid
func isPasswordValid(password string) error {
	if len(password) < 8 {
		return errors.New("Пароль повинен містити хоча б 8 символівІ")
	}

	// Check for uppercase, lowercase, and numeric characters
	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)
	// Iterate over password and set flags when conditions are met
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}
	// Return an error if a condition is not met
	if !hasUpper {
		return errors.New("Пароль повинен містити хоча б одну велику літеру")
	}
	if !hasLower {
		return errors.New("Пароль повинен містити хоча б одну маленьку літеру")
	}
	if !hasDigit {
		return errors.New("Пароль повинен містити хоча б одну цифру")
	}
	return nil
}

// isNickNameValid checks if the given nickName is valid
func isNickNameValid(nickName string) error {
	if len(nickName) < 3 {
		return errors.New("Нікнейм повинен містити хоча б 3 символи")
	}
	return nil
}

// isMailValid checks if the given mail is valid
func isMailValid(mail string) bool {
	return strings.HasSuffix(mail, os.Getenv("MAIL_DOMAIN"))
}

// hashPassword hashes the given string
func hashString(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return []byte(""), err
	}
	return hashedPassword, nil
}

// checkIfUserBanned checks if the given user is banned
func checkIfMailIsBanned(userService *userService.UserUseCase, mail string) error {
	return userService.IsBanned(context.Background(), mail)
}

// validateUserInput validates the email and password fields of a user
func validateUserInput(user models.User) error {
	if err := isPasswordValid(user.Password); err != nil {
		return err
	}
	if err := isNickNameValid(user.NickName); err != nil {
		return err
	}
	if !isMailValid(user.Email) {
		return errors.New("Невірний формат мейла")
	}
	return nil
}
