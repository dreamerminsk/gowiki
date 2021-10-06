package storage

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/dreamerminsk/gowiki/model"
)

type storage struct {
	DB *gorm.DB
	mu *sync.Mutex
}

type Storage interface {
	Create(value interface{}) (tx *gorm.DB)
	GetCategoryByID(ID uint) (*model.Category, error)
	GetForumByID(ID uint) (*model.Forum, error)
	GetForums() ([]*model.Forum, error)
	UpdateForum(*model.Forum) error
	GetUserByID(ID uint) (*model.User, error)
}

var (
	instance *storage
	once     sync.Once
)

func newStorage() (*storage, error) {
	db, err := gorm.Open(sqlite.Open("nnmclub.gorm.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s := &storage{DB: db, mu: &sync.Mutex{}}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Forum{})
	db.AutoMigrate(&model.Topic{})
	return s, nil
}

func New() Storage {
	once.Do(func() {
		instance, _ = newStorage()
	})

	return instance
}

func (s *storage) Create(value interface{}) (tx *gorm.DB) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.DB.Create(value)
}

func (s *storage) GetCategoryByID(ID uint) (*model.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cat := &model.Category{}
	if err := s.DB.Model(&model.Category{}).First(&cat, ID).Error; err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *storage) GetForumByID(ID uint) (*model.Forum, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	forum := &model.Forum{}
	if err := s.DB.Model(&model.Forum{}).First(&forum, ID).Error; err != nil {
		return nil, err
	}
	return forum, nil
}

func (s *storage) GetForums() ([]*model.Forum, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var forums []*model.Forum
	if err := s.DB.Model(&model.Forum{}).Find(&forums).Error; err != nil {
		return nil, err
	}
	return forums, nil
}

func (s *storage) UpdateForum(forum *model.Forum) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.DB.Model(&forum).Updates(
		model.Forum{Title: forum.Title, CatID: forum.CatID},
	).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) GetUserByID(ID uint) (*model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user := &model.User{}
	if err := s.DB.Model(&model.User{}).First(&user, ID).Error; err != nil {
		return nil, err
	}
	return user, nil
}
