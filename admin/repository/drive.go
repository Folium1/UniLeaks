package repository

import (
	"bytes"
	"fmt"
	"io"
	"leaks/config"
	"leaks/logger"

	"google.golang.org/api/drive/v3"
)

var logg = logger.NewLogger()

// DriveRepo is a repository for the Google Drive API
type DriveRepo struct {
	drive *drive.Service
}

// New creates a new instance of the repository.
func NewDriveRepo() *DriveRepo {
	dr, err := config.NewDriveClient()
	if err != nil {
		logg.Fatal(fmt.Sprint("Couldn't create a new drive client, err:", err))
	}
	return &DriveRepo{dr}
}

// DeleteFile deletes file from drive
func (r *DriveRepo) DeleteFile(fileId string) error {
	return r.drive.Files.Delete(fileId).Do()
}

// FilesList returns a list of files from drive ordered by dislikes
func (r *DriveRepo) FilesList() ([]*drive.File, error) {
	files, err := r.drive.Files.List().Fields("files(id, name, description, size, properties)").Do()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return nil, err
	}
	return files.Files, nil
}

// File returns file by its id
func (r *DriveRepo) File(fileID string) ([]byte, drive.File, error) {
	// Get the file metadata
	fileData, err := r.drive.Files.Get(fileID).Fields("name, description, size, properties").Do()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get file data, err:", err))
		return nil, drive.File{}, err
	}
	// Get the file content
	res, err := r.drive.Files.Get(fileID).Download()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get file content, err:", err))
		return nil, drive.File{}, err
	}
	defer res.Body.Close()

	// Copy the file content to a byte array
	buffer := bytes.Buffer{}
	_, err = io.Copy(&buffer, res.Body)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't copy file content to a byte array, err:", err))
		return nil, drive.File{}, err
	}
	return buffer.Bytes(), *fileData, nil
}

// GetUserFilesList returns a list of user files from drive
func (r *DriveRepo) GetUserFilesList(userId string) ([]*drive.File, error) {
	files, err := r.drive.Files.List().Q(fmt.Sprintf("properties has { key='userId' and value='%s' }", userId)).Fields("files(id, name, description, size, properties)").Do()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return nil, err
	}
	return files.Files, nil
}

// DeleteAllUserFiles deletes all user files from drive
func (r *DriveRepo) DeleteAllUserFiles(userId string) error {
	files, err := r.drive.Files.List().Fields("files(id)").Do()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return err
	}
	for _, file := range files.Files {
		if file.Properties["userId"] == userId {
			err := r.drive.Files.Delete(file.Id).Do()
			if err != nil {
				logg.Error(fmt.Sprint("Couldn't delete file, err:", err))
				return err
			}
		}
	}
	return nil
}
