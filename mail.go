package goemail

import (
	"context"
	"errors"

	"gopkg.in/gomail.v2"
)

const (
	ContentHTML = "html"
	ContentText = "text"
)

type email struct {
	dialer        *gomail.Dialer
	content       string
	contentType   string
	sender        string
	senderName    string
	returnEmail   string
	subject       string
	tags          string
	to            []string
	cc            []string
	bcc           []string
	imagesToEmbed []string
}

type MailDetails struct {
	To      []string
	Cc      []string
	Bcc     []string
	Subject string

	Sender      string
	SenderName  string
	ReturnEmail string
	Tags        string

	Content       Content
	ImagesToEmbed []string
}

type Content struct {
	Type         string
	Content      string
	TemplatePath string
	Data         interface{}
}

type ConnectionDetails struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
}

type sender interface {
	Send(ctx context.Context, e *email, retry int) error
}

func New(c ConnectionDetails) (*email, error) {
	err := c.validate()
	if err != nil {
		return nil, err
	}
	dialer = gomail.NewDialer(c.SmtpHost, c.SmtpPort, c.SmtpUser, c.SmtpPassword)
	return &email{dialer: dialer}, nil
}

func (c ConnectionDetails) validate() error {
	errStr := ""
	if c.SmtpHost == "" {
		errStr += "host cannot be empty"
	}
	if c.SmtpPort == 0 {
		errStr += "port cannot be 0"
	}
	if c.SmtpPassword == "" {
		errStr += "SmtpPassword cannot be empty"
	}
	if c.SmtpUser == "" {
		errStr += "SmtpUser cannot be empty"
	}

	return nil
}

func (e *email) validate() error {

	if e.content == "" {
		return errors.New("content missing")
	}

	if e.subject != "" && len(e.subject) > 500 {
		return errors.New("subject too long")
	}

	if e.tags != "" && len(e.tags) > 500 {
		return errors.New("tags too long")
	}

	if e.sender != "" && len(e.sender) > 500 {
		return errors.New("sender email too long")
	}

	if e.senderName != "" && len(e.senderName) > 500 {
		return errors.New("sender name email too long")
	}

	if e.returnEmail != "" && len(e.returnEmail) > 500 {
		return errors.New("returnEmail too long")
	}
	return nil
}

func sendBy(ctx context.Context, s sender, e *email, retry int) error {
	return s.Send(ctx, e, retry)
}

func (e *email) Send(ctx context.Context, retry int) error {
	return sendBy(ctx, &gomailSender{}, e, retry)
}

func (e *email) SetTo(to []string) {
	e.to = to
}

func (e *email) SetCC(cc []string) {
	e.cc = cc
}

func (e *email) SetBCC(bcc []string) {
	e.bcc = bcc
}

func (e *email) SetTags(tags string) {
	e.tags = tags
}

func (e *email) SetSubject(subject string) {
	e.subject = subject
}

func (e *email) SetSender(sender string) {
	e.sender = sender
}

func (e *email) SetSenderName(senderName string) {
	e.senderName = senderName
}

func (e *email) SetReturnEmail(returnEmail string) {
	e.returnEmail = returnEmail
}

func (e *email) SetImagesToEmbed(images []string) {
	e.imagesToEmbed = images
}

func (e *email) SetContent(content Content) error {
	if content.Type == ContentHTML {
		c, err := processTemplate(content.TemplatePath, content.Data)
		if err != nil {
			return err
		}
		e.content = c
		e.contentType = ContentHTML
	} else if content.Type == ContentText {
		e.content = content.Content
	}
	return nil
}
