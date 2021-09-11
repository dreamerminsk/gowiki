package web

import (
	"io"
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (Linux; Android 10; LM-X420) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Mobile Safari/537.36"

type WebClient struct {
	Client *http.Client
}

func NewWebClient() *WebClient {
	return &WebClient{
		Client: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

func (wc *WebClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", defaultUserAgent)
	return wc.Client.Do(req)
}

func (wc *WebClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("User-Agent", defaultUserAgent)
	return wc.Client.Do(req)
}
