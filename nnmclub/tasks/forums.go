package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/nnmclub/client"
	"github.com/dreamerminsk/gowiki/nnmclub/storage"
	"github.com/dreamerminsk/gowiki/tasks"
	"gorm.io/gorm"
)

func InitForums(ctx context.Context, t *tasks.Task) {
	forums, err := client.GetForums(ctx)
	if err != nil {
		log.Logf("ERROR : %s", err)
	}
	g := storage.New()
	for _, forum := range forums {
		if _, err := g.GetForumByID(forum.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				g.Create(forum)
				log.Logf("INSERT FORUM: %d - %s", forum.ID, forum.Title)
			}
		}
	}
}

func UpdateForums(ctx context.Context, t *tasks.Task) {
	g := storage.New()

	forums, err := g.GetForums()
	if err != nil {
		log.Logf("ERROR : %s", err)
	}

	t.MsgChan <- fmt.Sprintf("forums: %d", len(forums))

	for _, forum := range forums {
		if forum.CatID == 0 {
			newForum, err := client.GetForum(ctx, forum.ID)
			if err != nil {
				log.Logf("ERROR : %s", err)
				continue
			}
			g.UpdateForum(newForum)
		}
	}
}
