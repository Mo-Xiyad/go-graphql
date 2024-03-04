package model

import (
	"gorm.io/gorm"
)

// User corresponds to a user in the database.
type Company struct {
	gorm.Model
	Name           string
	Email          string
	PhoneNumber    string
	OrganizationID string
}
