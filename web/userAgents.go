package web

import (
	"math/rand"
	"time"
)

var uaRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

var userAgents = [...]string{
	"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.50 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/96.0.4664.36 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/605.1.15 (KHTML, like Gecko)",
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
}

func randomUserAgent() string {
	return userAgents[uaRand.Intn(len(userAgents))]
}
