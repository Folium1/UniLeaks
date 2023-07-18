package user

import (
	"context"
	"errors"
	"fmt"

	errHandler "leaks/pkg/err"
	"leaks/pkg/logger"
	"leaks/pkg/models"
	"leaks/pkg/user"
	repository "leaks/pkg/user/repository/mysql"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var logg = logger.NewLogger()

type UserUseCase struct {
	creator      user.Creator
	userReceiver user.Receiver
	mailLister   user.Lister
}

func New() UserUseCase {
	newRepo := repository.New()
	return UserUseCase{
		creator:      newRepo,
		userReceiver: newRepo,
		mailLister:   newRepo,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, newUser models.User) (int, error) {
	userId, err := u.creator.CreateUser(ctx, newUser)
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

func (u *UserUseCase) GetById(ctx context.Context, id int) (models.User, error) {
	receivedUser, err := u.userReceiver.UserById(ctx, id)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get user, err: ", err))
		return models.User{}, errors.New("Помилка отримання данних, спробуйте ще раз")
	}
	return receivedUser, nil
}

func (u *UserUseCase) UserByNick(ctx context.Context, nick string) (models.User, error) {
	receivedUser, err := u.userReceiver.UserByNick(ctx, nick)
	if err != nil {
		if errHandler.IsDuplicateEntryError(err) {
			return models.User{}, errors.New("Юзер з таким ніком вже зареєстрований")
		} else if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("Невірний нікнейм або пароль")
		}
		logg.Error(err.Error())
		return models.User{}, errors.New("Сталась помилка на сервері, спробуйте ще раз")
	}
	return receivedUser, nil
}

// IsBanned checks if the user with the given email is banned. If the user is banned, it returns an error = UserIsBannedErr,else returns nil.
func (u *UserUseCase) IsBanned(ctx context.Context, userMail string) error {
	mails, err := u.mailLister.BannedMailHashes(ctx)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get banned mails, err: ", err))
		return errHandler.ServerErr
	}
	for _, mail := range mails {
		if err = bcrypt.CompareHashAndPassword([]byte(mail), []byte(userMail)); err == nil {
			return errHandler.UserIsBannedErr
		}
	}
	return nil
}
