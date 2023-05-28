package auth

import (
	"context"
	"errors"
	repository "uniLeaks/auth/repository/redis"
	"uniLeaks/models"

	"github.com/redis/rueidis"
)

type UseCase struct {
	repo repository.Repo
}

// NewUseCase returns a new instance of the auth use case.
func NewUseCase() UseCase {
	return UseCase{repository.New()}
}

// SaveToken saves the given token to Redis with the appropriate expiration time.
func (u UseCase) SaveToken(ctx context.Context, token models.Token) error {
	err := u.repo.SaveToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

// DeleteToken deletes the given token from Redis.
func (u UseCase) DeleteToken(ctx context.Context, token models.Token) error {
	err := u.repo.DeleteToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

// GetUserId returns the user ID associated with the given token from Redis.
func (u UseCase) UserId(ctx context.Context, token models.Token) (int, error) {
	id, err := u.repo.UserId(ctx, token)
	if err != nil {
		if err == rueidis.Nil {
			return -1, errors.New("No user with that token")
		}
		return 0, err
	}
	return id, nil
}
