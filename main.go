package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dreamerminsk/gowiki/tasks"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	ticker := time.NewTicker(time.Duration(60) * time.Second)
	defer ticker.Stop()

	keyChan := make(chan os.Signal, 1)
	signal.Notify(keyChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := [...]*tasks.Task{
		tasks.New(tasks.InitCategories),
		tasks.New(tasks.InitForums),
		tasks.New(tasks.UpdateForums),
		tasks.New(tasks.UpdateTopics),
	}
	for _, t := range queue {
		t.Run(ctx)
	}

	start := time.Now()

	for {
		select {
		case <-keyChan:
			fmt.Println("\nCTRL-C: Завершаю работу.")
			return
		case <-ticker.C:
			current := time.Now()
			fmt.Println("working ", (current.Sub(start)))
			for _, t := range queue {
				fmt.Println(t)
			}
		}
	}

}
