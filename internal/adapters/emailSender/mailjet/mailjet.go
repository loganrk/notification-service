package mailjet

import (
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go"
)

type MailjetEmailSender struct {
	Client   *mailjet.Client
	From     string
	FromName string
}

func New(apiKey, apiSecret, from, fromName string) *MailjetEmailSender {
	client := mailjet.NewMailjetClient(apiKey, apiSecret)
	return &MailjetEmailSender{
		Client:   client,
		From:     from,
		FromName: fromName,
	}
}

func (m *MailjetEmailSender) SendEmail(to, subject, body string) error {
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
	fmt.Println("messages", messages)
	resp, err := m.Client.SendMailV31(&messages)
	if err != nil {
		return err
	}
	fmt.Println("resp", resp)

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
