package leaks

import (
	"uniLeaks/models"
)

type Repository interface {
	SaveFile(data models.LeakData) error
	GetList(data models.SubjectData) ([]models.LeakData, error)
	DislikeFile(fileId string) error
	LikeFile(fileId string) error
	GetFile(fileId string) ([]byte, error)
}
