package tasks

import (
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
)

func NewForums() {
	forums, err := GetForums()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	g := storage.New()
	for _, forum := range forums {
		fmt.Println("Title: ", forum.Title)
		g.Create(forum)
	}
}

func UpdateForums() {}
