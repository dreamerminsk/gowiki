package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/dreamerminsk/gowiki/nnmclub/client"
	"github.com/dreamerminsk/gowiki/nnmclub/model"
	"github.com/dreamerminsk/gowiki/nnmclub/storage"
	"github.com/dreamerminsk/gowiki/tasks"
)

func UpdateTopics(ctx context.Context, t *tasks.Task) {
	g := storage.New()

	var cats = map[client.Category]int{
		client.Music:                  1,
		client.HDMusic:                1,
		client.MusicCollections:       1,
		client.AnimeAndManga:          1,
		client.BooksAndMediaMaterials: 1,
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

func processTopicPage(ctx context.Context, g storage.Storage, catID client.Category, page int) error {
	topics, err := client.GetTopics(ctx, catID, page)
	if err != nil {
		return err
	}
	for _, topic := range topics {
		fmt.Println("ID: ", topic.ID)
		fmt.Println("Title: ", topic.Title)
		fmt.Println("Author: ", topic.Author)
		fmt.Println("Published: ", topic.Published.Format(time.RFC3339))
		fmt.Println("Likes: ", topic.Likes)
		fmt.Println("Comments: ", topic.Comments)
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
		fmt.Printf("SELECT ERROR: %s, %s\r\n", reflect.TypeOf(err), err)
		err = g.Create(topic).Error
		if err != nil {
			fmt.Printf("INSERT ERROR: %s, %s\r\n", reflect.TypeOf(err), err)
			return err
		}
	} else if topic.Likes != oldTopic.Likes {
		fmt.Printf("\tDIFF Likes: %d\r\n", topic.Likes-oldTopic.Likes)
		err = g.UpdateTopic(topic)
		if err != nil {
			fmt.Printf("UPDATE ERROR: %s, %s\r\n", reflect.TypeOf(err), err)
			return err
		}
	} else if topic.Comments != oldTopic.Comments {
		fmt.Printf("\tDIFF Comments: %d\r\n", topic.Comments-oldTopic.Comments)
		err = g.UpdateTopic(topic)
		if err != nil {
			fmt.Printf("UPDATE ERROR: %s, %s\r\n", reflect.TypeOf(err), err)
			return err
		}
	}
	return nil
}

func RandDuration(min, max int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Second * time.Duration((r.Intn(max-min) + min))
}
