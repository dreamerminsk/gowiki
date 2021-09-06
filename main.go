package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func wiki() {
	atps := getATPSeasons()
	for key, value := range atps {
		fmt.Println("Year: ", key, "Value: ", value)
	}
	wtas := getWTASeasons()
	for key, value := range wtas {
		fmt.Println("Year: ", key, "Value: ", value)
	}
}

func main() {
	topics := getTopics(nnmmusic)
	for key, topic := range topics {
		fmt.Println("ID: ", key)
		fmt.Println("Topic: ", topic.Title)
		fmt.Println("Author: ", topic.Author)
		fmt.Println("Published: ", topic.Published)
		fmt.Println("-------------------------")
	}
}
