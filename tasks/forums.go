package tasks

import (
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
)

func UpdateForums() {
	forums, err := GetForums()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	g, err := storage.NewStorage()
	if err != nil {
		fmt.Printf("Storage: %s", err.Error())
	}
	for _, forum := range forums {
		fmt.Println("Title: ", forum.Title)
		g.Create(forum)
	}
}
