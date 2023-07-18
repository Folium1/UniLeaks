package repository

import (
	"fmt"

	"leaks/pkg/config"
	logg "leaks/pkg/logger"
	"leaks/pkg/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

var logger = logg.NewLogger()

type Repository struct {
	db *gorm.DB
}

func New() *Repository {
	db, err := config.MysqlConn()
	if err != nil {
		logger.Fatal(fmt.Sprint("Couldn't connect to mysql, err: ", err))
	}
	return &Repository{db}
}

func (r *Repository) CreateUser(ctx context.Context, newUser models.User) (int, error) {
	err := r.db.WithContext(ctx).Create(&newUser).Last(&newUser).Error
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't create user, err: ", err))
		return -1, err
	}
	return newUser.ID, nil
}

func (r *Repository) UserById(ctx context.Context, id int) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Scan(&user)
	if result.Error != nil {
		logger.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *Repository) UserByNick(ctx context.Context, nick string) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("nick_name = ?", nick).First(&user)
	if result.Error != nil {
		logger.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *Repository) BannedMailHashes(ctx context.Context) ([]string, error) {
	var mails []string
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("is_banned = ?", true).Pluck("email", &mails)
	if result.Error != nil {
		logger.Error(fmt.Sprint("Couldn't get banned mails, err: ", result.Error))
		return nil, result.Error
	}
	return mails, nil
}
