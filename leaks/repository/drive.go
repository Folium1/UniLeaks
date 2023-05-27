package repository

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
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
	semesterStr  = "semester"
	moduleNumStr = "moduleNum"
	fileSize     = "size"
	isModuleStr  = "is_module"
	isExamStr    = "is_exam"
	likeStr      = "likes"
	dislikeStr   = "dislikes"
)

type Repo struct {
	driveService *drive.Service
}

// New creates a new Repo instance using the Google Drive API token
func New() *Repo {
	dr, err := config.NewDriveClient()
	if err != nil {
		log.Println(err)
	}
	return &Repo{dr}
}

// SaveFile is a function that saves a new file to Google Drive with the provided LeakData.
func (r Repo) SaveFile(data *models.LeakData) error {
	// Create a new file with the provided data
	driveFile := &drive.File{
		Name: data.File.Name,
		Properties: map[string]string{
			"userId":     data.UserData.UserId,
			subjStr:      data.Subject.Subject,
			facultyStr:   data.Subject.Faculty,
			yearStr:      fmt.Sprintf("%v", data.Subject.YearOfEducation),
			semesterStr:  fmt.Sprintf("%v", data.Subject.Semester),
			moduleNumStr: fmt.Sprintf("%v", data.Subject.ModuleNum),
			isExamStr:    strconv.FormatBool(data.Subject.IsExam),
			isModuleStr:  strconv.FormatBool(data.Subject.IsModuleTask),

			likeStr:    "0",
			dislikeStr: "0",
		},
		Description: data.File.Description,
	}
	// Create an io.Reader from the file content
	contentReader := bytes.NewReader(data.File.Content)
	_, err := r.driveService.Files.Create(driveFile).Media(contentReader).Do()
	if err != nil {
		return err
	}
	data.File.Content = nil
	runtime.GC()

	log.Println("Done")
	return nil
}

// buildQuery creates a query for searching files in Google Drive.
func buildQuery(data models.SubjectData) string {
	builder := strings.Builder{}
	builder.Write([]byte(" "))
	builder.Write([]byte(fmt.Sprintf("properties has {key='%s' and value='%s'}", facultyStr, data.Faculty)))
	builder.Write([]byte(fmt.Sprintf("properties has {key='%s' and value='%s'}", yearStr, data.YearOfEducation)))
	builder.Write([]byte(fmt.Sprintf("properties has {key='%s' and value='%s'}", semesterStr, fmt.Sprintf("%v", data.Semester))))
	builder.Write([]byte(fmt.Sprintf("properties has {key='%s' and value='%s'}", subjStr, data.Subject)))
	builder.Write([]byte(fmt.Sprintf("properties has {key='%s' and value='%s'}", moduleNumStr, fmt.Sprintf("%v", data.ModuleNum))))
	builder.Write([]byte(fmt.Sprintf("properties has {key='%s' and value='%s'}", isModuleStr, strconv.FormatBool(data.IsModuleTask))))
	builder.Write([]byte(fmt.Sprintf("properties has {key='%s' and value='%s'}", isExamStr, strconv.FormatBool(data.IsExam))))
	return builder.String()
}

// FilesList returns a list of files from Google Drive that match the provided LeakData.
func (r Repo) FilesList(data models.SubjectData) ([]*drive.File, error) {
	// Build the query for searching files
	query := buildQuery(data)
	// Define the fields that should be returned in the file list
	fields := "files(id, name, description, size, properties)"
	// Get the list of files from Google Drive
	files, err := r.driveService.Files.List().Q(query).Fields(googleapi.Field(fields)).Do()
	if err != nil {
		return nil, err
	}
	return files.Files, nil
}

// File returns file by its id
func (r Repo) File(fileID string) ([]byte, *drive.File, error) {
	// Get the file metadata
	fields := "name, description, size, properties"
	fileData, err := r.driveService.Files.Get(fileID).Fields(googleapi.Field(fields)).Do()
	if err != nil {
		return nil, nil, err
	}

	// Get the file content
	res, err := r.driveService.Files.Get(fileID).Download()
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	// Copy the file content to a byte array
	buffer := bytes.Buffer{}
	_, err = io.Copy(&buffer, res.Body)
	if err != nil {
		return nil, nil, err
	}
	return buffer.Bytes(), fileData, nil
}

// LikeFile increments the number of likes for particular file
func (r Repo) LikeFile(fileId string) error {
	// Get the current like count from the file's properties field
	file, err := r.driveService.Files.Get(fileId).Fields("properties").Do()
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
	_, err = r.driveService.Files.Update(fileId, &update).Do()
	if err != nil {
		return err
	}
	return nil
}

// LikeFile increments the number of dislikes for particular file
func (r Repo) DislikeFile(fileId string) error {
	// Get the current dislike count from the file's properties field
	file, err := r.driveService.Files.Get(fileId).Fields("properties").Do()
	if err != nil {
		return err
	}
	dislikesStr := file.Properties[dislikeStr]
	dislikes, err := strconv.Atoi(dislikesStr)
	if err != nil {
		return err
	}

	// Increment the dislike count and update the properties field of the file
	dislikes++
	properties := map[string]string{
		dislikesStr: strconv.Itoa(dislikes),
	}
	update := drive.File{
		Properties: properties,
	}
	_, err = r.driveService.Files.Update(fileId, &update).Do()
	if err != nil {
		return err
	}
	return nil
}

// AllFiles returns all files from Google Drive
func (r Repo) AllFiles() ([]*drive.File, error) {
	fields := "files(id, description, name, size, properties)"
	files, err := r.driveService.Files.List().Fields(googleapi.Field(fields)).Q("").Do()
	if err != nil {
		log.Printf("Failed to retrieve files by custom properties: %v", err)
		return nil, err
	}
	return files.Files, nil
}
