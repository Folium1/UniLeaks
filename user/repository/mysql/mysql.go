package repository

import (
	"fmt"
	"leaks/config"
	"leaks/logger"
	"leaks/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

var logg = logger.NewLogger()

type Repository struct {
	db *gorm.DB
}

// New creates a new instance of the Repository with a connection to the MySQL database.
func New() *Repository {
	db, err := config.MysqlConn()
	if err != nil {
		logg.Fatal(fmt.Sprint("Couldn't connect to mysql, err: ", err))
	}
	return &Repository{db}
}

// CreateUser creates a new user record in the database.
func (r *Repository) CreateUser(ctx context.Context, newUser models.User) (int, error) {
	err := r.db.WithContext(ctx).Create(&newUser).Last(&newUser).Error
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't create user, err: ", err))
		return -1, err
	}
	return newUser.ID, nil
}

// GetById returns the user record from the database with the specified ID.
func (r *Repository) GetById(ctx context.Context, id int) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Scan(&user)
	if result.Error != nil {
		logg.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return models.User{}, result.Error
	}
	return user, nil
}

// GetByMail returns the user record from the database with the specified email address.
func (r *Repository) GetByMail(ctx context.Context, mail string) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", mail).First(&user)
	if result.Error != nil {
		logg.Error(fmt.Sprint("Couldn't get user, err: ", result.Error))
		return models.User{}, result.Error
	}
	return user, nil
}
