package handler

import (
	"context"
)

// ActivationEmail processes a user activation via email.
func (h *handler) ActivationEmail(to, subject string, macros map[string]string) error {
	ctx := context.Background()

	err := h.usecases.User.ActivationEmail(ctx, to, subject, macros)
	if err != nil {
		h.logger.Errorw(ctx, "Failed to process Activation Email", "error", err)
		return err
	}

	h.logger.Infow(ctx, "Successfully processed Activation Email", "to", to)
	return nil
}

// ActivationPhone processes a user activation via phone (SMS).
func (h *handler) ActivationPhone(to string, macros map[string]string) error {
	ctx := context.Background()
	err := h.usecases.User.ActivationPhone(ctx, to, macros)
	if err != nil {
		h.logger.Errorw(ctx, "Failed to process Activation Phone", "error", err)
		return err
	}

	h.logger.Infow(ctx, "Successfully processed Activation Phone", "to", to)
	return nil
}

// PasswordResetEmail processes a password reset via email.
func (h *handler) PasswordResetEmail(to, subject string, macros map[string]string) error {
	ctx := context.Background()

	err := h.usecases.User.PasswordResetEmail(ctx, to, subject, macros)
	if err != nil {
		h.logger.Errorw(ctx, "Failed to process Password Reset Email", "error", err)
		return err
	}

	h.logger.Infow(ctx, "Successfully processed Password Reset Email", "to", to)
	return nil
}

// PasswordResetPhone processes a password reset via phone (SMS).
func (h *handler) PasswordResetPhone(to string, macros map[string]string) error {
	ctx := context.Background()

	err := h.usecases.User.PasswordResetPhone(ctx, to, macros)
	if err != nil {
		h.logger.Errorw(ctx, "Failed to process Password Reset Phone", "error", err)
		return err
	}

	h.logger.Infow(ctx, "Successfully processed Password Reset Phone", "to", to)
	return nil
}

// ActivationError logs errors that occur in the activation Kafka consumer pipeline.
func (h *handler) ActivationError(ctx context.Context, err error) {
	h.logger.Errorw(ctx, "Error in Activation Consumer", "error", err)
}

// PasswordResetError logs errors that occur in the password reset Kafka consumer pipeline.
func (h *handler) PasswordResetError(ctx context.Context, err error) {
	h.logger.Errorw(ctx, "Error in Password Reset Consumer", "error", err)
}
