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
	"github.com/dreamerminsk/gowiki/web"
)

func GetForums(ctx context.Context) ([]*model.Forum, error) {
	forums := make([]*model.Forum, 0)
	doc, err := web.New().GetDocument(ctx, "https://nnmclub.to/forum/index.php")
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if ref, ok := s.Attr("href"); ok {
			if strings.Contains(ref, "viewforum.php?f=") {
				u, _ := url.Parse(ref)
				m, _ := url.ParseQuery(u.RawQuery)
				forumID, _ := strconv.ParseInt(m["f"][0], 10, 32)
				forumTitle := s.Text()
				forums = append(forums, &model.Forum{
					ID:    uint(forumID),
					CatID: 0,
					Title: forumTitle,
				})
			}
		}
	})
	return forums, nil
}

func GetForum(ctx context.Context, forumID uint) (*model.Forum, error) {
	forum := &model.Forum{
		ID:    forumID,
		CatID: 0,
		Title: "",
	}
	doc, err := web.New().GetDocument(ctx, fmt.Sprintf("https://nnmclub.to/forum/viewforum.php?f=%d", forumID))
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc.Find("a.maintitle").Each(func(i int, s *goquery.Selection) {
		forum.Title = strings.TrimSpace(s.Text())
	})
	doc.Find("span.nav a[href]").Each(func(i int, s *goquery.Selection) {
		if cat, ok := GetCategory(ctx, s); ok {
			forum.CatID = cat.ID
		}
	})
	return forum, nil
}
