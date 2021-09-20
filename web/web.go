package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"

"github.com/dreamerminsk/gowiki/log"
)

const defaultUserAgent = "Mozilla/5.0 (Linux; Android 10; LM-X420) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Mobile Safari/537.36"
const timeFormat = "2006-01-02T15:04:05"

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
	requests *uint64 = new(uint64)
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
        reqID := atomic.AddUint64(requests, 1)
	log.log(fmt.Sprintf("%d - %s", reqID, url))
        ctx := context.WithValue(ctx, "reqID", reqID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.log(fmt.Sprintf("%d - %s", reqID, err))
		return nil, err
	}
	req.Header.Add("User-Agent", defaultUserAgent)
	return wc.Do(ctx, req)
}

func (wc *webClient) Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
        reqID := atomic.AddUint64(requests, 1)
	log.log(fmt.Sprintf("%d - %s", reqID, url))
        ctx := context.WithValue(ctx, "reqID", reqID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		log.log(fmt.Sprintf("%d - %s", reqID, err))
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("User-Agent", defaultUserAgent)
	return wc.Do(ctx, req)
}

func (wc *webClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
        reqID := ctx.Value("reqID").(uint64)
        log.log(fmt.Sprintf("%d - %s", reqID, url))
	err := wc.rateLimiter.Wait(ctx)
	if err != nil {
	        log.log(fmt.Sprintf("%d - %s", reqID, err))
		return nil, err
	}
	return wc.client.Do(req)
}
