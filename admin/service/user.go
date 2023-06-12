package service

import (
	"context"
	"fmt"
	adminRepo "leaks/admin/repository"
	errHandler "leaks/err"
	"leaks/models"
	"time"
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
	err := s.repo.BanUser(ctx, id)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't ban user: ", err))
		return errHandler.BanUserErr
	}
	return nil
}

// AllUsers returns all users
func (s *UserService) AllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	users, err := s.repo.AllUsers(ctx)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get all users: ", err))
		return nil, errHandler.UserListReceiveErr
	}
	return users, nil
}

// IsAdmin checks if user is admin
func (s *UserService) IsAdmin(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	isAdmin, err := s.repo.IsAdmin(ctx, id)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't check if user is admin: ", err))
		return false, errHandler.UserListReceiveErr
	}
	return isAdmin, nil
}

// GetByNick returns the user record from the database with the specified email address.
func (s *UserService) GetByNick(nick string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.repo.GetByNick(ctx, nick)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get user by nick: ", err))
		return user, errHandler.UserListReceiveErr
	}
	return user, nil
}

// GetBannedUsers returns all banned users
func (s *UserService) GetBannedUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	users, err := s.repo.GetBannedUsers(ctx)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get banned users: ", err))
		return nil, errHandler.UserListReceiveErr
	}
	return users, nil
}

// UnbanUser unbans user
func (s *UserService) UnbanUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.repo.UnbanUser(ctx, id)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't unban user: ", err))
		return err
	}
	return nil
}
