package handlers

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name             string
	Surname          string
	UniversityNumber int `gorm:"unique"`
	Token            string
}
