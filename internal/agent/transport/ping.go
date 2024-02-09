package transport

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func Ping(ctx context.Context, server string, interval time.Duration) error {
	count := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			resp, err := http.Get(server)
			if err != nil {
				count += 1
				if count >= 5 {
					resp.Body.Close()
					return errors.New("max retries reached")
				}
			} else {
				resp.Body.Close()
				count = 0
			}
			time.Sleep(interval)

		}
	}
}
