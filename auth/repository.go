package auth

import (
	"context"
	"uniLeaks/models"
)

type Repository interface {
	SaveToken(ctx context.Context, token models.Token) error
	DeleteToken(ctx context.Context, token models.Token) error
	GetUserId(ctx context.Context, token models.Token) (int, error)
}
