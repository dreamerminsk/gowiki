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
