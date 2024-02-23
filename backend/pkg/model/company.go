package model

import (
	"gorm.io/gorm"
)

// User corresponds to a user in the database.
type CompanyModel struct {
	gorm.Model
	Name   string
	Domain string `gorm:"unique"`
}
