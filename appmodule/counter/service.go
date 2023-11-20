package counter

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"go-monorepo/cache"
	"go-monorepo/health"
)

// Service makes meaningless counter and let clients to query the number.
type Service struct {
	count int64
}

// Start starts service.
func (s *Service) Start(ctx context.Context) {
	redis := cache.Redis()

	go func() {
		defer func() {
			redis.Set("count", s.count, 0)
		}()

		// create tracker to track this goroutine's healthy.
		info := health.NewInfo(
			"naive-count-worker",
			10*time.Second,
			health.ProbeReady,
		)

		// report this goroutine is exited and handles panic.
		defer info.Down()

		ticker := time.NewTicker(time.Second)

		count, err := redis.Get("count").Int64()
		if err != nil {
			if !errors.Is(cache.Nil, err) {
				log.Panicln("failed to restore counter")
			}
		}
		s.count = count

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				atomic.AddInt64(&s.count, 1)
			}

			// mark this goroutine is still alive.
			info.Up()
		}
	}()
}

// Query returns the number.
func (s *Service) Query() int64 {
	return atomic.LoadInt64(&s.count)
}
