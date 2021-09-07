package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	topics := getTopics(Music)
	for key, topic := range topics {
		fmt.Println("ID: ", key)
		fmt.Println("Title: ", topic.Title)
		fmt.Println("Author: ", topic.Author)
		fmt.Println("Published: ", topic.Published)
		fmt.Println("Likes: ", topic.Likes)
		fmt.Println("Magnet: ", topic.Magnet)
		fmt.Println("-------------------------")
	}
}
