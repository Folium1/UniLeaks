package repository

import (
	"bytes"
	"fmt"
	"io"

	"leaks/pkg/config"
	logg "leaks/pkg/logger"

	"google.golang.org/api/drive/v3"
)

var l = logg.NewLogger()

type DriveRepo struct {
	drive *drive.Service
}

func NewDriveRepo() *DriveRepo {
	dr, err := config.NewDriveClient()
	if err != nil {
		l.Fatal(fmt.Sprint("Couldn't create a new drive client, err:", err))
	}
	return &DriveRepo{dr}
}

func (r *DriveRepo) DeleteFile(fileId string) error {
	return r.drive.Files.Delete(fileId).Do()
}

func (r *DriveRepo) FilesOrderedByDislikes() ([]*drive.File, error) {
	files, err := r.drive.Files.List().Fields("files(id, name, description, size, properties)").Do()
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return nil, err
	}
	return files.Files, nil
}

func (r *DriveRepo) File(fileID string) ([]byte, drive.File, error) {
	fileData, err := r.drive.Files.Get(fileID).Fields("name, description, size, properties").Do()
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get file data, err:", err))
		return nil, drive.File{}, err
	}

	res, err := r.drive.Files.Get(fileID).Download()
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get file content, err:", err))
		return nil, drive.File{}, err
	}
	defer res.Body.Close()

	buffer := bytes.Buffer{}
	_, err = io.Copy(&buffer, res.Body)
	if err != nil {
		l.Error(fmt.Sprint("Couldn't copy file content to a byte array, err:", err))
		return nil, drive.File{}, err
	}

	return buffer.Bytes(), *fileData, nil
}

func (r *DriveRepo) GetUserFilesList(userId string) ([]*drive.File, error) {
	files, err := r.drive.Files.List().Q(fmt.Sprintf("properties has { key='userId' and value='%s' }", userId)).Fields("files(id, name, description, size, properties)").Do()
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return nil, err
	}
	return files.Files, nil
}

func (r *DriveRepo) DeleteAllUserFiles(userId string) error {
	files, err := r.drive.Files.List().Fields("files(id),properties").Do()
	if err != nil {
		l.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return err
	}
	for _, file := range files.Files {
		if file.Properties["userId"] == userId {
			err := r.drive.Files.Delete(file.Id).Do()
			if err != nil {
				l.Error(fmt.Sprint("Couldn't delete file, err:", err))
				return err
			}
		}
	}
	return nil
}
