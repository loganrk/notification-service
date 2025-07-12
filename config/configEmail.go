package config

import "time"

type Email interface {
	GetMailjetAPIKey() string
	GetMailjetAPISecret() string
	GetMailjetFromEmail() string
	GetMailjetFromName() string
	GetMailjetRateLimitEnabled() bool
	GetMailjetRateLimitMaxRequest() int
	GetMailjetRateLimitWindowSize() time.Duration
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
func (e email) GetMailjetRateLimitEnabled() bool {
	return e.Mailjet.RateLimit.Enabled
}

func (e email) GetMailjetRateLimitMaxRequest() int {
	return e.Mailjet.RateLimit.MaxRequests
}
func (e email) GetMailjetRateLimitWindowSize() time.Duration {
	return e.Mailjet.RateLimit.WindowSize
}
