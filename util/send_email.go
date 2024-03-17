package util

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	CONFIG_SMTP_HOST = "smtp.gmail.com"
	CONFIG_SMTP_PORT = 587
	MIME             = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type Request struct {
	to      []string
	subject string
	body    string
	Con     config.Config
	Log     *logrus.Logger
}

func NewRequest(to []string, subject string, Con config.Config, Log *logrus.Logger) *Request {
	return &Request{
		to:      to,
		subject: subject,
		Con:     Con,
		Log:     Log,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {
	body := "From: " + r.Con.SenderName + "\n" + "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", r.Con.AuthEmail, r.Con.AuthPassword, CONFIG_SMTP_HOST), r.Con.AuthEmail, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *Request) Send(templateName string, items interface{}) error {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		r.Log.Errorf("Failed to send the email to %s\n", r.to)
		return err
	}
	if ok := r.sendMail(); ok {
		r.Log.Errorf("Email has been sent to %s\n", r.to)
		return nil
	} else {
		r.Log.Errorf("Failed to send the email to %s\n", r.to)
		return err
	}
}
