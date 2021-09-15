package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Title string
}

type Forum struct {
	gorm.Model
	categoryID Category
	Title      string
}

type Topic struct {
	gorm.Model
	Title     string
	Author    string
	Published time.Time
	Magnet    string
	Likes     int64
}
