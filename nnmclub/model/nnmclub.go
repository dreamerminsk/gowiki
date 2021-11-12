package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}

type Category struct {
	ID        uint `gorm:"primarykey"`
	Title     string
	Forums    []Forum `gorm:"foreignKey:CatID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}

type Forum struct {
	ID        uint `gorm:"primarykey"`
	CatID     uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}

type Topic struct {
	ID        uint `gorm:"primarykey"`
	Title     string
	Category  int64
	Author    string
	Published time.Time
	Magnet    string
	Likes     int64
	Comments  int64
	Size      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}
