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

	// Read config file path, name, and type from environment variables
	configPath := os.Getenv("CONFIG_FILE_PATH")
	configName := os.Getenv("CONFIG_FILE_NAME")
	configType := os.Getenv("CONFIG_FILE_TYPE")

	// Initialize application configuration
	appConfig, err := config.StartConfig(configPath, config.File{
		Name: configName,
		Ext:  configType,
	})
	if err != nil {
		log.Println("failed to load config:", err)
		return
	}

	// Initialize logger
	loggerIns, err := initLogger(appConfig.GetLogger())
	if err != nil {
		log.Println("failed to initialize logger:", err)
		return
	}

	messageReceiverIns, err := initMessageReceiver(appConfig.GetKafka())
	if err != nil {
		loggerIns.Errorw(context.Background(), "failed to initialize kafka", "error", err)
		return
	}

	emailSenderIns, err := initEmailSender(appConfig.GetEmail())
	if err != nil {
		loggerIns.Errorw(context.Background(), "failed to initialize kafka", "error", err)
		return
	}

	// Initialize user service
	userService, err := userUsecase.New(loggerIns, emailSenderIns, appConfig.GetUser())
	if err != nil {
		loggerIns.Errorw(context.Background(), "failed to initialize user usecase", "error", err)
		return
	}

	services := port.SvrList{User: userService}
	handlerIns := handler.New(loggerIns, services)

	messageReceiverIns.ConsumeActivation(context.Background(), handlerIns.Activation, handlerIns.ActivationError)
	messageReceiverIns.ConsumePasswordReset(context.Background(), handlerIns.PasswordReset, handlerIns.PasswordResetError)

}

// initLogger creates a new zap-based logger with the given config.
func initLogger(conf config.Logger) (port.Logger, error) {
	loggerConf := logger.Config{
		Level:          conf.GetLoggerLevel(),
		Encoding:       conf.GetLoggerEncodingMethod(),
		EncodingCaller: conf.GetLoggerEncodingCaller(),
		OutputPath:     conf.GetLoggerPath(),
	}
	return logger.New(loggerConf)
}

// initMessageReceiver creates a kafka instance with decrypted brokers.
func initMessageReceiver(conf config.Kafka) (port.MessageReceiver, error) {
	cipherKey := os.Getenv("CIPHER_CRYPTO_KEY")
	cipherIns := cipher.New(cipherKey)
	var brokers []string

	for _, brokerEnc := range conf.GetBrokers() {
		broker, err := cipherIns.Decrypt(brokerEnc)
		if err != nil {
			return nil, err
		}
		brokers = append(brokers, broker)
	}

	return messageReceiver.New(brokers, conf)
}
func initEmailSender(conf config.Email) (port.EmailSender, error) {
	cipherKey := os.Getenv("CIPHER_CRYPTO_KEY")
	cipherIns := cipher.New(cipherKey)

	// Decrypt Host
	host, err := cipherIns.Decrypt(conf.GetSMTPHost())
	if err != nil {
		return nil, err
	}

	// Decrypt Password
	password, err := cipherIns.Decrypt(conf.GetSMTPPassword())
	if err != nil {
		return nil, err
	}

	return emailSender.New(conf.GetSMTPFrom(), host, password, conf.GetSMTPPort()), nil
}
