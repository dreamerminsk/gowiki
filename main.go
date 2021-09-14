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

	go tasks.UpdateCategories()

	go tasks.UpdateForums()

	go tasks.UpdateTopics()

	for {
		select {
		case <-keyChan:
			fmt.Println("CTRL-C: Завершаю работу. Всего записей: ")
			return
		case <-ticker.C:
			fmt.Printf("Всего  / Повторов  ( записей/сек) \n")

		}
	}

}
