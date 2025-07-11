package user

import (
	"context"
	"fmt"
	"os"

	"github.com/loganrk/worker-engine/config"
	"github.com/loganrk/worker-engine/internal/core/port"
	"github.com/loganrk/worker-engine/internal/utils"
)

// userusecase implements user-related operations such as sending activation and password reset emails.
type userusecase struct {
	logger           port.Logger      // Logger interface for structured logging
	activationTpl    string           // Email template content for activation emails
	passwordResetTpl string           // Email template content for password reset emails
	emailSender      port.EmailSender // Interface to send emails
}

// New initializes a new userusecase instance by loading email templates and setting dependencies.
func New(loggerIns port.Logger, emailSenderIns port.EmailSender, userConf config.User) (*userusecase, error) {
	// Read activation email template from file
	activationTpl, err := os.ReadFile(userConf.GetActivationTemplatePath())
	if err != nil {
		return nil, fmt.Errorf("failed to load activation template: %w", err)
	}

	// Read password reset email template from file
	passwordResetTpl, err := os.ReadFile(userConf.GetPasswordResetTemplatePath())
	if err != nil {
		return nil, fmt.Errorf("failed to load password reset template: %w", err)
	}

	// Return the fully initialized userusecase
	return &userusecase{
		logger:           loggerIns,
		emailSender:      emailSenderIns,
		activationTpl:    string(activationTpl),
		passwordResetTpl: string(passwordResetTpl),
	}, nil
}
func (u *userusecase) ActivationEmail(ctx context.Context, to, subject string, macros map[string]string) error {
	u.logger.Infow(ctx, "Processing Activation Email", "to", to, "subject", subject, "macros", macros)

	emailBody := utils.ReplaceMacros(u.activationTpl, macros)
	if err := u.emailSender.SendEmail(to, subject, emailBody); err != nil {
		u.logger.Errorw(ctx, "Failed to send activation email", "error", err)
		return err
	}
	return nil
}

func (u *userusecase) ActivationPhone(ctx context.Context, to string, macros map[string]string) error {
	u.logger.Infow(ctx, "Processing Activation SMS", "to", to, "macros", macros)

	// message := utils.ReplaceMacros(u.activationSMSTpl, macros)
	// if err := u.smsSender.SendSMS(to, message); err != nil {
	// 	u.logger.Errorw(ctx, "Failed to send activation SMS", "error", err)
	// 	return err
	// }
	return nil
}

func (u *userusecase) PasswordResetEmail(ctx context.Context, to, subject string, macros map[string]string) error {
	u.logger.Infow(ctx, "Processing Password Reset Email", "to", to, "subject", subject, "macros", macros)

	emailBody := utils.ReplaceMacros(u.passwordResetTpl, macros)
	if err := u.emailSender.SendEmail(to, subject, emailBody); err != nil {
		u.logger.Errorw(ctx, "Failed to send password reset email", "error", err)
		return err
	}
	return nil
}

func (u *userusecase) PasswordResetPhone(ctx context.Context, to string, macros map[string]string) error {
	u.logger.Infow(ctx, "Processing Password Reset SMS", "to", to, "macros", macros)

	// message := utils.ReplaceMacros(u.passwordResetSMSTpl, macros)
	// if err := u.smsSender.SendSMS(to, message); err != nil {
	// 	u.logger.Errorw(ctx, "Failed to send password reset SMS", "error", err)
	// 	return err
	// }
	return nil
}
