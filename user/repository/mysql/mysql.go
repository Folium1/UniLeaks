package repository

import (
	"log"

	"uniLeaks/config"
	"uniLeaks/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// New creates a new instance of the Repository with a connection to the MySQL database.
func New() *Repository {
	db, err := config.MysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{db}
}

// Create creates a new user record in the database.
func (r Repository) Create(ctx context.Context, newUser models.User) (int, error) {
	err := r.db.WithContext(ctx).Create(&newUser).Last(&newUser).Error
	if err != nil {
		return -1, err
	}
	return newUser.ID, nil
}

// GetById returns the user record from the database with the specified ID.
func (r Repository) GetById(ctx context.Context, id int) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Model(&user).Where("id = %v", id).Scan(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

// GetByMail returns the user record from the database with the specified email address.
func (r Repository) GetByMail(ctx context.Context, mail string) (models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", mail).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

// BanUser sets the banned flag on the user record with the specified ID.
func (r Repository) BanUser(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Set("is_banned", true).Where("id = %v", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
