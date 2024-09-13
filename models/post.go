package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title		string
	Content		string
	Comment		string
	Author		Author		`gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Category	[]Category
	Tags		[]Tag		`gorm:"many2many:post_tags;"`
}
