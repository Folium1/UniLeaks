package user

import (
	"context"
	"errors"
	"fmt"
	errHandler "leaks/pkg/err"
	"leaks/pkg/logger"
	"leaks/pkg/models"
	repository "leaks/pkg/user"

	"golang.org/x/crypto/bcrypt"
)

var logg = logger.NewLogger()

type UserUseCase struct {
	repo repository.Repository
}

// New creates a new UserUseCase with the given repository.
func New(repository repository.Repository) UserUseCase {
	return UserUseCase{repository}
}

// Create creates a new user in the repository.
func (u *UserUseCase) CreateUser(ctx context.Context, newUser models.User) (int, error) {
	userId, err := u.repo.CreateUser(ctx, newUser)
	// Check if the error is a duplicate entry error.
	if errHandler.IsDuplicateEntryError(err) {
		logg.Error(fmt.Sprint("Couldn't create user, err: ", err))
		return -1, errors.New("Юзер з таким мейлом або ніком вже існує")
	}
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't create user, err: ", err))
		return 0, errors.New("Помилка при регістрації, спробуйте ще раз")
	}
	return userId, nil
}

// GetById gets the user with the given id from the repository.
func (u *UserUseCase) GetById(ctx context.Context, id int) (models.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get user, err: ", err))
		return user, errors.New("Помилка отримання данних, спробуйте ще раз")
	}
	return user, nil
}

// GetByNick gets the user with the given email from the repository.
func (u *UserUseCase) GetByNick(ctx context.Context, nick string) (models.User, error) {
	user, err := u.repo.GetByNick(ctx, nick)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get user by mail, err: ", err))
		return models.User{}, errors.New("Юзер з таким ніком не знайдений")
	}
	return user, nil
}

// IsBanned checks if the user with the given email is banned. If the user is banned, it returns an error = UserIsBannedErr,else returns nil.
func (u *UserUseCase) IsBanned(ctx context.Context, userMail string) error {
	mails, err := u.repo.BannedMails(ctx)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get banned mails, err: ", err))
		return errHandler.ServerErr
	}
	// Check if the user is banned by comparing the user's email with the banned emails hash.
	// If the user is banned - return UserIsBannedErr.
	for _, mail := range mails {
		if err = bcrypt.CompareHashAndPassword([]byte(mail), []byte(userMail)); err == nil {
			return errHandler.UserIsBannedErr
		}
	}
	return nil
}
