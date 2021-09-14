package tasks

import (
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
)

func UpdateCategories() {
	cats, err := GetCategories()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	g, err := storage.NewStorage()
	if err != nil {
		fmt.Printf("Storage: %s", err.Error())
	}
	for _, cat := range cats {
		fmt.Println("Title: ", cat.Title)
		g.Create(cat)
	}
}
