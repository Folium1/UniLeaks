package user

import (
	"context"
	"leaks/pkg/models"
)

type Receiver interface {
	UserById(ctx context.Context, id int) (models.User, error)
	UserByNick(ctx context.Context, nick string) (models.User, error)
}

type Creator interface {
	CreateUser(ctx context.Context, newUser models.User) (int, error)
}

type Lister interface {
	BannedMailHashes(ctx context.Context) ([]string, error)
}
