package nnmclub

import (
	"fmt"
	"time"
)

const (
	timeFormat = "02-01-2006"
)

func GetForumUrl() string {
	return "https://nnmclub.to/forum/index.php"
}

func GetViewUserUrl(userID uint) string {
	return fmt.Sprintf("https://nnmclub.to/forum/profile.php?mode=viewprofile&u=%d", userID)
}

func GetViewForumUrl(forumID uint) string {
	return fmt.Sprintf("https://nnmclub.to/forum/viewforum.php?f=%d", forumID)
}

func GetTopicsByCatUrl(catID, page int) string {
	if page > 1 {
		return fmt.Sprintf("https://nnmclub.to/forum/portal.php?c=%d&start=%d#pagestart", catID, (page-1)*20)
	}
	return fmt.Sprintf("https://nnmclub.to/forum/portal.php?c=%d#pagestart", catID)
}

func GetTopicsByDateUrl(date time.Time, page int) string {
	if page > 1 {
		return fmt.Sprintf("https://nnmclub.to/?d=%s&start=%d#pagestart", date.Format(timeFormat), (page-1)*20)
	}
	return fmt.Sprintf("https://nnmclub.to/?d=%s", date.Format(timeFormat))
}
