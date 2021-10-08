package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/dreamerminsk/gowiki/model"
	"github.com/dreamerminsk/gowiki/storage"
	"github.com/dreamerminsk/gowiki/web/nnmclub"
)

func UpdateTopics(ctx context.Context) {
	g := storage.New()

	var cats = map[nnmclub.Category]int{
		nnmclub.Music:                  1,
		nnmclub.HDMusic:                1,
		nnmclub.MusicCollections:       1,
		nnmclub.AnimeAndManga:          1,
		nnmclub.BooksAndMediaMaterials: 1,
	}

	for {
		existValidPage := false
		for cat, page := range cats {
			if page > 0 {
				err := processTopicPage(ctx, g, cat, page)
				if err != nil {
					cats[cat] = -1
					continue
				}
				cats[cat] = cats[cat] + 1
				existValidPage = true
			}
		}
		if !existValidPage {
			break
		}
	}
}

func processTopicPage(ctx context.Context, g storage.Storage, catID nnmclub.Category, page int) error {
	topics, err := nnmclub.GetTopics(ctx, catID, page)
	if err != nil {
		return err
	}
	for _, topic := range topics {
		fmt.Println("ID: ", topic.ID)
		fmt.Println("Title: ", topic.Title)
		fmt.Println("Author: ", topic.Author)
		fmt.Println("Published: ", topic.Published.Format(time.RFC3339))
		fmt.Println("Likes: ", topic.Likes)
		fmt.Println("Magnet: ", topic.Magnet)
		err = insertOrUpdate(g, topic)
		if err != nil {
			return err
		}
		fmt.Println("-------------------------")
	}
	return nil
}

func insertOrUpdate(g storage.Storage, topic *model.Topic) error {
	oldTopic, err := g.GetTopicByID(topic.ID)
	if err != nil {
		fmt.Println("SELECT ERROR: ", reflect.TypeOf(err), err)
	}
	if oldTopic.ID == 0 {
		err = g.Create(topic).Error
		if err != nil {
			fmt.Println("INSERT ERROR: ", reflect.TypeOf(err), err)
			return err
		}
	} else if topic.Likes > oldTopic.Likes {
		fmt.Printf("\tDIFF LIKES: %d\r\n", topic.Likes-oldTopic.Likes)
		err = g.UpdateTopic(topic)
		if err != nil {
			fmt.Println("UPDATE ERROR: ", reflect.TypeOf(err), err)
			return err
		}
	}
	return nil
}

func RandDuration(min, max int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Second * time.Duration((r.Intn(max-min) + min))
}
