// Circuit Breaker automatically degrades service functions in response to a likely fault,
// preventing larger or cascading failures by eliminating recurring errors and providing
// reasonable error responses.
// ---------------------------------------------------------------------------------------
// Circuit Breaker pattern add functionality to control requests to service.
// Set number attemts of reaching service, delay, count failures and lock/unlock gorutine.

package stability

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

// Limits the number of failures requsting service by passing failureThreshold.
// Trying call a function while limit is off, then throw an error.
func Breaker(circuit Circuit, failureTreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock() // establish read lock

		d := consecutiveFailures - int(failureTreshold)

		// Set delay if attempts pass the limit, then retry to reach or throw an error.
		// Including authomatic reset mechanism with an exponential backoff in which the dura‐
		// tions of the delays between retries roughly doubles with each attempt.
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable (cached)")
			}
		}

		m.RUnlock() // release read lock

		response, err := circuit(ctx) // issue request proper

		m.Lock() // lock around shared resources
		defer m.Unlock()

		lastAttempt = time.Now() // record time of attempt

		if err != nil {
			consecutiveFailures++
			fmt.Println(consecutiveFailures)
			return response, err
		}

		consecutiveFailures = 0 // reset attempts to begin cycle again

		return response, nil
	}
}
