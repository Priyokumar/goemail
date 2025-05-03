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

func exponentialBackOffRetry(ctx context.Context, e *email, fn func(e *email) error, retry int) error {
	delay := 2 * time.Second
	for i := 0; i < retry; i++ {
		err := fn(e)
		if err == nil {
			return nil
		}
		backoff := delay * time.Duration(math.Pow(2, float64(i)))
		jitter := time.Duration(rand.Int63n(int64(i)))
		sleep := backoff + jitter

		select {
		case <-time.After(sleep):
			fmt.Println("send email retrying")
		case <-ctx.Done():
			return fmt.Errorf("send email timeout")
		}
	}
	return fmt.Errorf("retries are exhausted")
}
