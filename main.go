package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	res, err := http.Get("http://journaldev.com")
	if err != nil {
        log.Fatal(err)
	}
    defer res.Body.Close()
}