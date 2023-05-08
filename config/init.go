package config

import (
	"context"
	"log"
	"os"

	"uniLeaks/models"

	"github.com/joho/godotenv"
	"github.com/redis/rueidis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Couldn't load local variables, err:", err)
	}
}

// MysqlConn connects to mysql, logs if error has occured
func MysqlConn() (*gorm.DB, error) {
	mysqlConn := os.Getenv("MYSQL")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               mysqlConn,
		DefaultStringSize: 50,
	}), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	return db, nil
}

// InitMYSQL creates needed tables
func InitMYSQL() {
	db, err := MysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(models.User{})
	if err != nil {
		log.Fatal("Couldn't migrate mysql, err = ", err)
	}
}

type RedisConfig struct {
	Addres   string
	Password string
	DB       int
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), 0}
}

// NewRedisConfig creates a new Redis configuration object from environment variables
func (r RedisConfig) ConnectToRedis() rueidis.Client {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{r.Addres},
	})
	if err != nil {
		log.Fatal(err)
	}
	err = client.Do(context.Background(), client.B().Ping().Build()).Error()
	if err != nil {
		log.Fatal(err)
	}
	return client
}
