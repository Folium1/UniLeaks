package service

import (
	"log"
	"strconv"
	"uniLeaks/leaks"
	"uniLeaks/leaks/repository"
	"uniLeaks/models"
)

type Service struct {
	Repo repository.Repo
}

func New() leaks.Repository {
	return Service{repository.New()}
}

func (s Service) SaveFile(data models.LeakData) error {
	result, err := scanFile(data.File.OpenedFile)
	if !result {
		return leaks.VirusDetectedErr
	}
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.Repo.SaveFile(data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Service) GetList(data models.SubjectData) ([]models.LeakData, error) {
	files, err := s.Repo.GetList(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filesList := make([]models.LeakData, len(files))
	var file models.LeakData
	for _, i := range files {
		file.File.Name = i.Name
		file.File.Description = i.Description
		file.File.Id = i.Id
		file.File.Size = i.Size
		file.UserData.Dislikes, err = strconv.Atoi(i.Properties["dislikes"])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		file.UserData.Likes, err = strconv.Atoi(i.Properties["likes"])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		filesList = append(filesList, file)
	}
	return filesList, nil
}

func (s Service) GetFile(fileId string) ([]byte, error) {
	b, err := s.Repo.GetFile(fileId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return b, nil
}

func (s Service) DislikeFile(fileId string) error {
	err := s.Repo.DislikeFile(fileId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Service) LikeFile(fileId string) error {
	err := s.Repo.LikeFile(fileId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
