package repository

import (
	"context"
	"leaks/pkg/models"

	"google.golang.org/api/drive/v3"
)

type FileLister interface {
	GetUserFilesList(userId string) ([]*drive.File, error)
	FilesOrderedByDislikes() ([]*drive.File, error)
}

type FileDeleter interface {
	DeleteFile(fileId string) error
	DeleteAllUserFiles(userId string) error
}

type FileReceiver interface {
	File(fileID string) ([]byte, drive.File, error)
}

type UserStatusSetter interface {
	BanUser(context.Context, string) error
	UnbanUser(context.Context, string) error
}

type UserReceiver interface {
	GetUserById(ctx context.Context, id int) (models.User, error)
	GetByNick(ctx context.Context, nick string) (models.User, error)
}

type UserLister interface {
	AllUsers(context.Context) ([]*models.User, error)
	BannedUsers(ctx context.Context) ([]*models.User, error)
}
