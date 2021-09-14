package tasks

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/storage"
)

func UpdateTopics() {
	s, err := storage.NewSqliteStorage()
	if err != nil {
		fmt.Printf("Storage: %s", err.Error())
	}
	defer s.Close()
	for i := 1; i < 200; i++ {
		processTopicPage(s, Music, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(s, HDMusic, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(s, MusicCollections, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(s, AnimeAndManga, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(s, BooksAndMediaMaterials, i)
		time.Sleep(RandDuration(32, 128))
	}
}

func processTopicPage(s *storage.SqliteStorage, catID NnmClubCategory, page int) {
	topics := GetTopics(catID, page)
	for key, topic := range topics {
		fmt.Println("ID: ", key)
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