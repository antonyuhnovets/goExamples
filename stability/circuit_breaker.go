package stability

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Specifies the signature of the function that’s interacting with upstream service.
type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureTreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock() // Establish read lock

		d := consecutiveFailures - int(failureTreshold)

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

		consecutiveFailures = 0

		return response, nil
	}
}
