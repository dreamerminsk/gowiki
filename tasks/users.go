package tasks

import (
	"context"

	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/storage"
	"github.com/dreamerminsk/gowiki/web/nnmclub"
)

func InitUsers(ctx context.Context) {
	g := storage.New()
	forums, err := g.GetForums()
	if err != nil {
		log.Logf("ERROR : %s", err)
	}
	for _, forum := range forums {
		users, err := nnmclub.GetForumUsers2(ctx, forum.ID)
		if err != nil {
			log.Logf("ERROR : %s", err)

		}
		for _, user := range users {
			//if _, err := g.GetUserByID(user.ID); err != nil {
			//if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := g.Create(user).Error; err != nil {

			} else {
				log.Logf("INSERT USER: %d - %s", user.ID, user.Name)
			}
			//}
			//}
		}
	}
}
