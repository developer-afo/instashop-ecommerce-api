package service

import (
	"fmt"
	"time"

	"github.com/developer-afo/instashop-ecommerce-api/lib/config"
)

type SendEmailParams struct {
	To, Subject, Template string
	Variables             interface{}
}

type emailService struct {
	mail config.EmailInterface
}

type EmailServiceInterface interface {
	SendEmail(params SendEmailParams) error
}

func NewEmailService(mail config.EmailInterface) EmailServiceInterface {
	return &emailService{mail: mail}
}

func (e *emailService) SendEmail(params SendEmailParams) error {
	go func(p SendEmailParams) {
		err := e.mail.SendWithTemplate(p.To, p.Subject, fmt.Sprintf("templates/%s.html", p.Template), p.Variables)
		if err != nil {
			e.SetLogger(p.Subject, p.To, err.Error())
		}
	}(params)

	return nil
}

func (e *emailService) SetLogger(sub string, to string, message string) {
	currentTime := time.Now()

	logMessage := fmt.Sprintf("TimeStamp: %s, Message: %s, Subject: %s, To: %v",
		currentTime.Format("2006-01-02 15:04:05"), message, sub, to)
	fmt.Println(logMessage)
}
