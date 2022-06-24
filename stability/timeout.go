// Timeout allows a process to stop waiting for an answer
// once it’s clear that an answer may not be coming.

package stability

import "context"

// Example of a function that takes too many time
// and need to be stopped.
type SlowFunc func(string) (string, error)

// The same function with context as an argument.
type WithContext func(context.Context, string) (string, error)

// Timeout wraps the slow function and return it with context
// so it could be used with context.WithTimeout()
func Timeout(slow SlowFunc) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		// Channels for recieving output of slow function.
		chResult := make(chan string)
		chError := make(chan error)

		// Call func, send output to channels.
		go func() {
			result, err := slow(arg)
			chResult <- result
			chError <- err
		}()

		// Return result of function executing from channels.
		// Trace context cancellation.
		select {
		case result := <-chResult:
			return result, <-chError
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
