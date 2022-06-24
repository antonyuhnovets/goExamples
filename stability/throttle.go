// Throttle limits the frequency of a function call
// to some maximum number of invocations per unit of time.
// -------------------------------------------------------
// Throttle`s difference from Debounce is that Debonce manage
// cluster of activity and allows one call per cluster, while
// Throttle set limit of calls per unit of time.

package stability

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Function to regulate Effector(context.Context) (string, error) from Retry pattern.

// Function uses token bucket algorithm, which uses the analogy
// of a bucket that can hold some maximum number of tokens.
// When a function is called, a token is taken from the bucket,
// which then refills at some fixed rate.
func Throttle(effector Effector, max uint, refill uint, delay time.Duration) Effector {
	var tokens = max // calls (or requests) allowed per time unit (delay)
	var once sync.Once

	return func(ctx context.Context) (string, error) {
		// Trace context errors.
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		// Refill tokens after delay.
		once.Do(func() {
			ticker := time.NewTicker(delay) // send time to channel after delay

			go func() {
				defer ticker.Stop()

				// If time delay passed refilling tokens quality.
				// Exit if context closed.
				for {
					select {
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					case <-ctx.Done():
						return
					}
				}
			}()
		})
		// Error if calls number is more than allowed.
		if tokens <= 0 {
			return "", fmt.Errorf("too many calls")
		}
		// Call the function and decrement tokens amount.
		tokens--

		return effector(ctx)
	}
}
