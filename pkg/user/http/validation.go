package http

import (
	"context"
	"errors"
	"os"
	"strings"
	"unicode"

	errHandler "leaks/pkg/err"
	"leaks/pkg/models"
	userService "leaks/pkg/user/service"

	"golang.org/x/crypto/bcrypt"
)

func isPasswordValid(password string) error {
	if len(password) < 8 {
		return errors.New("Пароль повинен містити хоча б 8 символівІ")
	}
	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)
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

func isNickNameValid(nickName string) error {
	if len(nickName) < 3 || len(nickName) > 20 {
		return errors.New("Нікнейм повинен містити від 3 до 20 символів")
	}
	return nil
}

func isMailValid(mail string) bool {
	return strings.HasSuffix(mail, os.Getenv("MAIL_DOMAIN"))
}

func hashString(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return []byte(""), err
	}
	return hashedPassword, nil
}

func checkIfMailIsBanned(userService *userService.UserUseCase, mail string) error {
	return userService.IsBanned(context.Background(), mail)
}

func validateUserInput(userService *userService.UserUseCase, user models.User) error {
	if err := isPasswordValid(user.Password); err != nil {
		return err
	}
	if err := isNickNameValid(user.NickName); err != nil {
		return err
	}
	if !isMailValid(user.Email) {
		return errors.New("Невірний формат мейла")
	}
	if err := checkIfMailIsBanned(userService, user.Email); err != nil {
		logger.Error(err.Error())
		if err == errHandler.UserIsBannedErr {
			return errors.New("Юзер з таким мейлом, був забанений раніше")
		} else {
			return err
		}
	}
	return nil
}
