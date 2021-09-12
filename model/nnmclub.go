package model

import (
	"time"

	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model
	Title     string
	Author    string
	Published time.Time
	Magnet    string
	Likes     int64
}
