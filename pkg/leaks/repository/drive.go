package repository

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"leaks/pkg/config"
	logg "leaks/pkg/logger"
	"leaks/pkg/models"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// Constants for the Google Drive file properties keys
const (
	userIdStr        = "user_id"
	subjStr          = "subject"
	facultyStr       = "faculty"
	yearStr          = "year"
	semesterStr      = "semester"
	moduleNumStr     = "moduleNum"
	isModuleStr      = "is_module"
	isExamStr        = "is_exam"
	likeStr          = "likes"
	usersLikedStr    = "users_liked"
	dislikeStr       = "dislikes"
	usersDislikedStr = "users_disliked"
)

var logger = logg.NewLogger()

type Repo struct {
	driveService *drive.Service
}

// New creates a new Repo instance using the Google Drive API token
func New() *Repo {
	dr, err := config.NewDriveClient()
	if err != nil {
		logger.Fatal(fmt.Sprint("Couldn't create a new drive client, err:", err))
	}
	return &Repo{dr}
}

// SaveFile is a function that saves a new file to Google Drive with the provided LeakData.
func (r *Repo) SaveFile(data models.Leak) error {
	// Create a new file with the provided data
	driveFile := &drive.File{
		Name: data.File.Name,
		Properties: map[string]string{
			subjStr:      data.Subject.Subject,
			facultyStr:   data.Subject.Faculty,
			yearStr:      fmt.Sprintf("%v", data.Subject.YearOfEducation),
			semesterStr:  fmt.Sprintf("%v", data.Subject.Semester),
			moduleNumStr: fmt.Sprintf("%v", data.Subject.ModuleNum),
			isExamStr:    strconv.FormatBool(data.Subject.IsExam),
			isModuleStr:  strconv.FormatBool(data.Subject.IsModuleTask),

			userIdStr:        data.User.UserId,
			likeStr:          "0",
			usersLikedStr:    "",
			dislikeStr:       "0",
			usersDislikedStr: "",
		},
		Description: data.File.Description,
	}
	// Create an io.Reader from the file content
	contentReader := bytes.NewReader(data.File.Content)
	_, err := r.driveService.Files.Create(driveFile).Media(contentReader).Do()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't create a new file, err:", err))
		return err
	}
	data.File.Content = nil
	return nil
}

// buildQuery creates a query for searching files in Google Drive.
func buildQuery(data models.Subject) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("properties has {key='%s' and value='%s'} and ", facultyStr, data.Faculty))
	builder.WriteString(fmt.Sprintf("properties has {key='%s' and value='%s'} and ", yearStr, data.YearOfEducation))
	builder.WriteString(fmt.Sprintf("properties has {key='%s' and value='%v'} and ", semesterStr, data.Semester))
	builder.WriteString(fmt.Sprintf("properties has {key='%s' and value='%s'} and ", subjStr, data.Subject))
	builder.WriteString(fmt.Sprintf("properties has {key='%s' and value='%v'} and ", moduleNumStr, data.ModuleNum))
	builder.WriteString(fmt.Sprintf("properties has {key='%s' and value='%t'} and ", isModuleStr, data.IsModuleTask))
	builder.WriteString(fmt.Sprintf("properties has {key='%s' and value='%t'}", isExamStr, data.IsExam))
	return builder.String()
}

// FilesList returns a list of files from Google Drive that match the provided LeakData.
func (r *Repo) FilesList(data models.Subject) ([]*drive.File, error) {
	// Build the query for searching files
	query := buildQuery(data)
	// Get the list of files from Google Drive
	files, err := r.driveService.Files.List().Q(query).Fields("files(id, name, description, size, properties, createdTime)").OrderBy("name").Do()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return nil, err
	}
	return files.Files, nil
}

func (r *Repo) DownloadFile(fileID string) ([]byte, drive.File, error) {
	fileData, err := r.driveService.Files.Get(fileID).Fields("name, description, size, properties").Do()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get file data, err:", err))
		return nil, drive.File{}, err
	}
	res, err := r.driveService.Files.Get(fileID).Download()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get file content, err:", err))
		return nil, drive.File{}, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	buffer := bytes.Buffer{}
	_, err = io.Copy(&buffer, res.Body)
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't copy file content, err:", err))
		return nil, drive.File{}, err
	}
	return buffer.Bytes(), *fileData, nil
}

