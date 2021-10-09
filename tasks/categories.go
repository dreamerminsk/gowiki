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

func InitCategories(ctx context.Context, t *Task) {
	cats, err := nnmclub.GetCategories(ctx)
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
