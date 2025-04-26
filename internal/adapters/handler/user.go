package handler

import (
	"context"
	"encoding/json"

	"github.com/loganrk/message-service/internal/core/domain"
)

func (h *handler) Activation(ctx context.Context, message []byte) error {
	h.logger.Infow(ctx, "Received Activation Message", "message", string(message))

	var req domain.UserActivation
	if err := json.Unmarshal(message, &req); err != nil {
		h.logger.Errorw(ctx, "Failed to unmarshal activation message", "error", err)
		return err
	}

	// Call usecase
	if err := h.usecases.User.Activation(ctx, req); err != nil {
		h.logger.Errorw(ctx, "Failed to process activation usecase", "error", err)
		return err
	}

	h.logger.Infow(ctx, "Successfully processed Activation Message", "to", req.To)
	return nil
}

func (h *handler) PasswordReset(ctx context.Context, message []byte) error {
	h.logger.Infow(ctx, "Received Password Reset Message", "message", string(message))

	var req domain.UserPasswordReset
	if err := json.Unmarshal(message, &req); err != nil {
		h.logger.Errorw(ctx, "Failed to unmarshal password reset message", "error", err)
		return err
	}

	// Call usecase
	if err := h.usecases.User.PasswordReset(ctx, req); err != nil {
		h.logger.Errorw(ctx, "Failed to process password reset usecase", "error", err)
		return err
	}

	h.logger.Infow(ctx, "Successfully processed Password Reset Message", "to", req.To)
	return nil
}

func (h *handler) ActivationError(ctx context.Context, err error) {
	h.logger.Errorw(ctx, "Error in Activation Consumer", "error", err)
}

func (h *handler) PasswordResetError(ctx context.Context, err error) {
	h.logger.Errorw(ctx, "Error in Password Reset Consumer", "error", err)
}
