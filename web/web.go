package web

import (
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
	return wc.Client.Do(req)
}
