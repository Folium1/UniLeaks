package service

import (
	"errors"
	"fmt"
	adminRepo "leaks/admin/repository"
	errHandler "leaks/err"
	"leaks/logger"
	"leaks/models"
	"strconv"

	_ "google.golang.org/api/drive/v3"
)

var logg = logger.NewLogger()

type LeakService struct {
	repo *adminRepo.DriveRepo
}

// NewLeakService creates a new instance of the service.
func NewLeakService() LeakService {
	repo := adminRepo.NewDriveRepo()
	return LeakService{
		repo: repo,
	}
}

// FilesList returns a list of files from drive ordered by dislikes
func (l *LeakService) FilesList() ([]models.LeakData, error) {
	files, err := l.repo.FilesList()
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
			Size:        f.Size,
		}
		dislikes, err := strconv.Atoi(f.Properties["dislikes"])
		if err != nil {
			logg.Error(fmt.Sprint("Couldn't convert dislikes to int: ", err))
			return nil, errHandler.FileListReceiveErr
		}
		likes, err := strconv.Atoi(f.Properties["likes"])
		if err != nil {
			logg.Error(fmt.Sprint("Couldn't convert likes to int: ", err))
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

func (l *LeakService) DeleteFile(fileId string) error {
	err := l.repo.DeleteFile(fileId)
	if err != nil {
		return errors.New(fmt.Sprint("Couldn't delete file, err: ", err))
	}
	return nil
}

// File retrieves a specific file from the repository.
func (l *LeakService) File(fileID string) (models.LeakData, error) {
	// Check if the file exists
	if fileID == "" {
		logg.Error("File not found, id is empty")
		return models.LeakData{}, errHandler.FileNotFoundErr
	}
	// Retrieve the file from the repository
	b, fileData, err := l.repo.File(fileID)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get file: ", err))
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
		UserData: &models.UserFileData{},
		Subject:  &models.SubjectData{},
	}
	return leakData, nil
}
