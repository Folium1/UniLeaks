package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"uniLeaks/models"
	"uniLeaks/user"
	repository "uniLeaks/user/repository/mysql"
)

type UserUseCase struct {
	repo repository.Repository
}

// New creates a new UserUseCase with the given repository.
func New(repository repository.Repository) user.Repository {
	return &UserUseCase{repository}
}

// Create creates a new user in the repository.
func (u UserUseCase) Create(ctx context.Context, newUser models.User) (int, error) {
	userId, err := u.repo.Create(ctx, newUser)
	if err != nil {
		log.Println(err)
		return -1, errors.New("Couldn't create user")
	}
	return userId, nil
}

// GetById gets the user with the given id from the repository.
func (u UserUseCase) GetById(ctx context.Context, id int) (models.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		log.Println(err)
		return user, errors.New("Couldn't get user with that id")
	}
	return user, nil
}

// GetByMail gets the user with the given email from the repository.
func (u UserUseCase) GetByMail(ctx context.Context, mail string) (models.User, error) {
	user, err := u.repo.GetByMail(ctx, mail)
	if err != nil {
		log.Println(err)
		return models.User{}, errors.New("Couldn't get user, by mail")
	}
	return user, nil
}

// BanUser bans the user with the given id.
func (u UserUseCase) BanUser(ctx context.Context, id int) error {
	err := u.repo.BanUser(ctx, id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Couldn't ban user with id:%v", id)
	}
	return nil
}
