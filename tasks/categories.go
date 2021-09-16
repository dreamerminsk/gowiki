package tasks

import (
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
)

func InitOrUpdateCategories() {
	cats, err := GetCategories()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	g := storage.New()
	for _, cat := range cats {
		if _, err := g.GetCategoryByID(cat.ID); err != nil {
			g.Create(cat)
			fmt.Println("Title: ", cat.Title)
		}
	}
}

func UpdateCategories() {

}
