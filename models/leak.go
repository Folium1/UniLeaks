package models

type LeakData struct {
	File     *File
	Subject  *SubjectData
	UserData *UserFileData
}

type File struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Size        float64 `json:"size"`
	Content     []byte  `json:"file"`
}

type SubjectData struct {
	Faculty         string `json:"faculty"`
	Subject         string `json:"subject"`
	YearOfEducation string `json:"edu_year"`
	Semester        uint64 `json:"semester"`
	ModuleNum       uint64 `json:"module"`
	IsModuleTask    bool   `json:"is_module"`
	IsExam          bool   `json:"is_exam"`
}

type UserFileData struct {
	UserId   string `json:"user_id"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}

type LikeDislikeData struct {
	UserId string `json:"userId"`
	FileId string `json:"fileId"`
	Action string `json:"action"`
}
