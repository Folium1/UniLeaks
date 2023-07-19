package service

import (
	"errors"
	"fmt"

	errHandler "leaks/pkg/err"
	"leaks/pkg/leaks"
	"leaks/pkg/leaks/repository"
	logg "leaks/pkg/logger"
	"leaks/pkg/models"
)

var l = logg.NewLogger()

type Service struct {
	saver      leaks.Saver
	downloader leaks.Downloader
	lister     leaks.Lister
	evaluator  leaks.Evaluator
}

func New() *Service {
	newRepo := repository.New()
	return &Service{
		saver:      newRepo,
		downloader: newRepo,
		lister:     newRepo,
		evaluator:  newRepo,
	}
}

func (s *Service) SaveFile(data models.Leak) error {
	result, err := scanFile(data.File.Content)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't scan file: ", err))
		return errHandler.FileCheckErr
	}
	if !result {
		l.Error(fmt.Sprintf("Virus detected, uploaded by: %v", data.User.UserId))
		return errHandler.VirusDetectedErr
	}

	err = s.saver.SaveFile(data)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't save file: ", err))
		return errHandler.FileSaveErr
	}
	return nil
}

func (s *Service) FilesList(data models.Subject) ([]models.Leak, error) {
	driveFiles, err := s.lister.FilesList(data)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get files list: ", err))
		return nil, errHandler.FileListReceiveErr
	}
	filesList, err := parseDriveDataToModel(driveFiles)
	if err != nil {
		return nil, err
	}
	return filesList, nil
}

func (s *Service) File(fileID string) (models.Leak, error) {
	if fileID == "" {
		l.Error("File ID is empty")
		return models.Leak{}, errHandler.FileNotFoundErr
	}
	b, file, err := s.downloader.DownloadFile(fileID)
	if err != nil {
		if errors.Is(err, errHandler.FileNotFoundErr) {
			return models.Leak{}, errHandler.FileNotFoundErr
		}
		return models.Leak{}, errHandler.FileReceiveErr
	}
	fileData := &models.File{
		Id:      file.Id,
		Name:    file.Name,
		Content: b,
	}
	leak := models.Leak{
		File:    fileData,
		User:    &models.UserFileData{UserId: file.Properties["userId"]},
		Subject: &models.Subject{Faculty: file.Properties["faculty"], Subject: file.Properties["subject"]},
	}
	return leak, nil
}

func (s *Service) LikeDislikeFile(data models.LikeDislike) error {
	if data.FileId == "" {
		l.Error("File ID is empty")
		return errHandler.FileNotFoundErr
	}
	if data.Action == "like" {
		err := s.evaluator.LikeFile(data.FileId, data.UserId)
		if err != nil {
			l.Error(fmt.Sprint("Couldn't like file: ", err))
			return errHandler.LikeDislikeErr
		}
	} else if data.Action == "dislike" {
		err := s.evaluator.DislikeFile(data.FileId, data.UserId)
		if err != nil {
			l.Error(fmt.Sprint("Couldn't dislike file: ", err))
			return errHandler.LikeDislikeErr
		}
	} else {
		l.Error("Invalid action")
		return errHandler.InvalidActionErr
	}
	return nil
}

func (s *Service) AllFiles() ([]models.Leak, error) {
	driveFiles, err := s.lister.AllFiles()
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get all files: ", err))
		return nil, errHandler.FileListReceiveErr
	}
	files, err := parseDriveDataToModel(driveFiles)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *Service) UserFiles(userID string) ([]models.Leak, error) {
	driveFiles, err := s.lister.UserFiles(userID)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get all files: ", err))
		return nil, errHandler.FileListReceiveErr
	}
	filesList, err := parseDriveDataToModel(driveFiles)
	if err != nil {
		return nil, err
	}
	return filesList, nil
}
