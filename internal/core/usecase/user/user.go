package user

import (
	"context"

	"github.com/loganrk/message-service/internal/core/domain"
	"github.com/loganrk/message-service/internal/core/port"
)

type userusecase struct {
	logger port.Logger
}

// New initializes a new user service with required dependencies and returns it.
func New(loggerIns port.Logger) *userusecase {
	return &userusecase{
		logger: loggerIns,
	}
}
func (u *userusecase) Activation(ctx context.Context, req domain.UserActivation) error {
	u.logger.Infow(ctx, "Processing Activation Email",
		"to", req.To,
		"subject", req.Subject,
		"macros", req.Macros,
	)

	// TODO: Here you can actually call your EmailService to send email
	// Example:
	// err := u.emailSender.SendActivationEmail(req.To, req.Subject, req.Macros["name"], req.Macros["activationLink"])
	// if err != nil {
	//     u.logger.Errorw(ctx, "Failed to send activation email", "error", err)
	//     return err
	// }

	return nil
}

func (u *userusecase) PasswordReset(ctx context.Context, req domain.UserPasswordReset) error {
	u.logger.Infow(ctx, "Processing Password Reset Email",
		"to", req.To,
		"subject", req.Subject,
		"macros", req.Macros,
	)

	// TODO: Here you can actually call your EmailService to send email
	// Example:
	// err := u.emailSender.SendPasswordResetEmail(req.To, req.Subject, req.Macros["name"], req.Macros["activationLink"])
	// if err != nil {
	//     u.logger.Errorw(ctx, "Failed to send password reset email", "error", err)
	//     return err
	// }

	return nil
}
