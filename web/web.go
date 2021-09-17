package web

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const defaultUserAgent = "Mozilla/5.0 (Linux; Android 10; LM-X420) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Mobile Safari/537.36"

type webClient struct {
	client      *http.Client
	rateLimiter *rate.Limiter
}

type WebReader interface {
	Get(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error)
        Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

var (
	instance *webClient
	once     sync.Once
)

func newReader() *webClient {
	return &webClient{
		client: &http.Client{
			Timeout: time.Second * 60,
		},
		rateLimiter: rate.NewLimiter(rate.Every(60*time.Second), 1),
	}
}

func New() WebReader {
	once.Do(func() {
		instance = newReader()
	})

	return instance
}

func (wc *webClient) Get(ctx context.Context, url string) (*http.Response, error) {
	err := wc.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", defaultUserAgent)
	return wc.client.Do(req)
}

func (wc *webClient) Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	err := wc.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("User-Agent", defaultUserAgent)
	return wc.client.Do(req)
}

func (wc *webClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	err := wc.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return wc.client.Do(req)
}
