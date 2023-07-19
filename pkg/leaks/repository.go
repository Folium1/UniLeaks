package leaks

import (
	"leaks/pkg/models"

	"google.golang.org/api/drive/v3"
)

type Saver interface {
	SaveFile(data models.Leak) error
}

type Downloader interface {
	DownloadFile(fileID string) ([]byte, drive.File, error)
}

type Lister interface {
	FilesList(data models.Subject) ([]*drive.File, error)
	AllFiles() ([]*drive.File, error)
	UserFiles(userId string) ([]*drive.File, error)
}

type Evaluator interface {
	DislikeFile(fileId, userID string) error
	LikeFile(fileId, userID string) error
}
