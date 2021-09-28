package nnmclub

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/web"
	"golang.org/x/text/encoding/charmap"
	"gorm.io/gorm"
)

func GetForums(ctx context.Context) ([]*model.Forum, error) {
	forums := make([]*model.Forum, 0)
	res, err := web.New().Get(ctx, "https://nnmclub.to/forum/index.php")
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	decoder := charmap.Windows1251.NewDecoder()
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if ref, ok := s.Attr("href"); ok {
			if strings.Contains(ref, "viewforum.php?f=") {
				u, _ := url.Parse(ref)
				m, _ := url.ParseQuery(u.RawQuery)
				forumID, _ := strconv.ParseInt(m["f"][0], 10, 32)
				forumTitle, _ := decoder.String(s.Text())
				forums = append(forums, &model.Forum{
					Model: gorm.Model{ID: uint(forumID)},
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
		Model: gorm.Model{ID: forumID},
		CatID: 0,
		Title: "",
	}
	res, err := web.New().Get(ctx, fmt.Sprintf("https://nnmclub.to/forum/viewforum.php?f=%d", forumID))
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc.Find("a.maintitle").Each(func(i int, s *goquery.Selection) {
		decoder := charmap.Windows1251.NewDecoder()
		forum.Title, _ = decoder.String(s.Text())
	})
	doc.Find("span.nav a[href]").Each(func(i int, s *goquery.Selection) {
		if ref, ok := s.Attr("href"); ok {
			if strings.Contains(ref, "index.php?c=") {
				u, _ := url.Parse(ref)
				m, _ := url.ParseQuery(u.RawQuery)
				CatID, _ := strconv.ParseInt(m["c"][0], 10, 32)
				forum.CatID = uint(CatID)
			}
		}
	})
	return forum, nil
}
