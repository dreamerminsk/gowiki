package main

import (
	"context"
	_ "expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	nnmclub "github.com/dreamerminsk/gowiki/nnmclub/tasks"
	"github.com/dreamerminsk/gowiki/tasks"
)

func main() {

	go http.ListenAndServe(":8080", nil)

	ticker := time.NewTicker(time.Duration(60) * time.Second)
	defer ticker.Stop()

	keyChan := make(chan os.Signal, 1)
	signal.Notify(keyChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := [...]*tasks.Task{
		tasks.New(nnmclub.InitCategories),
		tasks.New(nnmclub.InitForums),
		//tasks.New(nnmclub.InitUsers),
		tasks.New(nnmclub.UpdateForums),
		tasks.New(nnmclub.UpdateTopics),
	}
	for _, t := range queue {
		go t.Run(ctx)
	}

	for {
		select {
		case <-keyChan:
			for _, t := range queue {
				fmt.Println(t)
			}
			return
		case <-ticker.C:
			for _, t := range queue {
				fmt.Println(t)
			}
		}
	}

}
