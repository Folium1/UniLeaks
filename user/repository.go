package user

import (
	"context"
	"leaks/models"
)

type Repository interface {
	CreateUser(ctx context.Context, newUser models.User) (int,error)
	GetById(ctx context.Context, id int) (models.User, error)
	GetByMail(ctx context.Context, mail string) (models.User, error)
}
