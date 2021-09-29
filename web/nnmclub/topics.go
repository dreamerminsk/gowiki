package nnmclub

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/utils"
	"github.com/dreamerminsk/gowiki/web"
	"golang.org/x/text/encoding/charmap"
)

func GetTopics(ctx context.Context, catID Category, page int) ([]*model.Topic, error) {
	topics := make([]*model.Topic, 0)
	var urlBuilder strings.Builder
	urlBuilder.WriteString("https://nnmclub.to/forum/portal.php?c=")
	urlBuilder.WriteString(strconv.FormatInt(int64(catID.EnumIndex()), 10))
	if page > 1 {
		urlBuilder.WriteString("&start=")
		urlBuilder.WriteString(strconv.FormatInt(int64(page-1)*20, 10))
		urlBuilder.WriteString("#pagestart")
	}
	url := urlBuilder.String()
	doc, err := web.New().GetDocument(ctx, url)
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
	decoder := charmap.Windows1251.NewDecoder()
	s.Find("td.pcatHead a").Each(func(i int, sl *goquery.Selection) {
		if title, ok := sl.Attr("title"); ok {
			titleString, _ := decoder.String(title)
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
		authorString, _ := decoder.String(sl.Text())
		topic.Author = authorString
	})
	s.Find("tbody > tr:nth-child(2) > td > span.genmed").Each(func(i int, sl *goquery.Selection) {
		timeString, _ := decoder.String(sl.Text())
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
