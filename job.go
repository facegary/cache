package cache

import (
	"time"
)

type ijob interface {
	Stop()
}

type job struct {
	cdone  chan bool
	ticker *time.Ticker
}

func newJob(f func(t time.Time), d time.Duration) ijob {
	ticker := time.NewTicker(d)
	cdone := make(chan bool)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-cdone:
				return
			case t := <-ticker.C:
				f(t)
			}
		}
	}(ticker)
	return &job{ticker: ticker, cdone: cdone}
}

func (j *job) Stop() {
	j.ticker.Stop()
	j.cdone <- true
}
