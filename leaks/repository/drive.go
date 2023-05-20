package repository

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"uniLeaks/config"
	"uniLeaks/models"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// Constants for the Google Drive file properties keys
const (
	subjStr      = "subject"
	facultyStr   = "faculty"
	yearStr      = "year"
	moduleNumStr = "moduleNum"
	fileSize     = "size"
	isModuleStr  = "is_module"
	isExamStr    = "is_exam"
	likeStr      = "likes"
	dislikeStr   = "dislikes"
)

type Repo struct {
	Service *drive.Service
}

// New creates a new Repo instance using the Google Drive API token
func New() Repo {
	dr, err := config.NewDriveClient()
	if err != nil {
		log.Println(err)
	}
	return Repo{dr}
}

// SaveFile is a function that saves a new file to Google Drive with the provided LeakData.
func (r Repo) SaveFile(data models.LeakData) error {
	driveFile := &drive.File{
		Name: data.File.OpenedFile.Name(),
		Properties: map[string]string{
			"userId": data.UserData.UserId,

			subjStr:      data.Subject.Subject,
			facultyStr:   data.Subject.Faculty,
			yearStr:      fmt.Sprintf("%v", data.Subject.YearOfEducation),
			moduleNumStr: fmt.Sprintf("%v", data.Subject.ModuleNum),
			isExamStr:    strconv.FormatBool(data.Subject.IsExam),
			isModuleStr:  strconv.FormatBool(data.Subject.IsModuleTask),

			likeStr:    "0",
			dislikeStr: "0",
		},
		Description: data.File.Description,
	}
	_, err := r.Service.Files.Create(driveFile).Media(data.File.OpenedFile).Do()
	if err != nil {
		return err
	}
	log.Println("Done")
	data.File.OpenedFile.Close()
	return nil
}

// GetList is a function that returns a list of files from Google Drive that match the provided LeakData.
func (r Repo) GetList(data models.SubjectData) ([]*drive.File, error) {
	// Build the query for searching files
	query := fmt.Sprintf("properties has {key='%s' and value='%s'} and properties has {key='%s' and value='%s'} and properties has {key='%s' and value='%s'} and properties has {key='%s' and value='%s'} and properties has {key='%s' and value='%s'} and properties has {key='%s' and value='%s'}",
		facultyStr, data.Faculty,
		yearStr, data.YearOfEducation,
		subjStr, data.Subject,
		moduleNumStr, fmt.Sprintf("%v", data.ModuleNum),
		isModuleStr, strconv.FormatBool(data.IsExam),
		isExamStr, strconv.FormatBool(data.IsExam),
	)
	// Define the fields that should be returned in the file list
	fields := "files(id, name, description, size, properties)"
	files, err := r.Service.Files.List().Q(query).Fields(googleapi.Field(fields)).Do()
	if err != nil {
		return nil, err
	}
	return files.Files, nil
}

// GetFile returns file by it's id
func (r Repo) GetFile(fileId string) ([]byte, error) {
	file, err := r.Service.Files.Get(fileId).Download()
	if err != nil {
		return nil, err
	}
	reader, err := ioutil.ReadAll(file.Body)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

// LikeFile increments the number of likes for particular file
func (r Repo) LikeFile(fileId string) error {
	// Get the current like count from the file's properties field
	file, err := r.Service.Files.Get(fileId).Fields("properties").Do()
	if err != nil {
		return err
	}
	likeStr := file.Properties[likeStr]
	likeCount, _ := strconv.Atoi(likeStr)

	// Increment the like count and update the properties field of the file
	likeCount++
	properties := map[string]string{
		likeStr: strconv.Itoa(likeCount),
	}
	update := drive.File{
		Properties: properties,
	}
	_, err = r.Service.Files.Update(fileId, &update).Do()
	if err != nil {
		return err
	}
	return nil
}

// LikeFile increments the number of dislikes for particular file
func (r Repo) DislikeFile(fileId string) error {
	// Get the current dislike count from the file's properties field
	file, err := r.Service.Files.Get(fileId).Fields("properties").Do()
	if err != nil {
		return err
	}
	dislikesStr := file.Properties[dislikeStr]
	dislikes, _ := strconv.Atoi(dislikesStr)

	// Increment the dislike count and update the properties field of the file
	dislikes++
	properties := map[string]string{
		dislikesStr: strconv.Itoa(dislikes),
	}
	update := drive.File{
		Properties: properties,
	}
	_, err = r.Service.Files.Update(fileId, &update).Do()
	if err != nil {
		return err
	}
	return nil
}
