package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

// Message is the structure expected from Kafka.
type message struct {
	Type    string            `json:"type"`    // e.g., "activation_email", "reset_phone"
	To      string            `json:"to"`      // email or phone
	Subject string            `json:"subject"` // only for email
	Macros  map[string]string `json:"macros"`  // template variables
}

// consumer is a Kafka adapter that handles two different consumer groups:
// one for user activation and one for password reset.
type consumer struct {
	activationTopic        string
	activationConsumer     sarama.ConsumerGroup
	activationPhoneHandler func(to string, macros map[string]string) error
	activationEmailHandler func(to, subject string, macros map[string]string) error

	passwordResetTopic        string
	passwordResetConsumer     sarama.ConsumerGroup
	passwordResetPhoneHandler func(to string, macros map[string]string) error
	passwordResetEmailHandler func(to, subject string, macros map[string]string) error

	groupID      string
	brokers      []string
	saramaConfig *sarama.Config
}

// New initializes the consumer with the provided Kafka connection details.
// It creates separate consumer groups for activation and password reset topics.
func New(brokers []string, groupID string) *consumer {

	// Create a common Sarama config
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest // Start from latest messages

	// Return initialized consumer
	return &consumer{
		groupID:      groupID,
		brokers:      brokers,
		saramaConfig: cfg,
	}
}

// RegisterActivationHandlers sets both activation handlers at once
func (c *consumer) RegisterActivation(
	activationTopic string,
	phoneHandler func(to string, macros map[string]string) error,
	emailHandler func(to, subject string, macros map[string]string) error,

) error {
	c.activationTopic = activationTopic
	c.activationPhoneHandler = phoneHandler
	c.activationEmailHandler = emailHandler

	// Create consumer group for activation topic
	activationConsumer, err := sarama.NewConsumerGroup(c.brokers, c.groupID, c.saramaConfig)
	if err != nil {
		return err
	}
	c.activationConsumer = activationConsumer

	return nil

}

// RegisterPasswordResetHandlers sets both password reset handlers at once
func (c *consumer) RegisterPasswordResetHandlers(
	passwordResetTopic string,
	phoneHandler func(to string, macros map[string]string) error,
	emailHandler func(to, subject string, macros map[string]string) error,
) error {
	c.passwordResetTopic = passwordResetTopic
	c.passwordResetPhoneHandler = phoneHandler
	c.passwordResetEmailHandler = emailHandler

	// Create consumer group for password reset topic
	passwordResetConsumer, err := sarama.NewConsumerGroup(c.brokers, c.groupID, c.saramaConfig)
	if err != nil {
		return err
	}
	c.passwordResetConsumer = passwordResetConsumer
	return nil
}

// ConsumeActivation starts consuming activation messages.
func (c *consumer) ListenActivationHResetTopic(ctx context.Context, errorHandler func(context.Context, error)) error {
	return c.consume(ctx, c.activationConsumer, c.activationTopic, c.routeActivation, errorHandler)
}

// ConsumePasswordReset starts consuming password reset messages.
func (c *consumer) ListenPasswordResetTopic(ctx context.Context, errorHandler func(context.Context, error)) error {
	return c.consume(ctx, c.passwordResetConsumer, c.passwordResetTopic, c.routePasswordReset, errorHandler)
}

// Internal router for activation topic messages.
func (c *consumer) routeActivation(ctx context.Context, msgBytes []byte) error {
	var msg message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return err
	}

	switch msg.Type {
	case "verification-phone":
		return c.activationPhoneHandler(msg.To, msg.Macros)
	case "verification-email":
		return c.activationEmailHandler(msg.To, msg.Subject, msg.Macros)
	default:
		return fmt.Errorf("unknown activation type: %s", msg.Type)
	}
}

// Internal router for password reset topic messages.
func (c *consumer) routePasswordReset(ctx context.Context, msgBytes []byte) error {
	var msg message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return err
	}

	switch msg.Type {
	case "password-reset-phone":
		return c.passwordResetPhoneHandler(msg.To, msg.Macros)
	case "password-reset-email":
		return c.passwordResetEmailHandler(msg.To, msg.Subject, msg.Macros)
	default:
		return fmt.Errorf("unknown password reset type: %s", msg.Type)
	}
}

// consume spawns a goroutine to consume from a Kafka topic.
func (c *consumer) consume(
	ctx context.Context,
	consumerGroup sarama.ConsumerGroup,
	topic string,
	messageHandler func(context.Context, []byte) error,
	errorHandler func(context.Context, error),
) error {
	go func() {
		defer consumerGroup.Close()

		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, &consumerHandler{
				messageHandler: messageHandler,
				errorHandler:   errorHandler,
			}); err != nil {
				errorHandler(ctx, err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

// consumerHandler delegates messages to a handler.
type consumerHandler struct {
	messageHandler func(context.Context, []byte) error
	errorHandler   func(context.Context, error)
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.messageHandler(session.Context(), message.Value); err != nil {
			h.errorHandler(session.Context(), err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}
