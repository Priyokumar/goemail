package main

import (
	"context"
	"log"

	"github.com/Priyokumar/goemail"
)

func main() {

	m, err := goemail.New(
		goemail.ConnectionDetails{
			SmtpHost:     "w",
			SmtpPort:     3,
			SmtpUser:     "w",
			SmtpPassword: "e",
		})

	if err != nil {
		log.Println(err)
		return
	}
	m.SetTo([]string{"priyon999@gmail.com"})
	m.SetSubject("Test")
	m.SetSenderName("test")
	m.SetSender("test@test.com")
	m.SetReturnEmail("tets@test.com")
	m.SetContent(
		goemail.Content{
			Type:    goemail.ContentText,
			Content: "This is test email.",
		},
	)
	err = m.Send(context.Background(), 1)

	if err != nil {
		log.Println(err)
		return
	}

}
