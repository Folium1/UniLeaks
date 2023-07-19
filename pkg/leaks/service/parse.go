package service

import (
	"strconv"

	errHandler "leaks/pkg/err"
	"leaks/pkg/models"

	"google.golang.org/api/drive/v3"
)

func parseDriveDataToModel(files []*drive.File) ([]models.Leak, error) {
	filesList := make([]models.Leak, 0, len(files))
	for _, f := range files {
		file := &models.File{
			Id:          f.Id,
			Name:        f.Name,
			Description: f.Description,
			Size:        float64(f.Size) / 1024 / 1024,
		}
		dislikesStr, ok := f.Properties["dislikes"]
		if !ok {
			l.Error("Couldn't find dislikes property")
			return nil, errHandler.FileListReceiveErr
		}
		dislikes, err := strconv.Atoi(dislikesStr)
		if err != nil {
			l.Error("Couldn't convert dislikes to int")
			return nil, errHandler.FileListReceiveErr
		}
		likes, err := strconv.Atoi(f.Properties["likes"])
		if err != nil {
			l.Error("Couldn't convert likes to int")
			return nil, errHandler.FileListReceiveErr
		}
		user := &models.UserFileData{
			Dislikes: dislikes,
			Likes:    likes,
		}
		leak := models.Leak{
			File:    file,
			User:    user,
			Subject: &models.Subject{},
		}
		filesList = append(filesList, leak)
	}
	return filesList, nil
}
