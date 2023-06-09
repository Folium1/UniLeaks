package user

import (
	"context"
	"errors"
	"fmt"
	"leaks/logger"
	"leaks/models"
	"leaks/user"
	repository "leaks/user"

	"github.com/go-sql-driver/mysql"
)

var logg = logger.NewLogger()

type UserUseCase struct {
	repo repository.Repository
}

// New creates a new UserUseCase with the given repository.
func New(repository repository.Repository) user.Repository {
	return &UserUseCase{repository}
}

func isDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062 // MySQL error code for duplicate entry
	}
	return false
}

// Create creates a new user in the repository.
func (u UserUseCase) CreateUser(ctx context.Context, newUser models.User) (int, error) {
	userId, err := u.repo.CreateUser(ctx, newUser)
	// Check if the error is a duplicate entry error.
	if isDuplicateEntryError(err) {
		logg.Error(fmt.Sprint("Couldn't create user, err: ", err))
		return -1, errors.New("Юзер з таким мейлом вже існує")
	}
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't create user, err: ", err))
		return 0, errors.New("Помилка при регістрації, спробуйте ще раз")
	}
	return userId, nil
}

// GetById gets the user with the given id from the repository.
func (u UserUseCase) GetById(ctx context.Context, id int) (models.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get user, err: ", err))
		return user, errors.New("Couldn't get user with that id")
	}
	return user, nil
}

// GetByMail gets the user with the given email from the repository.
func (u UserUseCase) GetByMail(ctx context.Context, mail string) (models.User, error) {
	user, err := u.repo.GetByMail(ctx, mail)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get user by mail, err: ", err))
		return models.User{}, errors.New("Couldn't get user, by mail")
	}
	return user, nil
}
