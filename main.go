package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Season struct {
	Year int64
	Wiki string
}

func getSeason(s *goquery.Selection) *Season {
	var season = new(Season)
	txt := s.Text()
	if year, err := strconv.ParseInt(txt, 10, 32); err == nil {
		season.Year = year
	}
	if ref, ok := s.Attr("href"); ok {
		season.Wiki = ref
	}
	return season
}

func getATPSeasons() map[int64]*Season {
	seasons := make(map[int64]*Season)

	res, err := http.Get("https://en.wikipedia.org/wiki/Template:ATP_seasons")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		txt := s.Text()
		_, err := strconv.ParseInt(txt, 10, 64)
		return err == nil
	}).Each(func(i int, s *goquery.Selection) {
		season := getSeason(s)
		seasons[season.Year] = season
	})
	return seasons
}

func getWTASeasons() map[int64]*Season {
	seasons := make(map[int64]*Season)

	res, err := http.Get("https://en.wikipedia.org/wiki/Template:WTA_seasons")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		txt := s.Text()
		_, err := strconv.ParseInt(txt, 10, 64)
		return err == nil
	}).Each(func(i int, s *goquery.Selection) {
		season := getSeason(s)
		seasons[season.Year] = season
	})
	return seasons
}

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
