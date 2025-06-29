package port

import "context"

type SvrList struct {
	User UserSvr
}

type UserSvr interface {
	ActivationEmail(ctx context.Context, to, subject string, macros map[string]string) error
	ActivationPhone(ctx context.Context, to string, macros map[string]string) error

	PasswordResetEmail(ctx context.Context, to, subject string, macros map[string]string) error
	PasswordResetPhone(ctx context.Context, to string, macros map[string]string) error
}
