package user

import (
	"context"
	"leaks/pkg/models"
)

type Repository interface {
	CreateUser(ctx context.Context, newUser models.User) (int, error)
	GetById(ctx context.Context, id int) (models.User, error)
	GetByNick(ctx context.Context, nick string) (models.User, error)
	BannedMails(ctx context.Context) ([]string, error)
}
