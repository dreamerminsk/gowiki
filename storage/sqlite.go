package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/dreamerminsk/gowiki/model"
)

type Storage struct {
	DB *gorm.DB
}

func NewStorage() (*Storage, error) {
	db, err := gorm.Open(sqlite.Open("nnmclub.gorm.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s := &Storage{DB: db}
	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Forum{})
	db.AutoMigrate(&model.Topic{})
	return s, nil
}

func (s *Storage) Create(value interface{}) (tx *gorm.DB) {
return s.DB.Create(value)
}
