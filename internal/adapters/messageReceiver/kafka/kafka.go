package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

// kafkaMessager is a Kafka adapter that handles two different consumer groups:
// one for user activation and one for password reset.
type kafkaMessager struct {
	activationConsumer    sarama.ConsumerGroup
	passwordResetConsumer sarama.ConsumerGroup

	activationTopic    string
	passwordResetTopic string
	groupID            string
	brokers            []string
	saramaConfig       *sarama.Config
}

// New initializes the kafkaMessager with the provided Kafka connection details.
// It creates separate consumer groups for activation and password reset topics.
func New(brokers []string, groupID string, activationTopic string, passwordResetTopic string) (*kafkaMessager, error) {

	// Create a common Sarama config
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest // Start from latest messages

	// Create consumer group for activation topic
	activationConsumer, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}

	// Create consumer group for password reset topic
	passwordResetConsumer, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}

	// Return initialized kafkaMessager
	return &kafkaMessager{
		activationConsumer:    activationConsumer,
		passwordResetConsumer: passwordResetConsumer,
		activationTopic:       activationTopic,
		passwordResetTopic:    passwordResetTopic,
		groupID:               groupID,
		brokers:               brokers,
		saramaConfig:          cfg,
	}, nil
}

// ConsumeActivation starts consuming activation messages from the Kafka topic.
func (k *kafkaMessager) ConsumeActivation(ctx context.Context, messageHandler func(context.Context, []byte) error, errorHandler func(context.Context, error)) error {
	return k.consume(ctx, k.activationConsumer, k.activationTopic, messageHandler, errorHandler)
}

// ConsumePasswordReset starts consuming password reset messages from the Kafka topic.
func (k *kafkaMessager) ConsumePasswordReset(ctx context.Context, messageHandler func(context.Context, []byte) error, errorHandler func(context.Context, error)) error {
	return k.consume(ctx, k.passwordResetConsumer, k.passwordResetTopic, messageHandler, errorHandler)
}

// consume starts a goroutine that listens for messages on the specified topic using the given consumer group.
// It continuously invokes the provided message handler and handles errors using the error handler.
func (k *kafkaMessager) consume(ctx context.Context, consumerGroup sarama.ConsumerGroup, topic string, messageHandler func(context.Context, []byte) error, errorHandler func(context.Context, error)) error {
	go func() {
		defer consumerGroup.Close()

		for {
			// Consume messages from the topic and delegate to consumerHandler
			if err := consumerGroup.Consume(ctx, []string{topic}, &consumerHandler{
				messageHandler: messageHandler,
				errorHandler:   errorHandler,
			}); err != nil {
				errorHandler(ctx, err)
			}

			// Exit if the context is cancelled
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

// consumerHandler implements sarama.ConsumerGroupHandler
// It delegates message processing and error handling to external functions.
type consumerHandler struct {
	messageHandler func(context.Context, []byte) error
	errorHandler   func(context.Context, error)
}

// Setup is called before consuming starts, can be used to initialize resources.
// No setup needed in this case.
func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup is called after consuming stops, can be used to clean up resources.
// No cleanup needed in this case.
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim handles the actual message consumption logic.
// It reads messages from the claim and delegates them to the message handler.
func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// Process message
		if err := h.messageHandler(session.Context(), message.Value); err != nil {
			h.errorHandler(session.Context(), err)
		}
		// Mark message as processed
		session.MarkMessage(message, "")
	}
	return nil
}
