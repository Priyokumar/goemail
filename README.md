# 📧 GoEmail – Simple and Fast Email Sending in Go

GoeEmail is a lightweight and easy-to-use package for sending emails using Go (Golang). It supports SMTP authentication, HTML/plain text messages, attachments, and more – all with a clean and developer-friendly API.

---

## ✨ Features

- 📤 Send plain text or HTML emails
- ✅ Supports SMTP authentication (username/password)
- 💡 Easy integration with your Go applications
- Support retry logic
- Support context timeout

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

 ```

 For sending text emails, you can use the `ContentText` type:
```
m.SetContent(
		goemail.Content{
			Type:    goemail.ContentText,
			Content: "This is a test email",
		},
	)
```

 For attaching files to emails:
```
m.SetAttachments([]string{"path/to/your/file.txt"})
```


