package email

import (
	"context"
	"fmt"
	"time"

	"encore.app/monitor"
	"encore.dev/pubsub"
	"github.com/resend/resend-go/v2"
)

var secrets struct {
	// ResendAPI.
	ResendAPIKey string
}

type EmailParams struct {
	// Text is the Slack message text to send.
	Text    string
	Subject string
}

const (
	EMAIL_FROM = "noreply@nosyn.dev"
)

//encore:api private
func Resend(ctx context.Context, p *EmailParams) error {
	client := resend.NewClient(secrets.ResendAPIKey)

	params := &resend.SendEmailRequest{
		To:      []string{"biem97@gmail.com"},
		From:    EMAIL_FROM,
		Text:    p.Text,
		Subject: p.Subject,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	fmt.Println(sent.Id)
	return nil
}

var _ = pubsub.NewSubscription(monitor.TransitionTopic, "email-notification", pubsub.SubscriptionConfig[*monitor.TransitionEvent]{
	Handler: func(ctx context.Context, event *monitor.TransitionEvent) error {
		subject := "Uptime Notification"

		// Compose our message.
		msg := fmt.Sprintf("*%s is down at %s!", event.Site.URL, time.Now())
		if event.Up {
			msg = fmt.Sprintf("*%s is back up at %s!", event.Site.URL, time.Now())
		}

		// Send an email.
		return Resend(ctx, &EmailParams{Text: msg, Subject: subject})
	},
})
