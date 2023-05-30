package models

type User struct {
	ID       int    `gorm:"primary_key" json:"id"`
	LeaksNum int    `gorm:"type:int" json:"leakNum"`
	Email    string `gorm:"type:varchar(100);unique" json:"email"`
	Password string `gorm:"type:varchar(200)" json:"password"`
	IsBanned bool   `gorm:"type:bool" json:"banned"`
	IsAdmin  bool   `gorm:"type:bool" json:"admin"`
}
