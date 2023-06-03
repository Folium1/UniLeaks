package service

import (
	"errors"
	"runtime"
	"strconv"
	"uniLeaks/leaks"
	"uniLeaks/models"
)

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
		return err
	}
	if !result {
		return leaks.ErrVirusDetected
	}

	err = s.repo.SaveFile(data)
	if err != nil {
		runtime.GC()
		return err
	}
	runtime.GC()
	return nil
}

// FilesList retrieves a list of files based on subject data from the repository.
func (s *Service) FilesList(data models.SubjectData) ([]models.LeakData, error) {
	files, err := s.repo.FilesList(data)
	if err != nil {
		return nil, err
	}
	filesList := make([]models.LeakData, 0, len(files))
	// Iterate over the files and create a list of LeakData
	for _, f := range files {
		file := &models.File{
			Id:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			Size:        f.Size,
		}
		dislikes, err := strconv.Atoi(f.Properties["dislikes"])
		if err != nil {
			return nil, err
		}
		likes, err := strconv.Atoi(f.Properties["likes"])
		if err != nil {
			return nil, err
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
		return models.LeakData{}, leaks.ErrFileNotFound
	}
	// Retrieve the file from the repository
	b, fileData, err := s.repo.File(fileID)
	if err != nil {
		if errors.Is(err, leaks.ErrFileNotFound) {
			return models.LeakData{}, leaks.ErrFileNotFound
		}
		return models.LeakData{}, err
	}
	// Create a new file and leak data
	file := &models.File{
		Id:      fileData.Id,
		Name:    fileData.Name,
		Content: b,
	}
	leakData := models.LeakData{
		File:     file,
		UserData: &models.UserFileData{},
		Subject:  &models.SubjectData{},
	}
	runtime.GC()
	return leakData, nil
}

// DislikeFile increments the dislike count of a file in the repository.
func (s *Service) DislikeFile(fileID string) error {
	err := s.repo.DislikeFile(fileID)
	if err != nil {
		return err
	}
	return nil
}

// LikeFile increments the like count of a file in the repository.
func (s *Service) LikeFile(fileID string) error {
	err := s.repo.LikeFile(fileID)
	if err != nil {
		return err
	}
	return nil
}

// AllFiles retrieves all files from the repository.
func (s *Service) AllFiles() ([]models.LeakData, error) {
	files, err := s.repo.AllFiles()
	if err != nil {
		return nil, err
	}
	// Create a list of LeakData
	filesList := make([]models.LeakData, 0, len(files))
	for _, f := range files {
		file := &models.File{
			Id:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			Size:        f.Size,
		}
		userData := &models.UserFileData{}
		dislikes, err := strconv.Atoi(f.Properties["dislikes"])
		if err != nil {
			return nil, err
		}
		likes, err := strconv.Atoi(f.Properties["likes"])
		if err != nil {
			return nil, err
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
		return nil, err
	}
	// Create a list of LeakData
	filesList := make([]models.LeakData, 0, len(files))
	for _, f := range files {
		file := &models.File{
			Id:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			Size:        f.Size,
		}
		userData := &models.UserFileData{}
		dislikes, err := strconv.Atoi(f.Properties["dislikes"])
		if err != nil {
			return nil, err
		}
		likes, err := strconv.Atoi(f.Properties["likes"])
		if err != nil {
			return nil, err
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
