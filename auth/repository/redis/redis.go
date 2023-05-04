package repository

import (
	"context"
	"errors"
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

func (r Repo) SaveToken(ctx context.Context, token models.Token) error {
	db := r.db.ConnectToRedis()
	defer db.Close()
	err := db.Do(ctx, db.B().Set().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Tk)).Value(fmt.Sprintf("%v", token.UserId)).Build()).Error()

	if err != nil {
		db.B().Discard()
	}
	err = db.Do(ctx, db.B().Expire().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Tk)).Seconds(int64(token.Exp)).Build()).Error()
	if err != nil {
		db.B().Discard()
	}

	return nil
}

func (r Repo) DeleteToken(ctx context.Context, token models.Token) error {
	db := r.db.ConnectToRedis()
	defer db.Close()
	err := db.Do(ctx, db.B().Del().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Tk)).Build()).Error()
	if err != nil {
		return err
	}
	return nil
}

func (r Repo) GetUserId(ctx context.Context, token models.Token) (int, error) {
	db := r.db.ConnectToRedis()
	defer db.Close()
	userId, err := db.Do(ctx, db.B().Get().Key(fmt.Sprintf("auth:%v:%v", token.TokenType, token.Tk)).Build()).AsInt64()
	if err == rueidis.Nil {
		return -1, errors.New("No user with that token")
	}
	if err != nil {
		return -1, err
	}
	return int(userId), nil

}
