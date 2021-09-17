package tasks

import (
	"errors"
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
	"gorm.io/gorm"
)

func InitCategories() {
	cats, err := GetCategories()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	g := storage.New()
	for _, cat := range cats {
		if _, err := g.GetCategoryByID(cat.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				g.Create(cat)
				fmt.Println("INSERT CATEGORY: ", cat.ID, " - ", cat.Title)
			}
		}
	}
}

func UpdateCategories() {

}
