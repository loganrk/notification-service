package mailjet

import (
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go"
)

type MailjetEmailer struct {
	Client   *mailjet.Client
	From     string
	FromName string
}

func New(apiKey, apiSecret, from, fromName string) *MailjetEmailer {
	client := mailjet.NewMailjetClient(apiKey, apiSecret)
	return &MailjetEmailer{
		Client:   client,
		From:     from,
		FromName: fromName,
	}
}

func (m *MailjetEmailer) SendEmail(to, subject, body string) error {
	toRecipients := mailjet.RecipientsV31{
		mailjet.RecipientV31{
			Email: to,
			Name:  to,
		},
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: m.From,
				Name:  m.FromName,
			},
			To:       &toRecipients,
			Subject:  subject,
			TextPart: body,
			HTMLPart: body,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	resp, err := m.Client.SendMailV31(&messages)
	if err != nil {
		return err
	}
	// Check how many messages were successfully sent
	if len(resp.ResultsV31) == 0 {
		return fmt.Errorf("no messages sent, response: %+v", resp)
	}

	// Optionally, check the status of each message
	for _, result := range resp.ResultsV31 {
		if result.Status != "success" && result.Status != "queued" {
			return fmt.Errorf("email sending failed, status: %s", result.Status)
		}
	}

	return nil
}
