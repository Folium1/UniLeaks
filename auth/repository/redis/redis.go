package repository

import (
	"context"
	"fmt"
	"log"
	"uniLeaks/config"
	"uniLeaks/models"

	"github.com/redis/rueidis"
)

type Repo struct {
	db *config.RedisConfig
}

// New returns a new instance of the auth repository.
func New() Repo {
	return Repo{config.NewRedisConfig()}
}

// SaveToken saves the given token to Redis with the appropriate expiration time.
func (r Repo) SaveToken(ctx context.Context, token models.Token) error {
	db := r.db.ConnectToRedis()
	defer db.Close()
	err := db.Do(ctx, db.B().Set().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Value)).Value(fmt.Sprintf("%v", token.UserId)).Build()).Error()

	if err != nil {
		log.Println("Redis:", err)
		db.B().Discard()
		return err
	}
	err = db.Do(ctx, db.B().Expire().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Value)).Seconds(int64(token.Exp)).Build()).Error()
	if err != nil {
		log.Println("Redis:", err)
		db.B().Discard()
		return err
	}

	return nil
}

// DeleteToken deletes the given token from Redis.
func (r Repo) DeleteToken(ctx context.Context, token models.Token) error {
	db := r.db.ConnectToRedis()
	defer db.Close()
	err := db.Do(ctx, db.B().Del().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Value)).Build()).Error()
	if err != nil {
		log.Println("Redis:", err)
		return err
	}
	return nil
}

// UserId returns the user ID associated with the given token from Redis.
func (r Repo) UserId(ctx context.Context, token models.Token) (int, error) {
	db := r.db.ConnectToRedis()
	defer db.Close()
	userId, err := db.Do(ctx, db.B().Get().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Value)).Build()).AsInt64()
	if err == rueidis.Nil {
		return -1, err
	}
	if err != nil {
		log.Println("Redis:", err)
		return -1, err
	}
	return int(userId), nil

}
