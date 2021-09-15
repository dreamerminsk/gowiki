package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Title  string
	Forums []Forum `gorm:"foreignKey:CatID"`
}

type Forum struct {
	gorm.Model
	CatID Category
	Title string
}

type Topic struct {
	gorm.Model
	Title     string
	Author    string
	Published time.Time
	Magnet    string
	Likes     int64
}
