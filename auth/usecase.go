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

func NewUseCase() UseCase {
	return UseCase{repository.New()}
}

func (u UseCase) SaveToken(ctx context.Context, token models.Token) error {
	err := u.repo.SaveToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (u UseCase) DeleteToken(ctx context.Context, token models.Token) error {
	err := u.repo.DeleteToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (u UseCase) GetUserId(ctx context.Context, token models.Token) (int, error) {
	id, err := u.repo.GetUserId(ctx, token)
	if err != nil {
		if err == rueidis.Nil {
			return -1, errors.New("No user with that token")
		}
		return 0, err
	}
	if err != nil {
		return 0, errors.New("Couldn't parse id into int")
	}
	return id, nil
}
