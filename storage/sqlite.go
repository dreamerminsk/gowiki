package storage

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/dreamerminsk/gowiki/model"
)

type storage struct {
	DB *gorm.DB
}

type Storage interface {
	Create(value interface{}) (tx *gorm.DB)
}

var (
	instance *storage
	once     sync.Once
)

func NewStorage() (*storage, error) {
	db, err := gorm.Open(sqlite.Open("nnmclub.gorm.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s := &storage{DB: db}
	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Forum{})
	db.AutoMigrate(&model.Topic{})
	return s, nil
}

func New() Storage {
	once.Do(func() {
		instance, _ = NewStorage()
	})

	return instance
}

func (s *storage) Create(value interface{}) (tx *gorm.DB) {
	return s.DB.Create(value)
}