func (r *Repo) LikeFile(fileID, userID string) error {
	file, err := r.getFileProperties(fileID)
	if err != nil {
		return err
	}

	if userHasAction(file.Properties[usersLikedStr], userID) {
		likeCount, err := strconv.Atoi(file.Properties[likeStr])
		if err != nil {
			return err
		}
		likeCount--

		properties := map[string]string{
			usersLikedStr: strings.ReplaceAll(file.Properties[usersLikedStr], userID+",", ""),
			likeStr:       strconv.Itoa(likeCount),
		}
		if err := r.updateFileProperties(fileID, properties); err != nil {
			return err
		}
	} else if userHasAction(file.Properties[usersDislikedStr], userID) {
		return errors.New("user has already disliked the file")
	} else {
		likeCount, err := strconv.Atoi(file.Properties[likeStr])
		if err != nil {
			return err
		}
		likeCount++

		properties := map[string]string{
			usersLikedStr: file.Properties[usersLikedStr] + userID + ",",
			likeStr:       strconv.Itoa(likeCount),
		}
		if err := r.updateFileProperties(fileID, properties); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) DislikeFile(fileID, userID string) error {
	file, err := r.getFileProperties(fileID)
	if err != nil {
		return err
	}

	if userHasAction(file.Properties[usersLikedStr], userID) {
		return errors.New("user has already liked the file")
	} else if userHasAction(file.Properties[usersDislikedStr], userID) {
		// Decrement the dislike count
		dislikeCount, err := strconv.Atoi(file.Properties[dislikeStr])
		if err != nil {
			return err
		}
		dislikeCount--

		properties := map[string]string{
			usersDislikedStr: strings.ReplaceAll(file.Properties[usersDislikedStr], userID+",", ""),
			dislikeStr:       strconv.Itoa(dislikeCount),
		}
		if err := r.updateFileProperties(fileID, properties); err != nil {
			return err
		}
	} else {

		dislikeCount, err := strconv.Atoi(file.Properties[dislikeStr])
		if err != nil {
			return err
		}
		dislikeCount++

		properties := map[string]string{
			usersDislikedStr: file.Properties[usersDislikedStr] + userID + ",",
			dislikeStr:       strconv.Itoa(dislikeCount),
		}
		if err := r.updateFileProperties(fileID, properties); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) getFileProperties(fileID string) (*drive.File, error) {
	file, err := r.driveService.Files.Get(fileID).Fields("properties").Do()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get file data, err:", err))
		return nil, err
	}
	return file, nil
}

func (r *Repo) updateFileProperties(fileID string, properties map[string]string) error {
	update := drive.File{
		Properties: properties,
	}
	_, err := r.driveService.Files.Update(fileID, &update).Do()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't update file data, err:", err))
		return err
	}
	return nil
}

func userHasAction(actions, userID string) bool {
	return strings.Contains(actions, userID)
}

func (r *Repo) AllFiles() ([]*drive.File, error) {
	fields := "files(id, description, name, size, properties)"
	files, err := r.driveService.Files.List().Fields(googleapi.Field(fields)).Q("").OrderBy("createdTime desc").Do()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return nil, err
	}
	return files.Files, nil
}

func (r *Repo) UserFiles(userId string) ([]*drive.File, error) {
	fields := "files(id, description, name, size, properties)"
	files, err := r.driveService.Files.List().Fields(googleapi.Field(fields)).Q(fmt.Sprintf("properties has {key='%s' and value='%s'}", userIdStr, userId)).Do()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't get files list, err:", err))
		return nil, err
	}
	return files.Files, nil
}

func (r *Repo) DeleteFile(fileId string) error {
	err := r.driveService.Files.Delete(fileId).Do()
	if err != nil {
		logger.Error(fmt.Sprint("Couldn't delete file, err:", err))
		return err
	}
	return nil
}
