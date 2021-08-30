package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func main() {
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
		txt := s.Text()
		fmt.Printf("ATP season %d: %s\n\n", i, txt)
	})
}
