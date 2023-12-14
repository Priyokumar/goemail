package goemail

import (
	"bytes"
	"errors"
	"html/template"
	"log"
)

const ContentHTML = "html"

var (
	host         string
	port         int
	smtpUser     string
	smtpPassword string
)

type email struct {
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

	Template      Template
	ImagesToEmbed []string

	Content     string
	ContentType string
}

type Template struct {
	Multilevel    bool
	TemplatePaths []string
	Data          interface{}
}

type ConnectionDetails struct {
	Host         string
	Port         int
	SmtpUser     string
	SmtpPassword string
}

type sender interface {
	Send(e *email) error
}

func Init(con ConnectionDetails) {
	host = con.Host
	port = con.Port
	smtpPassword = con.SmtpPassword
	smtpUser = con.SmtpUser
}

func New(mail MailDetails) (*email, error) {
	if host == "" || port == 0 {
		return nil, errors.New("call init func")
	}
	e := new(email)
	e.contentType = mail.ContentType
	if e.contentType == ContentHTML {
		c, err := processTemplate(mail.Template)
		if err != nil {
			return nil, err
		}
		e.content = c
		e.contentType = ContentHTML
	} else {
		e.content = mail.Content
	}

	e.bcc = mail.Bcc
	e.cc = mail.Cc
	e.to = mail.To
	e.tags = mail.Tags
	e.imagesToEmbed = mail.ImagesToEmbed
	e.returnEmail = mail.ReturnEmail
	e.sender = mail.Sender
	e.senderName = mail.SenderName
	e.subject = mail.Subject

	err := e.validate()
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *email) validate() error {

	// TODO
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

func sendBy(e *email, s sender) error {
	return s.Send(e)
}

func (e *email) Send() error {
	return sendBy(e, gomailSender{})
}

func processTemplate(template Template) (string, error) {
	log.Println("process template")
	if !template.Multilevel {
		return processSingleTemplate(template.TemplatePaths[0], template.Data)
	} else {
		return processMultilevelTemplate(template.TemplatePaths, template.Data)
	}
}

func processSingleTemplate(templatePath string, data interface{}) (string, error) {
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

func processMultilevelTemplate(templatePath []string, data interface{}) (string, error) {
	return "", nil
}
