package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `gorm:"-"`
	ID         uint64 `gorm:"primaryKey"`
	Name       string
	Email      string `gorm:"unique"`
}

// 1 =  before generating the code for gql run the following
// go get github.com/99designs/gqlgen@latest

// 2 = and then run
// go run github.com/99designs/gqlgen generate

// gqlgen generate
