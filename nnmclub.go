package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const nnmbooks = 5
const nnmmusic = 12

type Topic struct {
	ID    int64
	Title string
}

func getTopic(s *goquery.Selection) *Topic {
	var topic = new(Topic)
	topic.ID = -1 * rand.Int63n(100000)
	s.Find("td.pcatHead a").Each(func(i int, sl *goquery.Selection) {
		if title, ok := sl.Attr("title"); ok {
			topic.Title = title
		}
		if href, ok := sl.Attr("href"); ok {
			u, _ := url.Parse(href)
			m, _ := url.ParseQuery(u.RawQuery)
			topic.ID, _ = strconv.ParseInt(m["t"][0], 10, 64)
		}
	})
	return topic
}

func getTopics(topicID int64) map[int64]*Topic {
	topics := make(map[int64]*Topic)

	res, err := http.Get("https://nnmclub.to/forum/portal.php?c=" + strconv.FormatInt(topicID, 10))
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

	doc.Find("table.pline").FilterFunction(func(i int, s *goquery.Selection) bool {
		isTopic := false
		s.Find("a").Each(func(i int, sl *goquery.Selection) {
			if ref, ok := sl.Attr("href"); ok {
				if strings.HasPrefix(ref, "magnet:") {
					isTopic = true
				}
			}
		})
		return isTopic
	}).Each(func(i int, s *goquery.Selection) {
		topic := getTopic(s)
		topics[topic.ID] = topic
	})
	return topics
}