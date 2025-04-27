package config

type User interface {
	GetActivationTemplatePath() string
	GetPasswordResetTemplatePath() string
}

func (u user) GetActivationTemplatePath() string {

	return u.Activation.TemplatePath
}

func (u user) GetPasswordResetTemplatePath() string {

	return u.PasswordReset.TemplatePath
}
