package service

import (
	"context"
	"fmt"

	adminRepo "leaks/pkg/admin"
	repo "leaks/pkg/admin/repository"
	errHandler "leaks/pkg/err"
	"leaks/pkg/models"
	"time"
)

type AdminUserService struct {
	receiver     adminRepo.UserReceiver
	lister       adminRepo.UserLister
	statusSetter adminRepo.UserStatusSetter
}

func NewAdminUserService() *AdminUserService {
	r := repo.NewUserRepository()
	return &AdminUserService{
		receiver:     r,
		lister:       r,
		statusSetter: r,
	}
}

func (s *AdminUserService) AllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	users, err := s.lister.AllUsers(ctx)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get all users: ", err))
		return nil, errHandler.UserListReceiveErr
	}
	return users, nil
}

func (s *AdminUserService) IsAdmin(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.receiver.GetUserById(ctx, id)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't check if user is admin: ", err))
		return models.User{}, errHandler.UserListReceiveErr
	}
	return user, nil
}

func (s *AdminUserService) GetByNick(nickName string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.receiver.GetByNick(ctx, nickName)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get user by nick: ", err))
		return user, errHandler.UserListReceiveErr
	}
	return user, nil
}

func (s *AdminUserService) BannedUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	users, err := s.lister.BannedUsers(ctx)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get banned users: ", err))
		return nil, errHandler.UserListReceiveErr
	}
	return users, nil
}

func (s *AdminUserService) BanUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.statusSetter.BanUser(ctx, id)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't ban user: ", err))
		return errHandler.BanUserErr
	}
	return nil
}

func (s *AdminUserService) UnbanUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.statusSetter.UnbanUser(ctx, id)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't unban user: ", err))
		return err
	}
	return nil
}
