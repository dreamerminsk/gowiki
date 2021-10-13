package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/storage"
	"github.com/dreamerminsk/gowiki/web/nnmclub"
)

func InitUsers(ctx context.Context, t *Task) {
	g := storage.New()

	forums, err := g.GetForums()
	if err != nil {
		log.Logf("ERROR : %s", err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(forums), func(i, j int) {
		forums[i], forums[j] = forums[j], forums[i]
	})

	newUsers := 0
	lastUser := &model.User{}
	for idx, forum := range forums {
		users, next, err := nnmclub.GetForumUsers2(ctx, forum.ID, 1)
		if err != nil {
			log.Logf("ERROR : %s", err)

		}

		if next {
			users2, next2, err := nnmclub.GetForumUsers2(ctx, forum.ID, 2)
			if err != nil {
				log.Logf("ERROR : %s", err)
			}
			users = append(users, users2...)
		}

		if next2 {
			users3, _, err := nnmclub.GetForumUsers2(ctx, forum.ID, 2)
			if err != nil {
				log.Logf("ERROR : %s", err)
			}
			users = append(users, users3...)
		}

		for _, user := range users {
			//if _, err := g.GetUserByID(user.ID); err != nil {
			//if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := g.Create(user).Error; err != nil {

			} else {
				newUsers++
				lastUser = user
				log.Logf("INSERT USER: %d - %s", user.ID, user.Name)
			}
			//}
			//}
			t.MsgChan <- fmt.Sprintf("forums: %d of %d, new users: %d, last user: %d, %s",
				idx+1, len(forums), newUsers, lastUser.ID, lastUser.Name)
		}
	}
}
