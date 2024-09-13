package models

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);unique"`
	Desc string `gorm:"type:varchar(255);unique"`
	Jam  string `gorm:"type:varchar(255);unique"`
}
