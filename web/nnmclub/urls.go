package nnmclub

import (
	"fmt"
)

func GetForumUrl() string {
	return "https://nnmclub.to/forum/index.php"
}

func GetViewForumUrl(forumID uint) string {
	return fmt.Sprintf("https://nnmclub.to/forum/viewforum.php?f=%d", forumID)
}

func GetCatTopicsUrl(catID, page int) string {
	if page > 1 {
		return fmt.Sprintf("https://nnmclub.to/forum/portal.php?c=%d&start=%d#pagestart", catID, (page-1)*20)
	}
	return fmt.Sprintf("https://nnmclub.to/forum/portal.php?c=%d#pagestart", catID)
}
