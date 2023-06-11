package models

type User struct {
	ID       int    `gorm:"primary_key" json:"id"`
	NickName string `gorm:"type:varchar(20);unique" json:"nickname"`
	Email    string `gorm:"type:varchar(100)" json:"email"`
	Password string `gorm:"type:varchar(200)" json:"password"`
	IsBanned bool   `gorm:"type:bool" json:"banned"`
	IsAdmin  bool   `gorm:"type:bool" json:"admin"`
}
