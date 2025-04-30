package main

import (
	"log"

	"github.com/Priyokumar/goemail"
)

func main() {

	m, err := goemail.New(
		goemail.ConnectionDetails{
			Host:         "w",
			Port:         3,
			SmtpUser:     "w",
			SmtpPassword: "e",
		},
		goemail.MailDetails{
			To:          []string{"priyon999@gmail.com"},
			Subject:     "Test",
			Sender:      "priyon999@gmail.com",
			SenderName:  "Test",
			ReturnEmail: "priyon999@gmail.com",
			ContentType: goemail.ContentHTML,
			Template: goemail.Template{
				Multilevel:    false,
				TemplatePaths: []string{"index.html"},
				Data:          nil,
			},
		})

	if err != nil {
		log.Println(err)
		return
	}

	err = m.Send()

	if err != nil {
		log.Println(err)
		return
	}

}
