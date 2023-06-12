package repository

import (
	"context"
	"leaks/models"

	"google.golang.org/api/drive/v2"
)

type DriveRepository interface {
	FilesList() ([]*drive.File, error)
	DeleteFile(fileId string) error
	File(fileID string) (models.LeakData, error)
	DeleteAllUserFiles(userId string)
}

type UserRepository interface {
	BanUser(userId int) error
	AllUsers() ([]*models.User, error)
	IsAdmin(ctx context.Context, userId int) (bool, error)
	GetByNick(ctx context.Context, nick string) (models.User, error)
}
