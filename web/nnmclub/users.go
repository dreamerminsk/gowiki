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

func GetForumUsers(ctx context.Context, forumID uint) ([]*model.User, error) {
	users := make([]*model.User, 0)
	doc, err := web.New().GetDocument(ctx, GetViewForumUrl(forumID))
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if u, ok := ParseUser(ctx, s); ok {
			users = append(users, u)
		}
	})
	return users, nil
}

func ParseUser(ctx context.Context, s *goquery.Selection) (*model.User, bool) {
	if ref, ok := s.Attr("href"); ok {
		if strings.Contains(ref, "profile.php?mode=viewprofile&u=") {
			if userID, ok := utils.GetIntParam(ref, "u"); ok {
				return &model.User{
					ID:   uint(userID),
					Name: strings.TrimSpace(s.Text()),
				}, true
			}
		}
	}
	return nil, false
}
