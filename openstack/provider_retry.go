package openstack

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud/v2"
)

// requestTimeoutRetryDelay is the base delay between retries. It grows
// linearly with each attempt.
const requestTimeoutRetryDelay = 1 * time.Second

// requestTimeoutRetryFunc returns a gophercloud.RetryFunc that retries
// idempotent (GET and HEAD) requests answered with "408 Request Timeout".
//
// A load balancer such as HAProxy may answer a request with 408 when it is
// sent on a keep-alive connection at the moment the load balancer closes that
// idle connection. Without this, a single 408 while polling for a resource
// status (e.g. waiting for an instance to rebuild) fails the whole operation.
// Only GET and HEAD are retried, as they are idempotent and safe to repeat
// even if the request did reach the API.
//
// Any other error is returned unchanged, so callers still see the original
// error type and response code.
func requestTimeoutRetryFunc(maxRetries uint, baseDelay time.Duration) gophercloud.RetryFunc {
	return func(ctx context.Context, method, url string, _ *gophercloud.RequestOpts, err error, failCount uint) error {
		if method != http.MethodGet && method != http.MethodHead {
			return err
		}

		if !gophercloud.ResponseCodeIs(err, http.StatusRequestTimeout) {
			return err
		}

		// failCount is the request's total retry count, which gophercloud
		// shares with other retries such as 429 backoff, so the limit
		// applies to all retries of this request combined.
		if failCount > maxRetries {
			log.Printf("[DEBUG] Got 408 Request Timeout on %s %s, retry limit of %d reached: %s", method, url, maxRetries, err)

			return err
		}

		if ctx.Err() != nil {
			return err
		}

		delay := baseDelay * time.Duration(failCount)
		log.Printf("[DEBUG] Got 408 Request Timeout on %s %s, retrying in %s (retry %d, limit %d)", method, url, delay, failCount, maxRetries)

		t := time.NewTimer(delay)
		defer t.Stop()

		select {
		case <-ctx.Done():
			return err
		case <-t.C:
			return nil
		}
	}
}
