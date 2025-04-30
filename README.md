# 📧 GoEmail – Simple and Fast Email Sending in Go

GoeEmail is a lightweight and easy-to-use package for sending emails using Go (Golang). It supports SMTP authentication, HTML/plain text messages, attachments, and more – all with a clean and developer-friendly API.

---

## ✨ Features

- 📤 Send plain text or HTML emails
- ✅ Supports SMTP authentication (username/password)
- 💡 Easy integration with your Go applications

---

## 🚀 Getting Started

### Installation

```
go get github.com/Priyokumar/goemail
```

### Example:
```
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
 ```


