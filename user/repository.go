package user

import (
	"context"
	"uniLeaks/models"
)

type Repository interface {
	Create(ctx context.Context, newUser models.User) error
	GetById(ctx context.Context, id string) (models.User, error)
	GetByMail(ctx context.Context, mail string) (models.User, error)
}
