package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/dreamerminsk/gowiki/metrics"
	"github.com/dreamerminsk/gowiki/nnmclub/client"
	"github.com/dreamerminsk/gowiki/nnmclub/model"
	"github.com/dreamerminsk/gowiki/nnmclub/storage"
	"github.com/dreamerminsk/gowiki/tasks"
)

func UpdateTopics(ctx context.Context, t *tasks.Task) {
	g := storage.New()

	var cats = map[client.Category]int{
		client.NewMovies:              1,
		client.ForeignMovies:          1,
		client.ForeignTVSeries:        1,
		client.DomesticTVSeries:       1,
		client.DomesticMovies:         1,
		client.Music:                  1,
		client.HDMusic:                1,
		client.MusicCollections:       1,
		client.AnimeAndManga:          1,
		client.BooksAndMediaMaterials: 1,
		client.HDUHDAnd3DMovies:       1,
		client.DocAndTVShows:          1,
		client.SportsAndHumor:         1,
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
	metrics.GetOrRegisterValues("UpdateTopics", nil).Add("Topic", fmt.Sprintf("%s-%d", catID, page))
	topics, err := client.GetTopics(ctx, catID, page)
	if err != nil {
		return err
	}
	pub := time.Now().AddDate(-100, 0, 0)
	for _, topic := range topics {
		fmt.Println("ID: ", topic.ID)
		fmt.Println("Title: ", topic.Title)
		fmt.Println("Author: ", topic.Author)
		fmt.Println("Published: ", topic.Published.Format(time.RFC3339))
		fmt.Println("Likes: ", topic.Likes)
		fmt.Println("Comments: ", topic.Comments)
		fmt.Println("Magnet: ", topic.Magnet)
		fmt.Println("Size: ", topic.Size)
		err = insertOrUpdate(g, topic)
		if err != nil {
			return err
		}
		fmt.Println("-------------------------")
		if topic.Published.After(pub) {
			pub = topic.Published
		}
	}
	metrics.GetOrRegisterValues("UpdateTopics.Published", nil).Add("Last", pub.Format(time.RFC3339))
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
		} else {
			metrics.GetOrRegisterCounter("UpdateTopics.NewLikes", nil).Inc(topic.Likes)
			metrics.GetOrRegisterCounter("UpdateTopics.NewComments", nil).Inc(topic.Comments)
		}
	} else {
		metrics.GetOrRegisterCounter("UpdateTopics.Likes", nil).Inc(oldTopic.Likes)
		metrics.GetOrRegisterCounter("UpdateTopics.Comments", nil).Inc(oldTopic.Comments)
		if topic.Likes != oldTopic.Likes {
			metrics.GetOrRegisterCounter("UpdateTopics.NewLikes", nil).Inc(topic.Likes - oldTopic.Likes)
			fmt.Printf("\tDIFF Likes: %d\r\n", topic.Likes-oldTopic.Likes)
			topic.CreatedAt = oldTopic.CreatedAt
			err = g.UpdateTopic(topic)
			if err != nil {
				fmt.Printf("UPDATE ERROR: %s, %s\r\n", reflect.TypeOf(err), err)
				return err
			}
		}
		if topic.Comments != oldTopic.Comments {
			metrics.GetOrRegisterCounter("UpdateTopics.NewComments", nil).Inc(topic.Comments - oldTopic.Comments)
			fmt.Printf("\tDIFF Comments: %d\r\n", topic.Comments-oldTopic.Comments)
			topic.CreatedAt = oldTopic.CreatedAt
			err = g.UpdateTopic(topic)
			if err != nil {
				fmt.Printf("UPDATE ERROR: %s, %s\r\n", reflect.TypeOf(err), err)
				return err
			}
		}
		if topic.Size != oldTopic.Size {
			fmt.Printf("\tDIFF Size: %s\r\n", topic.Size)
			topic.CreatedAt = oldTopic.CreatedAt
			err = g.UpdateTopic(topic)
			if err != nil {
				fmt.Printf("UPDATE ERROR: %s, %s\r\n", reflect.TypeOf(err), err)
				return err
			}
		}
	}
	return nil
}

func RandDuration(min, max int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Second * time.Duration((r.Intn(max-min) + min))
}
