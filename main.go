package main

import (
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

	go tasks.InitOrUpdateCategories()

	go tasks.NewForums()

	go tasks.UpdateTopics()

	start := time.Now()

	for {
		select {
		case <-keyChan:
			fmt.Println("\nCTRL-C: Завершаю работу.")
			return
		case <-ticker.C:
			current := time.Now()
			fmt.Println("Всего  / Повторов  ( записей/сек) \n", (current.Sub(start)))

		}
	}

}
