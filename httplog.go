package mailgun

import (
	"log"
	"net/http"
	"os"
	"time"
)

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	return nil
}

type loggedRoundTripper struct {
	rt  http.RoundTripper
	log HTTPLogger
}

func (c *loggedRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	c.log.LogRequest(request)
	startTime := time.Now()
	response, err := c.rt.RoundTrip(request)
	duration := time.Since(startTime)
	c.log.LogResponse(request, response, err, duration)
	return response, err
}

// NewLoggedTransport ...
func NewLoggedTransport(rt http.RoundTripper, log HTTPLogger) http.RoundTripper {
	return &loggedRoundTripper{rt: rt, log: log}
}

// HTTPLogger ...
type HTTPLogger interface {
	LogRequest(*http.Request)
	LogResponse(*http.Request, *http.Response, error, time.Duration)
}

type httpLog struct {
	log *log.Logger
}

func newLogger() *httpLog {
	return &httpLog{
		log: log.New(os.Stderr, "✪", log.LstdFlags),
	}
}

func (l *httpLog) LogRequest(req *http.Request) {
	l.log.Printf("▶ %s %s", req.Method, req.URL.String())
}

func (l *httpLog) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		l.log.Println(err)
	} else {
		l.log.Printf("◀ method=%s status=%d durationMs=%d %s, %s", req.Method, res.StatusCode, duration, req.URL.String(), req.Referer())
	}
}
