package port

import (
	"context"
)

type Hanlder interface {
	ActivationEmail(to, subject string, macros map[string]string) error
	ActivationPhone(to string, macros map[string]string) error

	PasswordResetEmail(to, subject string, macros map[string]string) error
	PasswordResetPhone(to string, macros map[string]string) error

	ActivationError(ctx context.Context, err error)
	PasswordResetError(ctx context.Context, err error)
}

// Cipher defines the interface for encrypting and decrypting strings.
type Cipher interface {
	Encrypt(text string) (string, error)       // Encrypts a plain text string and returns the encrypted version
	Decrypt(cryptoText string) (string, error) // Decrypts an encrypted string and returns the original plain text
	GetKey() string                            // Returns the key used for encryption/decryption
}

// Logger defines the interface for structured and leveled logging.
type Logger interface {
	Debug(ctx context.Context, messages ...any) // Logs debug messages
	Info(ctx context.Context, messages ...any)  // Logs informational messages
	Warn(ctx context.Context, messages ...any)  // Logs warning messages
	Error(ctx context.Context, messages ...any) // Logs error messages
	Fatal(ctx context.Context, messages ...any) // Logs fatal messages and exits the application

	Debugf(ctx context.Context, template string, args ...any) // Logs formatted debug messages
	Infof(ctx context.Context, template string, args ...any)  // Logs formatted informational messages
	Warnf(ctx context.Context, template string, args ...any)  // Logs formatted warning messages
	Errorf(ctx context.Context, template string, args ...any) // Logs formatted error messages
	Fatalf(ctx context.Context, template string, args ...any) // Logs formatted fatal messages and exits the application

	Debugw(ctx context.Context, msg string, keysAndValues ...any) // Logs structured debug messages
	Infow(ctx context.Context, msg string, keysAndValues ...any)  // Logs structured informational messages
	Warnw(ctx context.Context, msg string, keysAndValues ...any)  // Logs structured warning messages
	Errorw(ctx context.Context, msg string, keysAndValues ...any) // Logs structured error messages
	Fatalw(ctx context.Context, msg string, keysAndValues ...any) // Logs structured fatal messages and exits the application

	Sync(ctx context.Context) error // Flushes any buffered log entries
}

type MessageReceiver interface {
	RegisterPasswordResetHandlers(
		passwordResetTopic string,
		phoneHandler func(to string, macros map[string]string) error,
		emailHandler func(to, subject string, macros map[string]string) error,
	) error
	RegisterActivation(
		activationTopic string,
		phoneHandler func(to string, macros map[string]string) error,
		emailHandler func(to, subject string, macros map[string]string) error,

	) error
	ListenActivationHResetTopic(ctx context.Context, errorHandler func(ctx context.Context, err error)) error
	ListenPasswordResetTopic(ctx context.Context, errorHandler func(ctx context.Context, err error)) error
}

type EmailSender interface {
	SendEmail(to, subject, body string) error
}
