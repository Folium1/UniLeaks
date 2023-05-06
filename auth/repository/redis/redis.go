package repository

import (
	"context"
	"fmt"
	"uniLeaks/config"
	"uniLeaks/models"

	"github.com/redis/rueidis"
)

type Repo struct {
	db *config.RedisConfig
}

func New() Repo {
	return Repo{config.NewRedisConfig()}
}

// SaveToken saves the given token to Redis with the appropriate expiration time.
func (r Repo) SaveToken(ctx context.Context, token models.Token) error {
	db := r.db.ConnectToRedis()
	defer db.Close()
	err := db.Do(ctx, db.B().Set().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Value)).Value(fmt.Sprintf("%v", token.UserId)).Build()).Error()

	if err != nil {
		db.B().Discard()
	}
	err = db.Do(ctx, db.B().Expire().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Value)).Seconds(int64(token.Exp)).Build()).Error()
	if err != nil {
		db.B().Discard()
	}

	return nil
}

// DeleteToken deletes the given token from Redis.
func (r Repo) DeleteToken(ctx context.Context, token models.Token) error {
	db := r.db.ConnectToRedis()
	defer db.Close()
	err := db.Do(ctx, db.B().Del().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Value)).Build()).Error()
	if err != nil {
		return err
	}
	return nil
}

// GetUserId returns the user ID associated with the given token from Redis.
func (r Repo) GetUserId(ctx context.Context, token models.Token) (int, error) {
	db := r.db.ConnectToRedis()
	defer db.Close()
	userId, err := db.Do(ctx, db.B().Get().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Value)).Build()).AsInt64()
	if err == rueidis.Nil {
		return -1, err
	}
	if err != nil {
		return -1, err
	}
	return int(userId), nil

}
