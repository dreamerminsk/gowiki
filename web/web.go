package web

import (
	"bufio"
	"context"
	"fmt"
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

var (
	instance *webClient
	once     sync.Once
	r        *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func init() {
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
		log.Logf("StatusCode: %d", res.StatusCode)
		return nil, fmt.Errorf("StatusCode: %d", res.StatusCode)
	}

	doc, err := decode(res.Body, "")
	if err != nil {
		log.Logf("%s", err)
		return nil, err
	}

	html, _ := doc.Html()
	metrics.GetOrRegisterHistogram("Web.Response.Length", nil, metrics.NewExpDecaySample(1028, 0.015)).Update(int64(len(html)))
	metrics.GetOrRegisterString("Web.Response.URL", nil).Update(fmt.Sprintf("%s", res.Request.URL))
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		metrics.GetOrRegisterString("Web.Response.Title", nil).Update(s.Text())
	})

	return doc, nil
}

func decode(body io.Reader, charset string) (*goquery.Document, error) {
	r := bufio.NewReader(body)
	if charset == "" {
		charset = detectContentCharset(r)
	}

	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}

	if name, _ := htmlindex.Name(e); name != "utf-8" {
		r = bufio.NewReader(e.NewDecoder().Reader(r))
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Logf("%s", err)
		return nil, err
	}

	return doc, nil
}

func detectContentCharset(r *bufio.Reader) string {
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
	metrics.GetOrRegisterString("Web.Request.URL", nil).Update(fmt.Sprintf("%d - %s", reqID, req.URL))
	log.Logf("%d - %s", reqID, req.URL)

	defer func() {
		metrics.GetOrRegisterCounter("Web.Requests", nil).Inc(1)
	}()

	waitstart := time.Now()
	err := wc.rateLimiter.WaitN(ctx, r.Intn(32000)+8000)
	metrics.GetOrRegisterTimer("Web.ReqTime", nil).UpdateSince(waitstart)
	if err != nil {
		log.Logf("%d - %s", reqID, err)
		return nil, err
	}

	return wc.client.Do(req)
}
