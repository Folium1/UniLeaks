package repository

import (
	"context"
	"uniLeaks/models"

	"google.golang.org/api/drive/v2"
)

type DriveRepository interface {
	FilesList() ([]*drive.File, error)
	DeleteFile(fileId string) error
	File(fileID string) (models.LeakData, error)
}

type UserRepository interface {
	BanUser(userId int) error
	AllUsers() ([]*models.User, error)
	IsAdmin(ctx context.Context, userId int) (bool, error)
}
