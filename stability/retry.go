// Retry accounts a possible transient fault by retrying a failed operation.

package stability

import (
	"context"
	"errors"
	"log"
	"time"
)

// Function that interacts with the service.
type Effector func(context.Context) (string, error)

// Takes effector and call it several time according to specified attempts with delay.
func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		// Call effector function.
		// If it fail try again "retries" times.
		for r := 0; ; r++ {
			result, err := effector(ctx)
			if err == nil || r >= retries {
				return result, err
			}

			// Output if attempt fail.
			log.Printf("%d attempt failed: retrying in %v", r+1, delay)

			// Continue after delay, if context is not cancelled.
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}

// (Example) emulator of transient error in role of Effector.
var count int

func EmulateTransientError(ctx context.Context) (string, error) {
	count++
	if count <= 3 {
		return "intentional fail", errors.New("error")
	} else {
		return "success", nil
	}
}
