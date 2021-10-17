package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dreamerminsk/gowiki/log"
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

	doc.Find("table.pline").FilterFunction(func(i int, s *goquery.Selection) bool {
		return isTopic(s)
	}).Each(func(i int, s *goquery.Selection) {
		topic := getTopic(s)
		topics = append(topics, topic)
	})
	return topics, nil
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
