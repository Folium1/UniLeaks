package repository

import (
	"context"
	"fmt"

	"leaks/pkg/config"
	"leaks/pkg/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepo {
	db, err := config.MysqlConn()
	if err != nil {
		l.Fatal(fmt.Sprint("Couldn't connect to mysql, err: ", err))
	}
	return &UserRepo{db}
}

func (r *UserRepo) AllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		l.Error(fmt.Sprint("Couldn't get all users, err: ", result.Error))
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, id int) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		l.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepo) GetByNick(ctx context.Context, nick string) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("nick_name = ?", nick).First(&user)
	if result.Error != nil {
		l.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepo) BannedUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Where("is_banned = ?", true).Find(&users).Error
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get banned users, err: ", err))
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) BanUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Set("is_banned", true).Where("id = %v", id)
	if result.Error != nil {
		l.Error(fmt.Sprint("Couldn't ban user, err: ", result.Error))
		return result.Error
	}
	return nil
}

func (r *UserRepo) UnbanUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Set("is_banned", false).Where("id = %v", id)
	if result.Error != nil {
		l.Error(fmt.Sprint("Couldn't unban user, err: ", result.Error))
		return result.Error
	}
	return nil
}
