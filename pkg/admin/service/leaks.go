package service

import (
	"errors"
	"fmt"

	adminRepo "leaks/pkg/admin"
	repo "leaks/pkg/admin/repository"
	errHandler "leaks/pkg/err"

	"leaks/pkg/logger"
	"leaks/pkg/models"

	_ "google.golang.org/api/drive/v3"
)

var l = logger.NewLogger()

type AdminLeakService struct {
	lister   adminRepo.FileLister
	deleter  adminRepo.FileDeleter
	receiver adminRepo.FileReceiver
}

func NewAdminLeakService() *AdminLeakService {
	r := repo.NewDriveRepo()
	return &AdminLeakService{
		lister:   r,
		deleter:  r,
		receiver: r,
	}
}

func (a *AdminLeakService) FilesOrderedByDislikes() ([]models.Leak, error) {
	driveFiles, err := a.lister.FilesOrderedByDislikes()
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

func (a *AdminLeakService) DeleteFile(fileId string) error {
	err := a.deleter.DeleteFile(fileId)
	if err != nil {
		return errors.New(fmt.Sprint("Couldn't delete file, err: ", err))
	}
	return nil
}

func (a *AdminLeakService) File(fileID string) (models.Leak, error) {
	b, file, err := a.receiver.File(fileID)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get file: ", err))
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
		User:    &models.UserFileData{},
		Subject: &models.Subject{},
	}

	return leak, nil
}

func (a *AdminLeakService) GetFilesListUploadedByUser(userId string) ([]models.Leak, error) {
	driveFiles, err := a.lister.GetUserFilesList(userId)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get user files list: ", err))
		return nil, errHandler.FileListReceiveErr
	}
	filesList, err := parseDriveDataToModel(driveFiles)
	if err != nil {
		return nil, err
	}
	return filesList, nil
}

func (a *AdminLeakService) DeleteAllFilesUploadedByUser(userId string) error {
	err := a.deleter.DeleteAllUserFiles(userId)
	if err != nil {

		return errors.New(fmt.Sprint("Couldn't delete user files, err: ", err))
	}
	return nil
}
