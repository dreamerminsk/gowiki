package utils

import (
	"net/url"
	"strconv"

	"golang.org/x/net/html"
)

func GetParams(ref string) (url.Values, error) {
	u, err := url.Parse(html.UnescapeString(ref))
	if err != nil {
		return nil, err
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func GetParam(ref, name string) (value string, ok bool) {
	u, err := url.Parse(html.UnescapeString(ref))
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
