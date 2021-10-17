package web

import (
	"bufio"
	"context"
	"expvar"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/time/rate"

	"github.com/dreamerminsk/gowiki/log"
	"github.com/dreamerminsk/gowiki/metrics"
	"github.com/dreamerminsk/gowiki/metrics/exp"
)

type key int

const (
	keyReqID key = iota
)

type webClient struct {
	client        *http.Client
	lastRequestId *uint64
	rateLimiter   *rate.Limiter
}

type WebReader interface {
	GetDocument(ctx context.Context, url string) (*goquery.Document, error)
	Get(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error)
	doReq(ctx context.Context, req *http.Request) (*http.Response, error)
}

type WebStats struct {
	Requests    uint64
	WaitTime    float64
	WaitTimeMin float64
	WaitTimeAvg float64
	WaitTimeMax float64
}

var (
	instance *webClient
	once     sync.Once
	stats    *WebStats   = &WebStats{WaitTimeMin: 1000}
	statsM   *sync.Mutex = &sync.Mutex{}
	r        *rand.Rand  = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func webstats() interface{} {
	statsM.Lock()
	defer statsM.Unlock()
	stats.WaitTimeAvg = stats.WaitTime / float64(stats.Requests)
	return *stats
}

func init() {
	expvar.Publish("WebStats", expvar.Func(webstats))
	exp.Exp(metrics.DefaultRegistry)
}

func newReader() *webClient {
	return &webClient{
		client: &http.Client{
			Timeout: time.Second * 60,
		},
		lastRequestId: new(uint64),
		rateLimiter:   rate.NewLimiter(1000, 100000),
	}
}

func New() WebReader {
	once.Do(func() {
		instance = newReader()
	})

	return instance
}

func (wc *webClient) GetDocument(ctx context.Context, url string) (*goquery.Document, error) {
	res, err := wc.Get(ctx, url)
	if err != nil {
		log.Logf("%s", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Logf("%s", err)
		return nil, err
	}
	doc, err := decode(res.Body, "")
	if err != nil {
		log.Logf("%s", err)
		return nil, err
	}
	return doc, nil
}

func decode(body io.Reader, charset string) (*goquery.Document, error) {
	if charset == "" {
		charset = detectContentCharset(body)
	}

	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}

	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Logf("%s", err)
		return nil, err
	}

	return doc, nil
}

func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, _ := charset.DetermineEncoding(data, ""); name != "" {
			return name
		}
	}
	return "utf-8"
}

func (wc *webClient) Get(ctx context.Context, url string) (*http.Response, error) {
	reqID := atomic.AddUint64(wc.lastRequestId, 1)
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
	reqID := atomic.AddUint64(wc.lastRequestId, 1)
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

	defer func() {
		statsM.Lock()
		defer statsM.Unlock()
		stats.Requests++
		metrics.GetOrRegisterCounter("WebStats.Requests", nil).Inc(1)
	}()

	waitstart := time.Now()
	err := wc.rateLimiter.WaitN(ctx, r.Intn(32000)+8000)
	statsM.Lock()
	wait := float64(time.Since(waitstart).Seconds())
	stats.WaitTime += wait
	if wait > 1.0 && wait < stats.WaitTimeMin {
		stats.WaitTimeMin = wait
	}
	if wait > stats.WaitTimeMax {
		stats.WaitTimeMax = wait
	}
	statsM.Unlock()
	if err != nil {
		log.Logf("%d - %s", reqID, err)
		return nil, err
	}

	return wc.client.Do(req)
}
