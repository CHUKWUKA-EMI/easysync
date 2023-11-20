package services

import (
	"context"
	"log"
	"net/http"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

// EmailUser ...
type EmailUser struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

// EmailService ...
type EmailService struct {
	Sender    EmailUser `binding:"required"`
	Recipient EmailUser `binding:"required"`
	Subject   string
	Content   string
}

// SendEmail ...
func (e EmailService) SendEmail() *http.Response {
	var ctx context.Context
	config := brevo.NewConfiguration()
	config.AddDefaultHeader("api-key", os.Getenv("BREVO_API_KEY"))

	client := brevo.NewAPIClient(config)

	_, response, err := client.TransactionalEmailsApi.SendTransacEmail(
		ctx,
		brevo.SendSmtpEmail{
			Sender:      &brevo.SendSmtpEmailSender{Name: e.Sender.Name, Email: e.Sender.Email},
			To:          []brevo.SendSmtpEmailTo{{Email: e.Recipient.Email, Name: e.Recipient.Name}},
			Subject:     e.Subject,
			HtmlContent: e.Content,
		})
	if err != nil {
		log.Fatal("Error sending email", err.Error())
		os.Exit(1)
	}

	return response
}
