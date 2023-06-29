package repository

import (
	"context"
	"fmt"
	"leaks/pkg/config"
	"leaks/pkg/models"

	"gorm.io/gorm"
)

// UserRepo is the MySQL repository for acting with users.
type UserRepo struct {
	db *gorm.DB
}

// New creates a new instance of the Repository with a connection to the MySQL database.
func NewUserRepository() *UserRepo {
	db, err := config.MysqlConn()
	if err != nil {
		logger.Fatal(fmt.Sprint("Couldn't connect to mysql, err: ", err))
	}
	return &UserRepo{db}
}

// BanUser sets the banned flag on the user record with the specified ID.
func (r *UserRepo) BanUser(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Set("is_banned", true).Where("id = %v", id)
	if result.Error != nil {
		logger.Error(fmt.Sprint("Couldn't ban user, err: ", result.Error))
		return result.Error
	}
	return nil
}

// AllUsers returns all users
func (r *UserRepo) AllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		logger.Error(fmt.Sprint("Couldn't get all users, err: ", result.Error))
		return nil, result.Error
	}
	return users, nil
}

// IsAdmin checks if user is admin
func (r *UserRepo) IsAdmin(ctx context.Context, id int) (bool, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		logger.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return false, result.Error
	}
	return user.IsAdmin, nil
}

// GetByNick returns the user record from the database with the specified email address.
func (r *UserRepo) GetByNick(ctx context.Context, nick string) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("nick_name = ?", nick).First(&user)
	if result.Error != nil {
		logger.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return models.User{}, result.Error
	}
	return user, nil
}

// GetBannedUsers returns all banned users
func (r *UserRepo) GetBannedUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Where("is_banned = ?", true).Find(&users).Error
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get banned users, err: ", err))
		return nil, err
	}
	return users, nil
}

// UnbanUser sets the banned flag on the user record with the specified ID.
func (r *UserRepo) UnbanUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Set("is_banned", false).Where("id = %v", id)
	if result.Error != nil {
		logger.Error(fmt.Sprint("Couldn't unban user, err: ", result.Error))
		return result.Error
	}
	return nil
}
