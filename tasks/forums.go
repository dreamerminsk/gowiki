package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/dreamerminsk/gowiki/storage"
	"github.com/dreamerminsk/gowiki/web"
	"gorm.io/gorm"
)

func InitForums(ctx context.Context) {
	forums, err := web.GetForums(ctx)
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

func UpdateForums(ctx context.Context) {
	g := storage.New()
	forums, err := g.GetForums()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	for _, forum := range forums {
		fmt.Println("FORUM : ", forum)
	}
}
