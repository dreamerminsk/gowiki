package web

import (
	"io"
"sync"
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (Linux; Android 10; LM-X420) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Mobile Safari/537.36"

type webClient struct {
	client *http.Client
        mu sync.Mutex
}

type WebReader interface{
Get(url string) (resp *http.Response, err error)
Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

func New() *WebReader {
	return &webClient{
		client: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

func (wc *webClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", defaultUserAgent)
	return wc.client.Do(req)
}

func (wc *webClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("User-Agent", defaultUserAgent)
	return wc.client.Do(req)
}
