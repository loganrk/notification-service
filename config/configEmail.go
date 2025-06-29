package config

type Email interface {
	GetMailjetAPIKey() string
	GetMailjetAPISecret() string
	GetMailjetFromEmail() string
	GetMailjetFromName() string
}

func (e email) GetMailjetAPIKey() string {
	return e.Mailjet.APIKey
}

func (e email) GetMailjetAPISecret() string {
	return e.Mailjet.APISecret
}

func (e email) GetMailjetFromEmail() string {
	return e.Mailjet.FromEmail
}

func (e email) GetMailjetFromName() string {
	return e.Mailjet.FromName
}
