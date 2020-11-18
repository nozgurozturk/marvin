package app

import (
	"bytes"
	"github.com/nozgurozturk/marvin/server/internal/config"
	"html/template"
	"net/smtp"
)

func SendEmail(to string, subject string, body string) error {
	cnf := config.Get().SMTP
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	emailFrom := "From: " + cnf.From + "\r\n"
	emailTo := "To: " + to + "\r\n"
	sbj := "Subject: " + subject + "\r\n"
	msg := []byte(emailFrom + emailTo + sbj + mime + "\r\n" + body)

	err := smtp.SendMail(cnf.Host+cnf.Port,
		smtp.PlainAuth("", cnf.From, cnf.Password, cnf.Host),
		cnf.From, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}

func ParseHTMLTemplate(templateFileName string, data interface{}) (string, error) {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "",err
	}
	return buf.String(), nil
}