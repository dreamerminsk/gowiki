package web

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gorm.io/gorm"

	"github.com/PuerkitoBio/goquery"
	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/utils"
        "github.com/dreamerminsk/gowiki/log"
)

var client WebReader = New()

type NnmClubCategory int

const (
	AnimeAndManga NnmClubCategory = iota + 1
	HDMusic
	ForeignTVSeries
	DomesticTVSeries
	BooksAndMediaMaterials
	ForeignMovies
	ForChildrenAndParents
	Various
	PCGames
	NewMovies
	HDUHDAnd3DMovies
	Music
	DomesticMovies
	AndroidMobile
	EverythingForApple
	MultimediaDesignGraphics
	_
	EverythingForNIXSystems
	Applications
	GamesForConsoles
	TheaterMusicVideoMiscellaneous
	DocTVBrands
	DocAndTVShows
	SportsAndHumor
	MusicCollections
)

func (c NnmClubCategory) String() string {
	return [...]string{
		"Аниме и Манга",
		"Музыка HD",
		"Зарубежные сериалы",
		"Наши сериалы",
		"Книги и медиаматериалы",
		"Зарубежное кино",
		"Детям и родителям",
		"Разное",
		"Игры для ПК",
		"Новинки кино",
		"HD, UHD и 3D Кино",
		"Музыка",
		"Наше кино",
		"Android, Mobile",
		"Всё для Apple",
		"Мультимедиа, Дизайн, Графика",
		"_",
		"Все для *NIX систем",
		"Программы",
		"Игры для консолей",
		"Театр, МузВидео, Разное",
		"Док. TV-бренды",
		"Док. и телепередачи",
		"Спорт и Юмор",
		"Музыка (сборники)"}[c-1]
}

func (c NnmClubCategory) EnumIndex() int {
	return int(c)
}

func GetCategories(ctx context.Context) (map[uint]*model.Category, error) {
	categories := make(map[uint]*model.Category)
	res, err := client.Get(ctx, "https://nnmclub.to/forum/index.php")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	decoder := charmap.Windows1251.NewDecoder()
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if ref, ok := s.Attr("href"); ok {
			if strings.Contains(ref, "index.php?c=") {
				u, _ := url.Parse(ref)
				m, _ := url.ParseQuery(u.RawQuery)
				categoryID, _ := strconv.ParseInt(m["c"][0], 10, 32)
				categoryTitle, _ := decoder.String(s.Text())
				categories[uint(categoryID)] = &model.Category{
					Model: gorm.Model{ID: uint(categoryID)},
					Title: categoryTitle,
				}
			}
		}
	})
	return categories, nil
}

func GetForums(ctx context.Context) (map[uint]*model.Forum, error) {
	forums := make(map[uint]*model.Forum)
	res, err := client.Get(ctx, "https://nnmclub.to/forum/index.php")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	decoder := charmap.Windows1251.NewDecoder()
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if ref, ok := s.Attr("href"); ok {
			if strings.Contains(ref, "viewforum.php?f=") {
				u, _ := url.Parse(ref)
				m, _ := url.ParseQuery(u.RawQuery)
				forumID, _ := strconv.ParseInt(m["f"][0], 10, 32)
				forumTitle, _ := decoder.String(s.Text())
				forums[uint(forumID)] = &model.Forum{
					Model: gorm.Model{ID: uint(forumID)},
					CatID: 0,
					Title: forumTitle,
				}
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
	res, err := client.Get(ctx, fmt.Sprintf("https://nnmclub.to/forum/viewforum.php?f=%d", forumID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
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

func GetTopics(ctx context.Context, catID NnmClubCategory, page int) map[uint]*model.Topic {
	topics := make(map[uint]*model.Topic)
	var urlBuilder strings.Builder
	urlBuilder.WriteString("https://nnmclub.to/forum/portal.php?c=")
	urlBuilder.WriteString(strconv.FormatInt(int64(catID.EnumIndex()), 10))
	if page > 1 {
		urlBuilder.WriteString("&start=")
		urlBuilder.WriteString(strconv.FormatInt(int64(page-1)*20, 10))
		urlBuilder.WriteString("#pagestart")
	}
	url := urlBuilder.String()
	fmt.Println()
	fmt.Println(url)
	res, err := client.Get(ctx, url)
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
		return isTopic(s)
	}).Each(func(i int, s *goquery.Selection) {
		topic := getTopic(s)
		topics[topic.ID] = topic
	})
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
