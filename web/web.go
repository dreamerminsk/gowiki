package web

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/time/rate"

	"github.com/dreamerminsk/gowiki/log"
)

type key int

const (
	keyReqID key = iota
)



type webClient struct {
	client      *http.Client
	rateLimiter *rate.Limiter
}

type WebReader interface {
	Get(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error)
	doReq(ctx context.Context, req *http.Request) (*http.Response, error)
}

var (
	instance *webClient
	requests *uint64 = new(uint64)
	once     sync.Once
	r        *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func newReader() *webClient {
	return &webClient{
		client: &http.Client{
			Timeout: time.Second * 60,
		},
		rateLimiter: rate.NewLimiter(1000, 100000),
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
	log.Logf("%d - %s", reqID, url)
	ctx = context.WithValue(ctx, keyReqID, reqID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Logf("%d - %s", reqID, err)
		return nil, err
	}
	req.Header.Add("User-Agent", randomUserAgent())
	return wc.doReq(ctx, req)
}

func (wc *webClient) Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	reqID := atomic.AddUint64(requests, 1)
	log.Logf("%d - %s", reqID, url)
	ctx = context.WithValue(ctx, keyReqID, reqID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		log.Logf("%d - %s", reqID, err)
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("User-Agent", randomUserAgent())
	return wc.doReq(ctx, req)
}

func (wc *webClient) doReq(ctx context.Context, req *http.Request) (*http.Response, error) {
	reqID := ctx.Value(keyReqID).(uint64)
	log.Logf("%d - %s", reqID, req.URL)
	err := wc.rateLimiter.WaitN(ctx, r.Intn(64000)+32000)
	if err != nil {
		log.Logf("%d - %s", reqID, err)
		return nil, err
	}
	return wc.client.Do(req)
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
		log.Logf("%s", err)
		return nil, err
	}
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Logf("%s", err)
		return nil, err
	}
	decoder := charmap.Windows1251.NewDecoder()
	decoder.Reader(res.Body)
	return
}
