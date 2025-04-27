package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/loganrk/message-service/config"
)

type kafkaMessager struct {
	activationConsumer    sarama.ConsumerGroup
	passwordResetConsumer sarama.ConsumerGroup

	activationTopic    string
	passwordResetTopic string
	groupID            string
	brokers            []string
	saramaConfig       *sarama.Config
}

// New initializes kafkaMessager with two consumer groups
func New(brokers []string, conf config.Kafka) (*kafkaMessager, error) {
	groupID := conf.GetConsumerGroupName()

	// Create a common Sarama config
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Create two consumer groups
	activationConsumer, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}

	passwordResetConsumer, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}

	return &kafkaMessager{
		activationConsumer:    activationConsumer,
		passwordResetConsumer: passwordResetConsumer,
		activationTopic:       conf.GetActivationTopic(),
		passwordResetTopic:    conf.GetPasswordResetTopic(),
		groupID:               groupID,
		brokers:               brokers,
		saramaConfig:          cfg,
	}, nil
}

func (k *kafkaMessager) ConsumeActivation(ctx context.Context, messageHandler func(context.Context, []byte) error, errorHandler func(context.Context, error)) error {
	return k.consume(ctx, k.activationConsumer, k.activationTopic, messageHandler, errorHandler)
}

func (k *kafkaMessager) ConsumePasswordReset(ctx context.Context, messageHandler func(context.Context, []byte) error, errorHandler func(context.Context, error)) error {
	return k.consume(ctx, k.passwordResetConsumer, k.passwordResetTopic, messageHandler, errorHandler)
}

func (k *kafkaMessager) consume(ctx context.Context, consumerGroup sarama.ConsumerGroup, topic string, messageHandler func(context.Context, []byte) error, errorHandler func(context.Context, error)) error {
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

// internal consumerHandler
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
