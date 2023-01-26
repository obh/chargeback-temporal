package utils

import (
	"bytes"
	"html/template"
	"os"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Activities struct {
	SMTPHost     string
	SMTPPort     int
	SMTPStub     bool
	SMTPUser     string
	SMTPPassword string
}

var u, _ = os.LookupEnv("SMTP_USERNAME")
var p, _ = os.LookupEnv("SMTP_PASSWORD")
var a = &Activities{
	SMTPHost:     "smtp.mailtrap.io",
	SMTPPort:     2525,
	SMTPStub:     false,
	SMTPUser:     u,
	SMTPPassword: p,
}

func SendMail(from string, to string, subject string, htmlTemplate *template.Template,
	textTemplate *template.Template, input interface{}) error {
	var htmlContent, textContent bytes.Buffer

	err := htmlTemplate.Execute(&htmlContent, input)
	if err != nil {
		return err
	}

	err = textTemplate.Execute(&textContent, input)
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject).
		SetBody(mail.TextHTML, htmlContent.String()).
		AddAlternative(mail.TextPlain, textContent.String())

	if email.Error != nil {
		return email.Error
	}

	if a.SMTPStub {
		return nil
	}

	server := mail.NewSMTPClient()
	server.Host = a.SMTPHost
	server.Port = a.SMTPPort
	server.Username = a.SMTPUser
	server.Password = a.SMTPPassword
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		return err
	}

	return email.Send(client)
}
