package goemail

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/gomail.v2"
)

// ContentHTML represents the content type for HTML emails.
const (
	ContentHTML = "html"
	ContentText = "text"
)

// email represents the structure for composing and sending an email.
// It contains fields for the email's content, sender details, recipients,
// and additional metadata such as subject, tags, and embedded images.
type email struct {
	// dialer is used to establish a connection to the SMTP server for sending emails.
	dialer *gomail.Dialer

	// content holds the body of the email.
	content string

	// contentType specifies the MIME type of the email content (e.g., "text/plain" or "text/html").
	contentType string

	// sender is the email address of the sender.
	sender string

	// senderName is the name of the sender to be displayed in the email.
	senderName string

	// returnEmail is the email address to which replies should be sent.
	returnEmail string

	// subject is the subject line of the email.
	subject string

	// tags are optional metadata tags associated with the email.
	tags string

	// to is a list of recipient email addresses.
	to []string

	// cc is a list of email addresses to be included in the CC (carbon copy) field.
	cc []string

	// bcc is a list of email addresses to be included in the BCC (blind carbon copy) field.
	bcc []string

	// imagesToEmbed is a list of file paths for images to be embedded in the email.
	imagesToEmbed []string

	// attachments is a list of file paths for images to be attached in the email.
	attachments []string
}

// MailDetails represents the details of an email to be sent.
// It includes information about the recipients, sender, subject, and content.
//
// Fields:
// - To: A list of recipient email addresses.
// - Cc: A list of email addresses to be included in the CC (carbon copy) field.
// - Bcc: A list of email addresses to be included in the BCC (blind carbon copy) field.
// - Subject: The subject line of the email.
// - Sender: The email address of the sender.
// - SenderName: The name of the sender to be displayed in the email.
// - ReturnEmail: The email address to which replies should be sent.
// - Tags: Tags associated with the email for categorization or tracking purposes.
// - Content: The main content of the email, represented by a Content struct.
// - ImagesToEmbed: A list of file paths for images to be embedded in the email.
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

// Content represents the structure for email content.
// It includes the type of content, the actual content,
// an optional template path for rendering, and any associated data.
//
// Fields:
// - Type: Specifies the type of the content (e.g., "text/plain", "text/html").
// - Content: The actual content of the email as a string.
// - TemplatePath: Path to the template file used for rendering the content (optional).
// - Data: Data to be used for populating the template (if applicable).
type Content struct {
	Type         string
	Content      string
	TemplatePath string
	Data         interface{}
}

// ConnectionDetails holds the configuration details required to establish
// a connection to an SMTP server for sending emails. It includes the
// server host, port, and authentication credentials.
type ConnectionDetails struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
}

// sender defines an interface for sending emails.
// It requires the implementation of a Send method, which takes a context
// and an email as parameters and returns an error if the sending process fails.
type sender interface {
	Send(ctx context.Context, e *email) error
}

// New creates a new email instance using the provided connection details.
// It validates the connection details and initializes a dialer for sending emails.
//
// Parameters:
//   - c: ConnectionDetails containing SMTP host, port, user, and password.
//
// Returns:
//   - *email: A pointer to the initialized email instance.
//   - error: An error if the connection details validation fails or any other issue occurs.
func New(c ConnectionDetails) (*email, error) {
	err := c.validate()
	if err != nil {
		return nil, err
	}
	dialer = gomail.NewDialer(c.SmtpHost, c.SmtpPort, c.SmtpUser, c.SmtpPassword)
	return &email{dialer: dialer}, nil
}

// validate checks the ConnectionDetails struct for missing or invalid fields.
// It ensures that the SmtpHost, SmtpPort, SmtpPassword, and SmtpUser fields
// are properly set. If any of these fields are invalid, an error string is
// constructed and returned. Otherwise, it returns nil.
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

// validate checks the fields of the email struct for validity.
// It returns an error if any of the following conditions are met:
// - The content field is empty.
// - The subject field exceeds 500 characters.
// - The tags field exceeds 500 characters.
// - The sender field exceeds 500 characters.
// - The senderName field exceeds 500 characters.
// - The returnEmail field exceeds 500 characters.
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

