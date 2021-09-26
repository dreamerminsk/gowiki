package utils

import (
	"net/url"
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
