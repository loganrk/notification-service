package user

import (
	"context"
	"fmt"
	"os"

	"github.com/loganrk/message-service/config"
	"github.com/loganrk/message-service/internal/core/domain"
	"github.com/loganrk/message-service/internal/core/port"
	"github.com/loganrk/message-service/internal/utils"
)

type userusecase struct {
	logger           port.Logger
	activationTpl    string
	passwordResetTpl string
	emailSender      port.EmailSender
}

// New initializes a new user service with required dependencies and returns it.
func New(loggerIns port.Logger, emailSenderIns port.EmailSender, userConf config.User) (*userusecase, error) {

	activationTpl, err := os.ReadFile(userConf.GetActivationTemplatePath())
	if err != nil {
		return nil, fmt.Errorf("failed to load activation template: %w", err)
	}

	passwordResetTpl, err := os.ReadFile(userConf.GetPasswordResetTemplatePath())
	if err != nil {
		return nil, fmt.Errorf("failed to load password reset template: %w", err)
	}

	return &userusecase{
		logger:           loggerIns,
		emailSender:      emailSenderIns,
		activationTpl:    string(activationTpl),
		passwordResetTpl: string(passwordResetTpl),
	}, nil
}

func (u *userusecase) Activation(ctx context.Context, req domain.UserActivation) error {

	if req.Type == "activation-email" {

		u.logger.Infow(ctx, "Processing Activation Email",
			"to", req.To,
			"subject", req.Subject,
			"macros", req.Macros,
		)

		emailBody := utils.ReplaceMacros(u.activationTpl, req.Macros)

		err := u.emailSender.SendEmail(req.To, req.Subject, emailBody)
		if err != nil {
			u.logger.Errorw(ctx, "Failed to send activation email", "error", err)
			return err
		}

	}

	return nil
}

func (u *userusecase) PasswordReset(ctx context.Context, req domain.UserPasswordReset) error {

	if req.Type == "activation-email" {

		u.logger.Infow(ctx, "Processing Password Reset Email",
			"to", req.To,
			"subject", req.Subject,
			"macros", req.Macros,
		)

		emailBody := utils.ReplaceMacros(u.passwordResetTpl, req.Macros)

		err := u.emailSender.SendEmail(req.To, req.Subject, emailBody)
		if err != nil {
			u.logger.Errorw(ctx, "Failed to send password reset email", "error", err)
			return err
		}
	}

	return nil
}
