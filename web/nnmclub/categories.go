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
	"gorm.io/gorm"
)

type Category int

const (
	AnimeAndManga Category = iota + 1
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

func (c Category) String() string {
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

func (c Category) EnumIndex() int {
	return int(c)
}

func GetCategories(ctx context.Context) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)
	doc, err := web.New().GetDocument(ctx, "https://nnmclub.to/forum/index.php")
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if ref, ok := s.Attr("href"); ok {
			if strings.Contains(ref, "index.php?c=") {
				u, _ := url.Parse(ref)
				m, _ := url.ParseQuery(u.RawQuery)
				categoryID, _ := strconv.ParseInt(m["c"][0], 10, 32)
				categoryTitle := s.Text()
				categories = append(categories, &model.Category{
					Model: gorm.Model{ID: uint(categoryID)},
					Title: categoryTitle,
				})
			}
		}
	})
	return categories, nil
}
