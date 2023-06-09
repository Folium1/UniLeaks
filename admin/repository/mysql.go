package repository

import (
	"context"
	"fmt"
	"leaks/config"
	"leaks/models"

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
		logg.Fatal(fmt.Sprint("Couldn't connect to mysql, err: ", err))
	}
	return &UserRepo{db}
}

// BanUser sets the banned flag on the user record with the specified ID.
func (r *UserRepo) BanUser(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Set("is_banned", true).Where("id = %v", id)
	if result.Error != nil {
		logg.Error(fmt.Sprint("Couldn't ban user, err: ", result.Error))
		return result.Error
	}
	return nil
}

// AllUsers returns all users
func (r *UserRepo) AllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		logg.Error(fmt.Sprint("Couldn't get all users, err: ", result.Error))
		return nil, result.Error
	}
	return users, nil
}

// IsAdmin checks if user is admin
func (r *UserRepo) IsAdmin(ctx context.Context, id int) (bool, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		logg.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return false, result.Error
	}
	return user.IsAdmin, nil
}
