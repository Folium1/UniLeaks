package repository

import (
	// "log"

	// "uniLeaks/config"
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
func (r Repository) Create(ctx context.Context, newUser models.User) error {
	result := r.db.WithContext(ctx).Create(newUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetById returns the user record from the database with the specified ID.
func (r Repository) GetById(ctx context.Context, id string) (models.User, error) {
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
	result := r.db.WithContext(ctx).Model(&user).Select("id", "password", "banned").Where("email = %v", mail).Scan(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

// BanUser sets the banned flag on the user record with the specified ID.
func (r Repository) BanUser(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Set("banned", true).Where("id = %v", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
