package nnmclub

import "fmt"

func GetForumUrl() string {
	return "https://nnmclub.to/forum/index.php"
}

func GetViewForumUrl(forumID uint) string {
	return fmt.Sprintf("https://nnmclub.to/forum/viewforum.php?f=%d", forumID)
}
