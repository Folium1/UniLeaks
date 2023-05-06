package user

import (
	"context"
	"uniLeaks/models"
)

type Repository interface {
	Create(ctx context.Context, newUser models.User) (int,error)
	GetById(ctx context.Context, id int) (models.User, error)
	GetByMail(ctx context.Context, mail string) (models.User, error)
	BanUser(ctx context.Context, id int) error
}
