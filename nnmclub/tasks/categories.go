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

func InitCategories(ctx context.Context, t *tasks.Task) {
	cats, err := client.GetCategories(ctx)
	if err != nil {
		log.Logf("ERROR : %s", err)
	}
	g := storage.New()
	t.MsgChan <- fmt.Sprintf("categories: %d", len(cats))
	for _, cat := range cats {
		if _, err := g.GetCategoryByID(cat.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				g.Create(cat)
				log.Logf("INSERT CATEGORY: %d - %s", cat.ID, cat.Title)
			}
		}
	}
}

func UpdateCategories() {

}
