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

func GetCategories(ctx context.Context) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)
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
			if strings.Contains(ref, "index.php?c=") {
				u, _ := url.Parse(ref)
				m, _ := url.ParseQuery(u.RawQuery)
				categoryID, _ := strconv.ParseInt(m["c"][0], 10, 32)
				categoryTitle, _ := decoder.String(s.Text())
				categories = append(categories, &model.Category{
					Model: gorm.Model{ID: uint(categoryID)},
					Title: categoryTitle,
				})
			}
		}
	})
	return categories, nil
}
