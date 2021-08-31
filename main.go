package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
"wiki"

	"github.com/PuerkitoBio/goquery"
)

const nnmbooks = "https://nnmclub.to/forum/portal.php?c=5"
const nnmmusic = "https://nnmclub.to/forum/portal.php?c=12"


func main() {
	atps := getATPSeasons()
	for key, value := range atps {
		fmt.Println("Year: ", key, "Value: ", value)
	}
	wtas := getWTASeasons()
	for key, value := range wtas {
		fmt.Println("Year: ", key, "Value: ", value)
	}
}
