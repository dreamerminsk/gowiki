package nnmclub

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/utils"
	"github.com/dreamerminsk/gowiki/web"
)

func GetForums(ctx context.Context) ([]*model.Forum, error) {
	forums := make([]*model.Forum, 0)
	doc, err := web.New().GetDocument(ctx, GetForumUrl())
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if f, ok := ParseForum(ctx, s); ok {
			if len(f.Title) > 0 {
				forums = append(forums, f)
			}
		}
	})
	return forums, nil
}

func ParseForum(ctx context.Context, s *goquery.Selection) (*model.Forum, bool) {
	if ref, ok := s.Attr("href"); ok {
		if strings.Contains(ref, "viewforum.php?f=") {
			if fID, ok := utils.GetIntParam(ref, "f"); ok {
				return &model.Forum{
					ID:    uint(fID),
					Title: strings.TrimSpace(s.Text()),
				}, true
			}
		}
	}
	return nil, false
}

func GetForum(ctx context.Context, forumID uint) (*model.Forum, error) {
	forum := &model.Forum{
		ID:    forumID,
		CatID: 0,
		Title: "",
	}
	doc, err := web.New().GetDocument(ctx, GetViewForumUrl(forumID))
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc.Find("a.maintitle").Each(func(i int, s *goquery.Selection) {
		forum.Title = strings.TrimSpace(s.Text())
	})
	if len(forum.Title) == 0 {
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			if ref, ok := s.Attr("href"); ok {
				if strings.Contains(ref, fmt.Sprintf("viewforum.php?f=%d", forumID)) {
					forum.Title = strings.TrimSpace(s.Text())
				}
			}
		})
	}
	doc.Find("span.nav a[href]").Each(func(i int, s *goquery.Selection) {
		if cat, ok := ParseCategory(ctx, s); ok {
			forum.CatID = cat.ID
		}
	})
	return forum, nil
}
