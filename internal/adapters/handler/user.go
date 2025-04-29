package handler

import (
	"context"
	"encoding/json"

	"github.com/loganrk/message-service/internal/core/domain"
)

// Activation handles incoming user activation messages.
// It logs the received message, unmarshals it into a UserActivation struct,
// and then calls the corresponding usecase logic.
func (h *handler) Activation(ctx context.Context, message []byte) error {
	h.logger.Infow(ctx, "Received Activation Message", "message", string(message))

	// Parse the incoming JSON message into the UserActivation struct.
	var req domain.UserActivation
	if err := json.Unmarshal(message, &req); err != nil {
		h.logger.Errorw(ctx, "Failed to unmarshal activation message", "error", err)
		return err
	}

	// Call the business logic (usecase) to handle the activation.
	if err := h.usecases.User.Activation(ctx, req); err != nil {
		h.logger.Errorw(ctx, "Failed to process activation usecase", "error", err)
		return err
	}

	h.logger.Infow(ctx, "Successfully processed Activation Message", "to", req.To)
	return nil
}

// PasswordReset handles incoming password reset messages.
// It logs the received message, unmarshals it into a UserPasswordReset struct,
// and then calls the corresponding usecase logic.
func (h *handler) PasswordReset(ctx context.Context, message []byte) error {
	h.logger.Infow(ctx, "Received Password Reset Message", "message", string(message))

	// Parse the incoming JSON message into the UserPasswordReset struct.
	var req domain.UserPasswordReset
	if err := json.Unmarshal(message, &req); err != nil {
		h.logger.Errorw(ctx, "Failed to unmarshal password reset message", "error", err)
		return err
	}

	// Call the business logic (usecase) to handle the password reset.
	if err := h.usecases.User.PasswordReset(ctx, req); err != nil {
		h.logger.Errorw(ctx, "Failed to process password reset usecase", "error", err)
		return err
	}

	h.logger.Infow(ctx, "Successfully processed Password Reset Message", "to", req.To)
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
