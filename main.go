package main

import (
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	s, err := NewStorage()
	if err != nil {
		fmt.Printf("Storage: %s", err.Error())
	}
	defer s.Close()
	for i := 5; i < 200; i++ {
		processTopicPage(s, Music, i)
		time.Sleep(75 * time.Second)
	}

}

func processTopicPage(s *Storage, catID NnmClubCategory, page int) {
	topics := getTopics(Music, page)
	for key, topic := range topics {
		fmt.Println("ID: ", key)
		fmt.Println("Title: ", topic.Title)
		fmt.Println("Author: ", topic.Author)
		fmt.Println("Published: ", topic.Published)
		fmt.Println("Likes: ", topic.Likes)
		fmt.Println("Magnet: ", topic.Magnet)
		err := s.addTopic(topic)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		fmt.Println("-------------------------")
	}
}
