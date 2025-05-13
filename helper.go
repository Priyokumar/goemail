package goemail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand"
	"time"
)

// processTemplate parses the template file at the specified path and executes it
// with the provided data, returning the resulting string or an error if any occurs.
//
// Parameters:
//   - templatePath: The file path to the template file to be parsed.
//   - data: The data to be injected into the template during execution.
//
// Returns:
//   - A string containing the rendered template.
//   - An error if the template parsing or execution fails.
func processTemplate(templatePath string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("template parsing error")
		return "", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		log.Println("template parsing error")
		return "", err
	}
	log.Println(buf.String())
	return buf.String(), nil
}

// exponentialBackOffRetry attempts to execute a given function with exponential backoff retries.
// It retries the function up to the specified number of times or until the context is canceled.
//
// Parameters:
//   - ctx: The context to manage cancellation and timeout.
//   - e: A pointer to an email object (custom type) passed to the function being retried.
//   - fn: The function to be executed. It takes an *email as input and returns an error.
//   - retry: The maximum number of retry attempts.
//
// Behavior:
//   - The function starts with an initial delay of 2 seconds and doubles the delay
//     with each retry attempt (exponential backoff).
//   - A random jitter is added to the delay to prevent synchronized retries in distributed systems.
//   - If the function succeeds (returns nil), the retry loop exits early.
//   - If the context is canceled or times out, the function returns an error indicating a timeout.
//   - If all retries are exhausted, the function returns an error indicating failure.
//
// Returns:
//   - nil if the function succeeds within the retry attempts.
//   - An error if the retries are exhausted or the context is canceled.
func exponentialBackOffRetry(ctx context.Context, e *email, fn func(e *email) error, retry int) error {
	delay := 2 * time.Second
	for i := 1; i <= retry; i++ {
		err := fn(e)
		if err == nil {
			return nil
		} else {
			fmt.Println("err")
		}
		backoff := delay * time.Duration(math.Pow(2, float64(i)))
		jitter := time.Duration(rand.Int63n(int64(i)))
		sleep := backoff + jitter

		select {
		case <-time.After(sleep):
			fmt.Println("send email retry ", i)
		case <-ctx.Done():
			return fmt.Errorf("send email timeout")
		}
	}
	return fmt.Errorf("retries are exhausted")
}
