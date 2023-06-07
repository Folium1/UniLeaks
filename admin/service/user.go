package service

import (
	"context"
	"time"
	adminRepo "uniLeaks/admin/repository"
	"uniLeaks/models"
)

type UserService struct {
	repo *adminRepo.UserRepo
}

// NewUserService creates a new instance of the service.
func NewUserService() UserService {
	repo := adminRepo.NewUserRepository()
	return UserService{
		repo: repo,
	}
}

// BanUser bans user
func (s *UserService) BanUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.BanUser(ctx, id)
}

// AllUsers returns all users
func (s *UserService) AllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.AllUsers(ctx)
}

// IsAdmin checks if user is admin
func (s *UserService) IsAdmin(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.IsAdmin(ctx, id)
}
