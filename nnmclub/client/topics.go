package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/metrics"
	"github.com/dreamerminsk/gowiki/nnmclub/model"
	"github.com/dreamerminsk/gowiki/utils"
	"github.com/dreamerminsk/gowiki/web"
)

func GetTopics(ctx context.Context, catID Category, page int) ([]*model.Topic, error) {
	topics := make([]*model.Topic, 0)

	doc, err := web.New().GetDocument(ctx, GetTopicsByCatUrl(catID.EnumIndex(), page))
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}

	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		metrics.GetOrRegisterValues("Web.Res", nil).Add("Title", s.Text())
	})

	doc.Find("table.pline").FilterFunction(func(i int, s *goquery.Selection) bool {
		return isTopic(s)
	}).Each(func(i int, s *goquery.Selection) {
		topic := getTopic(s)
		topics = append(topics, topic)
	})
	return topics, nil
}

func GoGetTopics(ctx context.Context, catID Category) chan model.Topic {
	topics := make(chan model.Topic, 20)

	go func() {
		url := GetTopicsByCatUrl(catID.EnumIndex(), 1)
		for {
			doc, err := web.New().GetDocument(ctx, url)
			if err != nil {
				log.Log(fmt.Sprintf("%s", err))
				close(topics)
			}

			doc.Find("table.pline").FilterFunction(func(i int, s *goquery.Selection) bool {
				return isTopic(s)
			}).Each(func(i int, s *goquery.Selection) {
				topic := getTopic(s)
				topics <- *topic
			})

			url = ""
			doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
				return strings.HasPrefix(s.Text(), "След.")
			}).Each(func(i int, s *goquery.Selection) {
				url = "https://nnmclub.to/forum/portal.php?c=12&start=20#pagestart"
			})
			if url == "" {
				close(topics)
				return
			}
		}
	}()

	return topics
}

func isTopic(s *goquery.Selection) bool {
	isTopic := false
	s.Find("a").Each(func(i int, sl *goquery.Selection) {
		if ref, ok := sl.Attr("href"); ok {
			if strings.HasPrefix(ref, "magnet:") {
				isTopic = true
			}
		}
	})
	return isTopic
}

func getTopic(s *goquery.Selection) *model.Topic {
	var topic = new(model.Topic)
	s.Find("td.pcatHead a").Each(func(i int, sl *goquery.Selection) {
		if title, ok := sl.Attr("title"); ok {
			titleString := title
			topic.Title = titleString
		}
		if href, ok := sl.Attr("href"); ok {
			u, _ := url.Parse(href)
			m, _ := url.ParseQuery(u.RawQuery)
			topicID, _ := strconv.ParseInt(m["t"][0], 10, 32)
			topic.ID = uint(topicID)
		}
	})
	s.Find("tbody > tr:nth-child(2) > td > span.genmed > b").Each(func(i int, sl *goquery.Selection) {
		authorString := sl.Text()
		topic.Author = authorString
	})
	s.Find("tbody > tr:nth-child(2) > td > span.genmed").Each(func(i int, sl *goquery.Selection) {
		timeString := sl.Text()
		topic.Published = utils.ParseTime(timeString)
	})

	s.Find("a.pcomm").Each(func(i int, sl *goquery.Selection) {
		if alt, ok := sl.Attr("title"); ok {
			if strings.Contains(alt, "Ответов") {
				topic.Comments, _ = strconv.ParseInt(strings.TrimSpace(sl.Text()), 10, 64)
			}
		}
	})

	s.Find("span.pcomm").Each(func(i int, sl *goquery.Selection) {
		if strings.HasSuffix(sl.Text(), "GB") ||
			strings.HasSuffix(sl.Text(), "MB") ||
			strings.HasSuffix(sl.Text(), "KB") {
			topic.Size = strings.TrimSpace(sl.Text())
		}
	})

	s.Find("span.pcomm").Each(func(i int, sl *goquery.Selection) {
		if _, ok := sl.Attr("id"); ok {
			topic.Likes, _ = strconv.ParseInt(sl.Text(), 10, 64)
		}
	})

	s.Find("a").Each(func(i int, sl *goquery.Selection) {
		if ref, ok := sl.Attr("href"); ok {
			if strings.HasPrefix(ref, "magnet:") {
				topic.Magnet = ref
			}
		}
	})
	return topic
}
