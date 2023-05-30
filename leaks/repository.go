package leaks

import (
	"uniLeaks/models"

	"google.golang.org/api/drive/v3"
)

type Repository interface {
	SaveFile(data *models.LeakData) error
	FilesList(data models.SubjectData) ([]*drive.File, error)
	DislikeFile(fileId string) error
	LikeFile(fileId string) error
	File(fileId string) ([]byte, *drive.File, error)
	AllFiles() ([]*drive.File, error)
}
