package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"notification/config"
	"path/filepath"
	"runtime"
	"strings"
)

type EmailNotifier struct {
	configs string
}

func NewEmailNotifier(configs string)*EmailNotifier {
	return &EmailNotifier{
		configs: configs,
	}
}

func (s EmailNotifier) NotifyUserRegistered(ctx context.Context, name string, email string) error {
		msg := newEmailMessage(
		"registration.html",
		"Welcome",
		struct{ Name string }{Name: name},
		[]string{email},
	)

	err := s.sendEmail(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s EmailNotifier) NotifyUserUpdated(ctx context.Context, name string, email string) error {
	msg := newEmailMessage(
		"update.html",
		"update",
		struct{ Name string }{Name: name},
		[]string{email},
	)

	err := s.sendEmail(msg)
	if err != nil {
		return err
	}

	return nil
}

func templatesDirPath() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error in generating env dir")
	}

	return filepath.Dir(f)
}

type emailMessage struct {
	tmplateFileName string
	receiver        []string
	args            interface{}
	subject         string
}


func newEmailMessage(tmplFileName, subject string, args interface{}, reciver []string) *emailMessage {
	return &emailMessage{
		tmplateFileName: tmplFileName,
		subject:         subject,
		args:            args,
		receiver:        reciver,
	}
}

func (s *EmailNotifier) readTemplate(templFileName string) *template.Template {
	templFileNameFullAddres := fmt.Sprintf("%s/%s/%s", templatesDirPath(), "templates", templFileName)

	tpl := template.Must(template.ParseFiles(templFileNameFullAddres))

	return tpl
}

func (s *EmailNotifier) sendEmail(msg *emailMessage) error {
	if config.GetENV() != "PRODUCTION" {
		log.Println("Email sent")
		log.Println(msg)
		return nil
	}

	template := s.readTemplate(msg.tmplateFileName)

	var body bytes.Buffer
	err := template.Execute(&body, msg.args)

	if err != nil {
		return err
	}

	emailMessage := fmt.Sprintf("Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+body.String(), msg.subject)

	configs := extractSMTPConfigs(s.configs)
	auth := smtp.PlainAuth("", configs[0], configs[1], configs[2])
	err = smtp.SendMail(configs[2]+":"+configs[3], auth, configs[0], msg.receiver, []byte(emailMessage))

	if err != nil {
		return err
	}

	log.Println("email successfully sent")

	return nil
}

func extractSMTPConfigs(configs string) []string {
	GetSMTPConfigs := strings.Split(configs, ",")

	return GetSMTPConfigs
}
