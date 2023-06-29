package leaks

import (
	"leaks/pkg/models"

	"google.golang.org/api/drive/v3"
)

type Repository interface {
	SaveFile(data models.LeakData) error
	FilesList(data models.SubjectData) ([]*drive.File, error)
	DislikeFile(fileId, userID string) error
	LikeFile(fileId, userID string) error
	File(fileID string) ([]byte, drive.File, error)
	AllFiles() ([]*drive.File, error)
	MyFiles(userId string) ([]*drive.File, error)
}
