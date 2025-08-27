package email

import (
	"bytes"
	"text/template"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type IEmail interface {
	SetSender(sender string)
	SetReceiver(to ...string)
	SetSubject(subject string)
	SetBodyHTML(path string, data interface{}) error
	Send() error
}

type Email struct {
	message  *gomail.Message
	dialer   *gomail.Dialer
	htmlPath string
}

func New() IEmail {
	return &Email{
		message: gomail.NewMessage(),
		dialer: gomail.NewDialer(
			viper.GetString("email.host"),
			viper.GetInt("email.port"),
			viper.GetString("email.user"),
			viper.GetString("email.password"),
		),
		htmlPath: viper.GetString("email.html_template_path"),
	}
}

func (g *Email) SetSender(sender string) {
	g.message.SetHeader("From", sender)
}

func (g *Email) SetReceiver(to ...string) {
	g.message.SetHeader("To", to...)
}

func (g *Email) SetSubject(subject string) {
	g.message.SetHeader("Subject", subject)
}

func (g *Email) SetBodyHTML(path string, data interface{}) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(g.htmlPath + path)
	if err != nil {
		return errx.InternalServerError("failed to parse template")
	}

	err = t.Execute(&body, data)
	if err != nil {
		return errx.InternalServerError("failed to execute template")
	}

	g.message.SetBody("text/html", body.String())
	return nil
}

func (g *Email) Send() error {
	if err := g.dialer.DialAndSend(g.message); err != nil {
		return err
	}
	return nil
}
