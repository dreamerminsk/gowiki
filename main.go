package main

import (
	"github.com/dreamerminsk/gowiki/tasks"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	tasks.UpdateCategories()

	tasks.UpdateForums()

	tasks.UpdateTopics()

}
