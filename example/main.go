package main

import (
	"context"
	"log"
	"time"

	"github.com/Priyokumar/goemail"
)

func main() {

	m, err := goemail.New(
		goemail.ConnectionDetails{
			SmtpHost:     "xxxxxx",
			SmtpPort:     000,
			SmtpUser:     "xxxx",
			SmtpPassword: "xxx",
		})

	if err != nil {
		log.Println(err)
		return
	}
	m.SetTo([]string{"xxxxxx"})
	m.SetSubject("Test")
	m.SetSenderName("test")
	m.SetSender("xxxxxx")
	m.SetReturnEmail("xxxxxx")
	m.SetContent(
		goemail.Content{
			Type:    goemail.ContentHTML,
			Content: "<div><h2>Test</h2></div>",
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*50))
	defer cancel()
	err = m.Send(ctx, 10)

	if err != nil {
		log.Println(err)
		return
	}

}
