package web

import (
	"io"
	"net/http"
	"time"
)

type WebClient struct {
	Client *http.Client
}

func NewWebClient() *WebClient {
	return &WebClient{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (wc *WebClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "")
	return wc.Client.Do(req)
}
func (wc *WebClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("User-Agent", "")
	return wc.Client.Do(req)
}
