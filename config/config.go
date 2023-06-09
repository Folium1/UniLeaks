package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"leaks/logger"
	"leaks/models"

	"github.com/joho/godotenv"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logg = logger.NewLogger()

// Init initializes the environment variables
func init() {
	err := godotenv.Load()
	if err != nil {
		logg.Fatal(fmt.Sprint("Couldn't load local variables, err:", err))
	}
}

// NewDriveClient creates a new Google Drive client
func NewDriveClient() (*drive.Service, error) {
	key, err := ioutil.ReadFile("leaks-386216-becb63cca935.json")
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't read key file, err:", err))
		return nil, err
	}

	// Create a new Drive API client
	ctx := context.Background()
	clientOption := option.WithCredentialsJSON(key)
	service, err := drive.NewService(ctx, clientOption)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't create a new drive service, err:", err))
		return nil, err
	}
	return service, nil
}

// MysqlConn creates a new MySQL connection
func MysqlConn() (*gorm.DB, error) {
	mysqlConn := os.Getenv("MYSQL")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               mysqlConn,
		DefaultStringSize: 50,
	}), &gorm.Config{})
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't connect to mysql, err:", err))
	}
	return db, nil
}

// InitMYSQL creates needed tables
func InitMYSQL() {
	db, err := MysqlConn()
	if err != nil {
		logg.Fatal(fmt.Sprint("Couldn't connect to mysql, err:", err))
	}
	// Create tables
	err = db.AutoMigrate(models.User{})
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't migrate mysql table, err:", err))
	}
}