package models

type LeakData struct {
	File     *File
	Subject  *SubjectData
	UserData *UserFileData
}

type File struct {
	Id              string `json:"Id"`
	FilePath        string `json:"path"`
	FileDescription string `json:"description"`
	FileSize        int64  `json:"size"`
	UploadAt        int64  `json:"uploat_at_unix"`
}

type SubjectData struct {
	Faculty         string `json:"faculty"`
	Subject         string `json:"subject"`
	YearOfEducation int    `json:"edu_year"`
	ModuleNum       uint   `json:"module"`
	IsModuleTask    bool   `json:"is_module"`
	IsExam          bool   `json:"is_exam"`
}

type UserFileData struct {
	FileId   string `json:"Id"`
	UserId   string `json:"user_id"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}
