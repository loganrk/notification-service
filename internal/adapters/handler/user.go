package handler

import "context"

func (h *handler) Activation(ctx context.Context, message []byte) error {
	h.logger.Infow(ctx, "Received Activation Message", "message", string(message))
	// TODO: Process the activation message
	return nil
}

func (h *handler) PasswordReset(ctx context.Context, message []byte) error {
	h.logger.Infow(ctx, "Received Password Reset Message", "message", string(message))
	// TODO: Process the password reset message
	return nil
}

func (h *handler) ActivationError(ctx context.Context, err error) {
	h.logger.Errorw(ctx, "Error in Activation Consumer", "error", err)
}

func (h *handler) PasswordResetError(ctx context.Context, err error) {
	h.logger.Errorw(ctx, "Error in Password Reset Consumer", "error", err)
}
