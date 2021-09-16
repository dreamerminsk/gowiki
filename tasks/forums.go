package tasks

import (
	"errors"
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
	"gorm.io/gorm"
)

func InitOrUpdateForums() {
	forums, err := GetForums()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	g := storage.New()
	for _, forum := range forums {
		if _, err := g.GetCategoryByID(forum.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				g.Create(forum)
				fmt.Println("INSERT FORUM: ", forum.ID, " - ", forum.Title)
			}
		}
	}
}

func UpdateForums() {}
