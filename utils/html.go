package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/dreamerminsk/gowiki/log"
	"golang.org/x/text/encoding/charmap"
)

func GetParam(ref, name string) (value string, ok bool) {
	u, err := url.Parse(ref)
	if err != nil {
		return "", false
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", false
	}
	return q.Get(name), true
}

func GetIntParam(ref, name string) (value int, ok bool) {
	if p, ok := GetParam(ref, name); ok {
		n, err := strconv.Atoi(p)
		if err != nil {
			return 0, false
		}
		return n, true
	}
	return 0, false
}

func NewDocumentFromReader(res *http.Response) (doc *goquery.Document, err error) {
	if res == nil {
		return nil, errors.New("response is nil")
	}
	defer res.Body.Close()
	if res.Request == nil {
		return nil, errors.New("response.Request is nil")
	}
	if res.StatusCode != http.StatusOK {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Log(fmt.Sprintf("%s", err))
		return nil, err
	}
	decoder := charmap.Windows1251.NewDecoder()
	decoder.Reader(res.Body)
	return
}
