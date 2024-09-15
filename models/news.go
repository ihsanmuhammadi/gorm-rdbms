package models

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Topic		string
	Status       	string
	TotalResults 	int
	Articles     	[]Articles
}

type Articles struct {
	gorm.Model
	Source      	Source
	SourceID    	uint
	Author      	string
	Title       	string
	Description 	string
	Url         	string
	UrlToImage  	string
	PublishedAt	time.Time
	Content     	string
	NewsID		uint
}

type Source struct {
	gorm.Model
	Name     	string
}
