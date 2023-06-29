package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"leaks/pkg/logger"
	"leaks/pkg/models"

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
	key, err := ioutil.ReadFile("client_credentials.json")
	if err != nil {
		logg.Fatal(fmt.Sprint("Couldn't read key file, err:", err))
		return nil, err
	}
	// Create a new Drive API client
	ctx := context.Background()
	clientOption := option.WithCredentialsJSON(key)
	service, err := drive.NewService(ctx, clientOption)
	if err != nil {
		logg.Fatal(fmt.Sprint("Couldn't create a new drive service, err:", err))
		return nil, err
	}
	return service, nil
}

// MysqlConn creates a new MySQL connection with retry mechanism
func MysqlConn() (*gorm.DB, error) {
	mysqlConn := os.Getenv("MYSQL")
	var db *gorm.DB
	var err error
	maxRetries := 10
	retryInterval := time.Second * 5
	// Retry loop
	for retries := 0; retries < maxRetries; retries++ {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: mysqlConn,
		}), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		logg.Info(fmt.Sprintf("Failed to connect to MySQL, err: %s. Retrying in %s...", err, retryInterval))
		time.Sleep(retryInterval)
	}
	logg.Fatal(fmt.Sprintf("Couldn't connect to MySQL after %d retries", maxRetries))
	return nil, fmt.Errorf("failed to connect to MySQL")
}

// InitMYSQL creates needed tables
func InitMYSQL() {
	db, err := MysqlConn()
	if err != nil {
		logg.Fatal(err.Error())
	}
	// Create tables
	err = db.AutoMigrate(models.User{})
	if err != nil {
		logg.Fatal(fmt.Sprint("Couldn't migrate mysql table, err:", err))
	}
}
