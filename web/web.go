package web

import "net/http"

type WebClient struct {
	Client *http.Client
}
