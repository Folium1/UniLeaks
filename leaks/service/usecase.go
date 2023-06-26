package service

import (
	"errors"
	"fmt"
	errHandler "leaks/err"
	"leaks/leaks"
	"leaks/logger"
	"leaks/models"
	"strconv"
)

var logg = logger.NewLogger()

type Service struct {
	repo leaks.Repository
}

// New creates a new instance of the service.
func New(repo leaks.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// SaveFile saves the file to the repository.
func (s *Service) SaveFile(data models.LeakData) error {
	result, err := scanFile(data.File.Content)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't scan file: ", err))
		return errHandler.FileCheckErr
	}
	if !result {
		logg.Error(fmt.Sprintf("Virus detected, uploaded by: %v", data.UserData.UserId))
		return errHandler.VirusDetectedErr
	}

	err = s.repo.SaveFile(data)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't save file: ", err))
		return errHandler.FileSaveErr
	}
	return nil
}

// FilesList retrieves a list of files based on subject data from the repository.
func (s *Service) FilesList(data models.SubjectData) ([]models.LeakData, error) {
	files, err := s.repo.FilesList(data)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get files list: ", err))
		return nil, errHandler.FileListReceiveErr
	}
	filesList := make([]models.LeakData, 0, len(files))
	// Iterate over the files and create a list of LeakData
	for _, f := range files {
		file := &models.File{
			Id:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			Size:        float64(f.Size) / 1024 / 1024,
		}
		dislikesStr, ok := f.Properties["dislikes"]
		if !ok {
			logg.Error("Couldn't find dislikes property")
			return nil, errHandler.FileListReceiveErr
		}
		dislikes, err := strconv.Atoi(dislikesStr)
		if err != nil {
			logg.Error("Couldn't convert dislikes to int")
			return nil, errHandler.FileListReceiveErr
		}
		likes, err := strconv.Atoi(f.Properties["likes"])
		if err != nil {
			logg.Error("Couldn't convert likes to int")
			return nil, errHandler.FileListReceiveErr
		}
		// Create a new file and user data
		userData := &models.UserFileData{
			Dislikes: dislikes,
			Likes:    likes,
		}
		leakData := models.LeakData{
			File:     file,
			UserData: userData,
			Subject:  &models.SubjectData{},
		}
		filesList = append(filesList, leakData)
	}

	return filesList, nil
}

// File retrieves a specific file from the repository.
func (s *Service) File(fileID string) (models.LeakData, error) {
	// Check if the file exists
	if fileID == "" {
		logg.Error("File ID is empty")
		return models.LeakData{}, errHandler.FileNotFoundErr
	}
	// Retrieve the file from the repository
	b, fileData, err := s.repo.File(fileID)
	if err != nil {
		if errors.Is(err, errHandler.FileNotFoundErr) {
			return models.LeakData{}, errHandler.FileNotFoundErr
		}
		return models.LeakData{}, errHandler.FileReceiveErr
	}
	// Create a new file and leak data
	file := &models.File{
		Id:      fileData.Id,
		Name:    fileData.Name,
		Content: b,
	}
	leakData := models.LeakData{
		File:     file,
		UserData: &models.UserFileData{UserId: fileData.Properties["userId"]},
		Subject:  &models.SubjectData{Faculty: fileData.Properties["faculty"], Subject: fileData.Properties["subject"]},
	}
	return leakData, nil
}

// / LikeDislikeFile likes or dislikes a file.
func (s *Service) LikeDislikeFile(data models.LikeDislikeData) error {
	// Check if the file exists
	if data.FileId == "" {
		logg.Error("File ID is empty")
		return errHandler.FileNotFoundErr
	}
	// Retrieve the file from the repository
	if data.Action == "like" {
		err := s.repo.LikeFile(data.FileId, data.UserId)
		if err != nil {
			logg.Error(fmt.Sprint("Couldn't like file: ", err))
			return errHandler.LikeDislikeErr
		}
	} else if data.Action == "dislike" {
		err := s.repo.DislikeFile(data.FileId, data.UserId)
		if err != nil {
			logg.Error(fmt.Sprint("Couldn't dislike file: ", err))
			return errHandler.LikeDislikeErr
		}
	} else {
		logg.Error("Invalid action")
		return errHandler.InvalidActionErr
	}
	return nil
}

// AllFiles retrieves all files from the repository.
func (s *Service) AllFiles() ([]models.LeakData, error) {
	files, err := s.repo.AllFiles()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get all files: ", err))
		return nil, errHandler.FileListReceiveErr
	}
	// Create a list of LeakData
	filesList := make([]models.LeakData, 0, len(files))
	// Iterate over the files and create a list of LeakData
	for _, f := range files {
		file := &models.File{
			Id:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			Size:        float64(f.Size) / 1024 / 1024,
		}
		userData := &models.UserFileData{}
		dislikes, err := strconv.Atoi(f.Properties["dislikes"])
		if err != nil {
			logg.Error("Couldn't convert dislikes to int")
			return nil, errHandler.FileListReceiveErr
		}
		likes, err := strconv.Atoi(f.Properties["likes"])
		if err != nil {
			logg.Error("Couldn't convert likes to int")
			return nil, errHandler.FileListReceiveErr
		}
		userData.Dislikes = dislikes
		userData.Likes = likes
		leakData := models.LeakData{
			File:     file,
			UserData: userData,
			Subject:  &models.SubjectData{},
		}
		filesList = append(filesList, leakData)
	}

	return filesList, nil
}

// MyFiles retrieves all files from the repository.
func (s *Service) MyFiles(userID string) ([]models.LeakData, error) {
	files, err := s.repo.MyFiles(userID)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get all files: ", err))
		return nil, errHandler.FileListReceiveErr
	}
	// Create a list of LeakData
	filesList := make([]models.LeakData, 0, len(files))
	// Iterate over the files and create a list of LeakData
	for _, f := range files {
		file := &models.File{
			Id:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			Size:        float64(f.Size) / 1024 / 1024,
		}
		userData := &models.UserFileData{}
		dislikes, err := strconv.Atoi(f.Properties["dislikes"])
		if err != nil {
			logg.Error("Couldn't convert dislikes to int")
			return nil, errHandler.FileListReceiveErr
		}
		likes, err := strconv.Atoi(f.Properties["likes"])
		if err != nil {
			logg.Error("Couldn't convert likes to int")
			return nil, errHandler.FileListReceiveErr
		}
		userData.Dislikes = dislikes
		userData.Likes = likes
		leakData := models.LeakData{
			File:     file,
			UserData: userData,
			Subject:  &models.SubjectData{},
		}
		filesList = append(filesList, leakData)
	}
	return filesList, nil
}
