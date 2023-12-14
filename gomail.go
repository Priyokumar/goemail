package goemail

import (
	"errors"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

type gomailSender struct{}

func (g gomailSender) Send(e *email) error {

	gm := getMessage(e)
	dialer := gomail.NewDialer(host, port, smtpUser, smtpPassword)

	for i := 0; i < 5; i++ {
		err := dialer.DialAndSend(gm)
		if err != nil {
			log.Println("failed to send email. ", err.Error())
			log.Println("Retry ", i+1)
		} else {
			log.Println("sent email")
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return errors.New("failed to sent email")
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
