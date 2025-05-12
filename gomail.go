package goemail

import (
	"context"

	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

type gomailSender struct {
	retry int
	ctx   context.Context
}

func (g *gomailSender) Send(ctx context.Context, e *email) error {
	err := e.validate()
	if err != nil {
		return err
	}
	return exponentialBackOffRetry(g.ctx, e, sendEmail, g.retry)
}

func sendEmail(e *email) error {
	gm := getMessage(e)
	return e.dialer.DialAndSend(gm)
}

func getMessage(e *email) *gomail.Message {
	gm := gomail.NewMessage()

	if e.contentType == ContentHTML {
		gm.SetBody("text/html", e.content)
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
