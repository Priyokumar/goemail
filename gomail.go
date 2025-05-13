package goemail

import (
	"context"

	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

// gomailSender represents a mail sender with a configurable retry mechanism.
// The retry field specifies the number of times the sender will attempt to resend
// an email in case of failure.
type gomailSender struct {
	retry int
}

// Send attempts to send an email using the gomailSender. It first validates the email
// and then retries sending it using an exponential backoff strategy if necessary.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellations.
//   - e: A pointer to the email object to be sent.
//
// Returns:
//   - An error if the email validation fails or if the email cannot be sent
//     after the configured retries.
func (g *gomailSender) Send(ctx context.Context, e *email) error {
	err := e.validate()
	if err != nil {
		return err
	}
	return exponentialBackOffRetry(ctx, e, sendEmail, g.retry)
}

// sendEmail sends an email using the provided email configuration.
// It constructs the email message using the getMessage function and
// sends it via the dialer associated with the email instance.
//
// Parameters:
//   - e: A pointer to an email instance containing the email details
//     and the dialer for sending the email.
//
// Returns:
//   - error: An error if the email fails to send, or nil if the email
//     is sent successfully.
func sendEmail(e *email) error {
	gm := getMessage(e)
	return e.dialer.DialAndSend(gm)
}

// getMessage constructs and returns a *gomail.Message based on the provided email struct.
// It sets the email body, headers, and optionally embeds images.
//
// Parameters:
//   - e: A pointer to an email struct containing the email details such as content, recipients, subject, etc.
//
// Behavior:
//   - Sets the email body based on the content type (HTML or plain text).
//   - Configures the "To", "Cc", and "Bcc" headers if recipients are provided.
//   - Sets the "Subject", "From", "X-SES-MESSAGE-TAGS", and "Return-Path" headers.
//   - Embeds images into the email if any are specified.
//
// Returns:
//   - A pointer to a gomail.Message object representing the constructed email.
func getMessage(e *email) *gomail.Message {
	gm := gomail.NewMessage()

	if e.contentType == ContentHTML {
		gm.SetBody("text/html", e.content)
	} else if e.contentType == ContentText {
		gm.SetBody("text/plain", e.content)
	}

	if len(e.to) > 0 {
		a := make([]string, 0)
		for _, v := range e.to {
			a = append(a, gm.FormatAddress(v, ""))
		}
		gm.SetHeader("To", a...)
	}

	if len(e.cc) > 0 {
		a := make([]string, 0)
		for _, v := range e.cc {
			a = append(a, gm.FormatAddress(v, ""))
		}
		gm.SetHeader("Cc", a...)
	}

	if len(e.bcc) > 0 {
		a := make([]string, 0)
		for _, v := range e.bcc {
			a = append(a, gm.FormatAddress(v, ""))
		}
		gm.SetHeader("Bcc", a...)
	}

	gm.SetHeader("Subject", e.subject)
	gm.SetHeader("From", gm.FormatAddress(e.sender, e.senderName))
	gm.SetHeader("X-SES-MESSAGE-TAGS", e.tags)
	gm.SetHeader("Return-Path", e.returnEmail)

	if len(e.imagesToEmbed) > 0 {
		for _, v := range e.imagesToEmbed {
			gm.Embed(v)
		}
	}
	return gm
}
