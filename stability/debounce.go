// Debounce limits the frequency of a function invocation
// so that only the first or last in a cluster of calls is actually performed.
// ---------------------------------------------------------------------------
// Limits how often a function can be called by restricting cluster of invocations.

package stability

import (
	"context"
	"sync"
	"time"
)

// Using Circuit(ctx context.Context) (string, error) function from Circuit Breaker pattern.

// Debounce-first compared to function-last.
// Tracks the last time function was called,
// returns a cached result if it called again before timeout (duration).
func DebounceFirst(circuit Circuit, dur time.Duration) Circuit {
	var treshold time.Time
	var result string
	var err error
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()

		// Set treshold after get result.
		defer func() {
			treshold = time.Now().Add(dur)
			m.Unlock()
		}()

		// Return cached result if calling before treshold time.
		if time.Now().Before(treshold) {
			return result, err
		}

		// Make a call if treshold passed.
		result, err = circuit(ctx)

		return result, err
	}
}

// Debounce-last compared to function-first.
// Set trashold time according to duration,
// delay of check if it passed (ticker).
// Make a call after treshold time passed.
func DebounceLast(circuit Circuit, dur time.Duration) Circuit {
	var treshold time.Time = time.Now()
	var result string
	var err error
	var ticker *time.Ticker
	var once sync.Once
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()

		treshold = time.Now().Add(dur) // Set trashold

		// Function executes only once
		once.Do(func() {
			// Set timer for 100 ms:
			// every 100 ms pass current time to channel in "ticker".
			ticker = time.NewTicker(time.Millisecond * 100)

			// Start gorutine.
			go func() {
				// Prepeare for next retry:
				// reset once to default, stop ticker.
				defer func() {
					m.Lock()
					once = sync.Once{}
					ticker.Stop()
					m.Unlock()
				}()

				for {
					select {
					// Recieve time after every "tick".
					// Make a call if time has passed the treshold.
					case <-ticker.C:
						m.Lock()
						if time.Now().After(treshold) {
							result, err = circuit(ctx)
							m.Unlock()
							return
						}
						m.Unlock()
					// Trace the context cancelation.
					case <-ctx.Done():
						m.Lock()
						result, err = "", ctx.Err()
						m.Unlock()
						return
					}
				}
			}()
		})

		return result, err
	}
}
