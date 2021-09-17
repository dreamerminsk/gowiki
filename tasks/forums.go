package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
	"gorm.io/gorm"
)

func InitForums(ctx context.Context) {
	forums, err := GetForums(ctx)
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	g := storage.New()
	for _, forum := range forums {
		if _, err := g.GetForumByID(forum.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				g.Create(forum)
				fmt.Println("INSERT FORUM: ", forum.ID, " - ", forum.Title)
			}
		}
	}
}

func UpdateForums() {}
