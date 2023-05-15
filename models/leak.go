package models

import "os"

type LeakData struct {
	File     *File
	Subject  *SubjectData
	UserData *UserFileData
}

type File struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Size        int64    `json:"size"`
	UploadAt    int64    `json:"uploat_at_unix"`
	OpenedFile  *os.File `json:"file"`
}

type SubjectData struct {
	Faculty         string `json:"faculty"`
	Subject         string `json:"subject"`
	YearOfEducation string `json:"edu_year"`
	ModuleNum       uint64 `json:"module"`
	IsModuleTask    bool   `json:"is_module"`
	IsExam          bool   `json:"is_exam"`
}

type UserFileData struct {
	UserId   string `json:"user_id"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}
