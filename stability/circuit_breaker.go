// Circuit Breaker automatically degrades service functions in response to a likely fault,
// preventing larger or cascading failures by eliminating recurring errors and providing
// reasonable error responses.
// ---------------------------------------------------------------------------------------
// Circuit breaker pattern add functionality to control requests to service.
// Set number attemts of reaching service, delay, count failures and lock/unlock gorutine.

package stability

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Specifies the signature of the function that’s interacting with upstream service.
type Circuit func(context.Context) (string, error)

// Limits the number of failures requsting service by passing failureThreshold.
// Trying call a function while limit is off, then throw an error.
func Breaker(circuit Circuit, failureTreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock() // Establish read lock

		d := consecutiveFailures - int(failureTreshold)

		// Set delay if attempts pass the limit, then retry to reach or throw an error.
		// Including authomatic reset mechanism with an exponential backoff in which the dura‐
		// tions of the delays between retries roughly doubles with each attempt.
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << 2)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock() // Release read lock

		response, err := circuit(ctx) // Issue request proper

		m.Lock() // Lock around shared resources
		defer m.Unlock()

		lastAttempt = time.Now() // Record time of attempt

		if err != nil {
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0 // Reset attempts to begin cycle again

		return response, nil
	}
}
