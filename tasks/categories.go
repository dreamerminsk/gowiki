package tasks

import (
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
)

func NewCategories() {
	cats, err := GetCategories()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	g := storage.New()
	for _, cat := range cats {
		fmt.Println("Title: ", cat.Title)
		g.Create(cat)
	}
}
