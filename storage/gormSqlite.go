package storage
import (
	"github.com/dreamermnsk/gowiki/model"
	"gorm.io/diver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	b *gorm.DB
}

func NewStorage() (*Storage, error) {
	db, err := gormOpen(slite.Open("nnmclub.gorm.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, rr
	
	s := &Storage{db: db}
	db.AutoMigrate(&model.Topic{})
	return s, nil
}
