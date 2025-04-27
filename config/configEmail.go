package config

type Email interface {
	GetSMTPFrom() string
	GetSMTPPassword() string
	GetSMTPHost() string
	GetSMTPPort() int
}

func (e email) GetSMTPFrom() string {
	return e.SMTP.From
}

func (e email) GetSMTPPassword() string {
	return e.SMTP.Password
}

func (e email) GetSMTPHost() string {
	return e.SMTP.Host
}

func (e email) GetSMTPPort() int {
	return e.SMTP.Port
}
