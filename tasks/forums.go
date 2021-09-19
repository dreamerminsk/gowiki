package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"

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
		if forum.CatID == 0 {
		    newForum, err := web.GetForum(ctx, forum.ID)
		    if err != nil {
			    fmt.Printf("[%s] [%s] %s\r\n", time.Now().Format(time.RFC3339), "tasks->UpdateForums", err)

		    }
		    g.UpdateForum(newForum)
                }
	}
}
