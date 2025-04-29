package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/loganrk/message-service/config"
	"github.com/loganrk/message-service/internal/core/port"
	userUsecase "github.com/loganrk/message-service/internal/core/usecase/user"

	cipher "github.com/loganrk/message-service/internal/adapters/cipher/aes"
	emailSender "github.com/loganrk/message-service/internal/adapters/emailSender/smtp"
	"github.com/loganrk/message-service/internal/adapters/handler"
	logger "github.com/loganrk/message-service/internal/adapters/logger/zapLogger"
	messageReceiver "github.com/loganrk/message-service/internal/adapters/messageReceiver/kafka"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Load config values from file
	configPath := os.Getenv("CONFIG_FILE_PATH")
	configName := os.Getenv("CONFIG_FILE_NAME")
	configType := os.Getenv("CONFIG_FILE_TYPE")

	// Initialize app config
	appConfig, err := config.StartConfig(configPath, config.File{
		Name: configName,
		Ext:  configType,
	})
	if err != nil {
		log.Println("failed to load config:", err)
		return
	}

	// Initialize cipher for encrypted data decryption
	cipherIns := initCipher()

	// Initialize zap-based logger
	loggerIns, err := initLogger(appConfig.GetLogger())
	if err != nil {
		log.Println("failed to initialize logger:", err)
		return
	}

	// Initialize Kafka message receiver
	messageReceiverIns, err := initMessageReceiver(appConfig.GetKafka(), cipherIns)
	if err != nil {
		loggerIns.Errorw(context.Background(), "failed to initialize kafka", "error", err)
		return
	}

	// Initialize SMTP email sender
	emailSenderIns, err := initEmailSender(appConfig.GetEmail(), cipherIns)
	if err != nil {
		loggerIns.Errorw(context.Background(), "failed to initialize email sender", "error", err)
		return
	}

	// Initialize user usecase/service with logger, email sender, and user config
	userServiceIns, err := initUserService(loggerIns, emailSenderIns, appConfig.GetUser())
	if err != nil {
		loggerIns.Errorw(context.Background(), "failed to initialize user usecase", "error", err)
		return
	}

	// Register service(s) to handler
	services := port.SvrList{User: userServiceIns}
	handlerIns := initHandler(loggerIns, services)

	// Start message receiver consumers for different event types
	messageReceiverIns.ConsumeActivation(context.Background(), handlerIns.Activation, handlerIns.ActivationError)
	messageReceiverIns.ConsumePasswordReset(context.Background(), handlerIns.PasswordReset, handlerIns.PasswordResetError)
}

// initCipher initializes the AES cipher using the secret key from environment variable.
func initCipher() port.Cipher {
	cipherKey := os.Getenv("CIPHER_CRYPTO_KEY")
	return cipher.New(cipherKey)
}

// initLogger initializes the zap logger with the provided configuration.
func initLogger(conf config.Logger) (port.Logger, error) {
	loggerConf := logger.Config{
		Level:          conf.GetLoggerLevel(),
		Encoding:       conf.GetLoggerEncodingMethod(),
		EncodingCaller: conf.GetLoggerEncodingCaller(),
		OutputPath:     conf.GetLoggerPath(),
	}
	return logger.New(loggerConf)
}

// initMessageReceiver decrypts the Kafka broker URLs and returns a Kafka receiver instance.
func initMessageReceiver(conf config.Kafka, cipherIns port.Cipher) (port.MessageReceiver, error) {
	var brokers []string

	// Decrypt each broker address
	for _, brokerEnc := range conf.GetBrokers() {
		broker, err := cipherIns.Decrypt(brokerEnc)
		if err != nil {
			return nil, err
		}
		brokers = append(brokers, broker)
	}

	// Pass individual config values to the Kafka adapter
	return messageReceiver.New(
		brokers,
		conf.GetConsumerGroupName(),
		conf.GetActivationTopic(),
		conf.GetPasswordResetTopic(),
	)
}

// initEmailSender decrypts SMTP credentials and initializes the email sender.
func initEmailSender(conf config.Email, cipherIns port.Cipher) (port.EmailSender, error) {
	// Decrypt host
	host, err := cipherIns.Decrypt(conf.GetSMTPHost())
	if err != nil {
		return nil, err
	}

	// Decrypt password
	password, err := cipherIns.Decrypt(conf.GetSMTPPassword())
	if err != nil {
		return nil, err
	}

	// Return new email sender instance
	return emailSender.New(conf.GetSMTPFrom(), host, password, conf.GetSMTPPort()), nil
}

// initHandler initializes the message handler with logger and available services.
func initHandler(logger port.Logger, services port.SvrList) port.Hanlder {
	return handler.New(logger, services)
}

// initUserService creates a new instance of the user service/usecase.
func initUserService(logger port.Logger, emailSender port.EmailSender, conf config.User) (port.UserSvr, error) {

	// Create and return the user service
	return userUsecase.New(logger, emailSender, conf)
}
