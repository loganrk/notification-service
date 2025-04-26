package port

import (
	"context"

	"github.com/loganrk/message-service/internal/core/domain"
)

type SvrList struct {
	User UserSvr
}

type UserSvr interface {
	Activation(ctx context.Context, req domain.UserActivation) error
	PasswordReset(ctx context.Context, req domain.UserPasswordReset) error
}
