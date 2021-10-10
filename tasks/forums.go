package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/storage"
	"github.com/dreamerminsk/gowiki/web/nnmclub"
	"gorm.io/gorm"
)

func InitForums(ctx context.Context, t *Task) {
	forums, err := nnmclub.GetForums(ctx)
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

func UpdateForums(ctx context.Context, t *Task) {
	g := storage.New()
	
	forums, err := g.GetForums()
	if err != nil {
		log.Logf("ERROR : %s", err)
	}
	
	t.MsgChan <- fmt.Sprintf("forums: %d", len(forums))
	
	for _, forum := range forums {
		if forum.CatID == 0 {
			newForum, err := nnmclub.GetForum(ctx, forum.ID)
			if err != nil {
				log.Logf("ERROR : %s", err)
continue
			}
			g.UpdateForum(newForum)
		}
	}
}
