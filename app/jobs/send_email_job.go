package jobs

import (
	"github.com/test/myapp/framework/queue"
	"log"
)

// SendEmailJob handles email sending
type SendEmailJob struct {
	queue.BaseJob
	Email   string
	Subject string
	Body    string
}

func (j *SendEmailJob) Handle() error {
	log.Printf("📧 Sending email to %s: %s", j.Email, j.Subject)
	// Email sending logic here
	return nil
}

func NewSendEmailJob(email, subject, body string) *SendEmailJob {
	return &SendEmailJob{
		BaseJob: queue.BaseJob{
			Name: "send_email",
			Payload: map[string]interface{}{
				"email":   email,
				"subject": subject,
				"body":    body,
			},
		},
		Email:   email,
		Subject: subject,
		Body:    body,
	}
}
