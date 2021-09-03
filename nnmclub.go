package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const nnmbooks = 5
const nnmmusic = 12

type Topic struct {
	ID        int64
	Title     string
	Author    string
	Published time.Time
}

const timeLayout = "02 Jan 2006 15:04:05"
const timePattern = `.*?(?P<Day>\d{2}) (?P<Month>\D{3}) (?P<Year>\d{4}) (?P<Hours>\d{2}):(?P<Minutes>\d{2}):(?P<Seconds>\d{2})`

var secondsEastOfUTC = int((3 * time.Hour).Seconds())
var beijing = time.FixedZone("Beijing Time",
	secondsEastOfUTC)

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
	s.Find("tbody > tr:nth-child(2) > td > span.genmed > b").Each(func(i int, sl *goquery.Selection) {
		topic.Author = sl.Text()
	})
	s.Find("tbody > tr:nth-child(2) > td > span.genmed").Each(func(i int, sl *goquery.Selection) {
		timeString := strings.Split(sl.Text(), "|")[1]
		var compRegEx = regexp.MustCompile(timePattern)
		match := compRegEx.FindStringSubmatch(timeString)
		for i, name := range compRegEx.SubexpNames() {
			if i > 0 && i <= len(match) {
				fmt.Println(name, ": ", match[i])
			}
		}
		t := time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
		topic.Published = t
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