// sendBy sends an email using the provided sender implementation.
//
// Parameters:
//   - ctx: The context for managing request deadlines, cancellations, and other values.
//   - s: The sender implementation responsible for sending the email.
//   - e: A pointer to the email object containing the email details.
//
// Returns:
//   - An error if the email could not be sent, or nil if the operation was successful.
func sendBy(ctx context.Context, s sender, e *email) error {
	return s.Send(ctx, e)
}

// Send sends the email using the specified context and retry count.
// It utilizes a gomailSender to handle the actual sending process.
//
// Parameters:
//   - ctx: The context to control the lifetime of the send operation.
//   - retry: The number of retry attempts in case of failure.
//
// Returns:
//   - An error if the email sending fails, or nil if successful.
func (e *email) Send(ctx context.Context, retry int) error {
	return sendBy(ctx, &gomailSender{retry: retry}, e)
}

// SetTo sets the recipient email addresses for the email.
// It takes a slice of strings representing the email addresses
// and assigns them to the "to" field of the email.
func (e *email) SetTo(to []string) {
	e.to = to
}

// SetCC sets the CC (carbon copy) recipients for the email.
// It accepts a slice of strings representing the email addresses
// to be included in the CC field.
//
// Parameters:
//   - cc: A slice of strings containing the email addresses to set as CC.
func (e *email) SetCC(cc []string) {
	e.cc = cc
}

// SetBCC sets the BCC (Blind Carbon Copy) recipients for the email.
// The BCC recipients will receive the email without their addresses being visible to other recipients.
//
// Parameters:
//   - bcc: A slice of strings containing the email addresses to be added as BCC recipients.
func (e *email) SetBCC(bcc []string) {
	e.bcc = bcc
}

// SetTags sets the tags for the email instance.
// It takes a string parameter `tags` which represents the tags to be associated with the email.
func (e *email) SetTags(tags string) {
	e.tags = tags
}

// SetSubject sets the subject of the email.
// It takes a string parameter 'subject' which represents the subject line
// to be assigned to the email.
func (e *email) SetSubject(subject string) {
	e.subject = subject
}

// SetSender sets the sender's email address for the email.
// The sender parameter should be a valid email address in string format.
// This method updates the sender field of the email instance.
func (e *email) SetSender(sender string) {
	e.sender = sender
}

// SetSenderName sets the sender's name for the email.
// It updates the senderName field of the email instance.
//
// Parameters:
//   - senderName: A string representing the name of the sender.
func (e *email) SetSenderName(senderName string) {
	e.senderName = senderName
}

// SetReturnEmail sets the return email address for the email instance.
// This email address is used as the "Return-Path" in the email header,
// which specifies where bounce messages should be sent.
//
// Parameters:
//   - returnEmail: A string representing the return email address.
func (e *email) SetReturnEmail(returnEmail string) {
	e.returnEmail = returnEmail
}

// SetImagesToEmbed sets the list of image file paths to be embedded in the email.
// The provided slice of strings should contain the file paths of the images
// that need to be included as embedded content in the email.
//
// Parameters:
//
//	images - A slice of strings representing the file paths of the images to embed.
func (e *email) SetImagesToEmbed(images []string) {
	e.imagesToEmbed = images
}

// SetAttachements sets the list of image file paths to be attached in the email.
// The provided slice of strings should contain the file paths of the images
// that need to be attached in the email.
//
// Parameters:
//
//	images - A slice of strings representing the file paths of the files to attached.
func (e *email) SetAttachements(files []string) {
	e.attachments = files
}

// SetContent sets the content of the email based on the provided Content struct.
// It supports both HTML and plain text content types. For HTML content, it can
// either use the provided content string or process a template file with the
// given data. If neither content nor a valid template path is provided for HTML,
// an error is returned.
//
// Parameters:
//   - content: A Content struct containing the type, content string, template path,
//     and optional data for template processing.
//
// Returns:
//   - error: An error is returned if the content type is HTML and neither content
//     nor a valid template path is provided, or if there is an issue processing
//     the template.
func (e *email) SetContent(content Content) error {
	e.contentType = content.Type
	if content.Type == ContentHTML {

		if content.Content != "" {
			e.content = content.Content
		} else if content.TemplatePath != "" {
			c, err := processTemplate(content.TemplatePath, content.Data)
			if err != nil {
				return err
			}
			e.content = c
		} else {
			return fmt.Errorf("not content provided")
		}

	} else if content.Type == ContentText {
		e.content = content.Content
	}
	return nil
}
