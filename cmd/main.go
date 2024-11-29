package main

import (
	"log"

	goemail "github.com/Priyokumar/goemail"
)

func main() {

	goemail.Init(goemail.ConnectionDetails{
		Host:         "w",
		Port:         3,
		SmtpUser:     "w",
		SmtpPassword: "e",
	})

	m, err := goemail.New(goemail.MailDetails{
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
