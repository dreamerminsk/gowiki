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

func processTopicPage(ctx context.Context, s *storage.SqliteStorage, catID web.NnmClubCategory, page int) {
	topics := web.GetTopics(ctx, catID, page)
	for _, topic := range topics {
		fmt.Println("ID: ", topic.ID)
		fmt.Println("Title: ", topic.Title)
		fmt.Println("Author: ", topic.Author)
		fmt.Println("Published: ", topic.Published.Format(time.RFC3339))
		fmt.Println("Likes: ", topic.Likes)
		fmt.Println("Magnet: ", topic.Magnet)
		insertOrUpdate(s, topic)
		fmt.Println("-------------------------")
	}
}

func insertOrUpdate(s *storage.SqliteStorage, topic *model.Topic) {
	oldTopic, selectErr := s.GetTopic(int(topic.ID))
	if selectErr != nil {
		fmt.Println("SELECT ERROR: ", reflect.TypeOf(selectErr), selectErr)
	}
	if oldTopic.ID == 0 {
		insertErr := s.AddTopic(topic)
		if insertErr != nil {
			fmt.Println("INSERT ERROR: ", reflect.TypeOf(insertErr), insertErr)
		}
	} else if topic.Likes > oldTopic.Likes {
		fmt.Printf("\tDIFF LIKES: %d\r\n", topic.Likes-oldTopic.Likes)
		updateErr := s.UpdateTopic(topic)
		if updateErr != nil {
			fmt.Println("UPDATE ERROR: ", reflect.TypeOf(updateErr), updateErr)
		}
	}
}

func RandDuration(min, max int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Second * time.Duration((r.Intn(max-min) + min))
}
