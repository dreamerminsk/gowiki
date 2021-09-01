package main

import (
	"fmt"
	//"log"
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
		fmt.Println("ID: ", key, "Topic: ", topic.Title)
	}
}
