package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/storage"
	"github.com/dreamerminsk/gowiki/web"
)

func UpdateTopics(ctx context.Context) {
	s, err := storage.NewSqliteStorage()
	if err != nil {
		fmt.Printf("Storage: %s", err.Error())
	}
	defer s.Close()
	for i := 1; i < 200; i++ {
		processTopicPage(ctx, s, web.Music, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(ctx, s, web.HDMusic, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(ctx, s, web.MusicCollections, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(ctx, s, web.AnimeAndManga, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(ctx, s, web.BooksAndMediaMaterials, i)
		time.Sleep(RandDuration(32, 128))
	}
}

func processTopicPage(ctx context.Context, s *storage.SqliteStorage, catID web.NnmClubCategory, page int) error {
	topics, err := web.GetTopics(ctx, catID, page)
        if err != nil {
            return err
        }
	for _, topic := range topics {
		fmt.Println("ID: ", topic.ID)
		fmt.Println("Title: ", topic.Title)
		fmt.Println("Author: ", topic.Author)
		fmt.Println("Published: ", topic.Published.Format(time.RFC3339))
		fmt.Println("Likes: ", topic.Likes)
		fmt.Println("Magnet: ", topic.Magnet)
		err = insertOrUpdate(s, topic)
                if err != nil {
                    return err
                }
		fmt.Println("-------------------------")
	}
}

func insertOrUpdate(s *storage.SqliteStorage, topic *model.Topic) error {
	oldTopic, err := s.GetTopic(int(topic.ID))
	if err != nil {
		fmt.Println("SELECT ERROR: ", reflect.TypeOf(err), err)
                return err
	}
	if oldTopic.ID == 0 {
		err = s.AddTopic(topic)
		if err != nil {
			fmt.Println("INSERT ERROR: ", reflect.TypeOf(err), err)
                        return err
		}
	} else if topic.Likes > oldTopic.Likes {
		fmt.Printf("\tDIFF LIKES: %d\r\n", topic.Likes-oldTopic.Likes)
		err = s.UpdateTopic(topic)
		if err != nil {
			fmt.Println("UPDATE ERROR: ", reflect.TypeOf(err), err)
                        return err
		}
	}
    return nil
}

func RandDuration(min, max int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Second * time.Duration((r.Intn(max-min) + min))
}
