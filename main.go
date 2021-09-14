package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/storage"
	"github.com/dreamerminsk/gowiki/tasks"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	g, err := storage.NewStorage()
	if err != nil {
		fmt.Printf("Storage: %s", err.Error())
	}

	cats, err := tasks.GetCategories()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	for _, cat := range cats {
		fmt.Println("Title: ", cat.Title)
		g.DB.Create(cat)
	}

	forums, err := tasks.GetForums()
	if err != nil {
		fmt.Println("ERROR : ", err)
	}
	for _, forum := range forums {
		fmt.Println("Title: ", forum.Title)
		g.DB.Create(forum)
	}

	s, err := NewStorage()
	if err != nil {
		fmt.Printf("Storage: %s", err.Error())
	}
	defer s.Close()
	for i := 1; i < 200; i++ {
		processTopicPage(s, tasks.Music, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(s, tasks.HDMusic, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(s, tasks.MusicCollections, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(s, tasks.AnimeAndManga, i)
		time.Sleep(RandDuration(32, 128))
		processTopicPage(s, tasks.BooksAndMediaMaterials, i)
		time.Sleep(RandDuration(32, 128))
	}

}

func processTopicPage(s *SqliteStorage, catID tasks.NnmClubCategory, page int) {
	topics := tasks.GetTopics(catID, page)
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

func insertOrUpdate(s *SqliteStorage, topic *model.Topic) {
	oldTopic, selectErr := s.getTopic(int(topic.ID))
	if selectErr != nil {
		fmt.Println("SELECT ERROR: ", reflect.TypeOf(selectErr), selectErr)
	}
	if oldTopic.ID == 0 {
		insertErr := s.addTopic(topic)
		if insertErr != nil {
			fmt.Println("INSERT ERROR: ", reflect.TypeOf(insertErr), insertErr)
		}
	} else if topic.Likes > oldTopic.Likes {
		fmt.Printf("\tDIFF LIKES: %d\r\n", topic.Likes-oldTopic.Likes)
		updateErr := s.updateTopic(topic)
		if updateErr != nil {
			fmt.Println("UPDATE ERROR: ", reflect.TypeOf(updateErr), updateErr)
		}
	}
}

func RandDuration(min, max int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Second * time.Duration((r.Intn(max-min) + min))
}
