package health

import (
	"context"
	"encoding/json"
	"time"

	"go-monorepo/logging"
)

// StartVarsPoller starts vars poller and triggers callback function.
func StartVarsPoller(ctx context.Context, fn func([]byte)) {
	go func() {
		for ctx.Err() == nil {
			time.Sleep(pollPeriod)
			_, _, vars := gatherInfos()
			b, err := json.Marshal(vars)
			if err != nil {
				logging.Get().Error(err)
				continue
			}

			fn(b)
		}
	}()
}
