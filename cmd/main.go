package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/loganrk/message-service/config"
	"github.com/loganrk/message-service/internal/core/port"
	userUsecase "github.com/loganrk/message-service/internal/core/usecase/user"

	emailSender "github.com/loganrk/message-service/internal/adapters/emailSender/mailjet"
	"github.com/loganrk/message-service/internal/adapters/handler"
	messageReceiver "github.com/loganrk/utils-go/adapters/message/kafka/consumer"

	cipher "github.com/loganrk/utils-go/adapters/cipher/aes"
	logger "github.com/loganrk/utils-go/adapters/logger/zapLogger"
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

	//Initialize Kafka message receiver
	messageReceiverIns, err := initMessageReceiver(appConfig.GetKafka(), appConfig.GetAppName(), handlerIns, cipherIns)
	if err != nil {
		loggerIns.Errorw(context.Background(), "failed to initialize kafka", "error", err)
		return
	}

	go messageReceiverIns.ListenActivationHResetTopic(context.Background(), handlerIns.ActivationError)

	go messageReceiverIns.ListenPasswordResetTopic(context.Background(), handlerIns.PasswordResetError)

	fmt.Println("server start")
	select {}
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
func initMessageReceiver(conf config.Kafka, appName string, handlerIns port.Hanlder, cipherIns port.Cipher) (port.MessageReceiver, error) {
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
	messageReceiverIns := messageReceiver.New(
		brokers,
		strings.Replace(conf.GetConsumerGroupName(), "{{appName}}", appName, 1),
	)

	//Start message receiver consumers for different event types
	err := messageReceiverIns.RegisterActivation(conf.GetActivationTopic(), handlerIns.ActivationPhone, handlerIns.ActivationEmail)
	if err != nil {
		return nil, err
	}
	err = messageReceiverIns.RegisterPasswordResetHandlers(conf.GetPasswordResetTopic(), handlerIns.PasswordResetPhone, handlerIns.PasswordResetEmail)
	if err != nil {
		return nil, err
	}

	return messageReceiverIns, nil

}

// initEmailSender decrypts SMTP credentials and initializes the email sender.
func initEmailSender(conf config.Email, cipherIns port.Cipher) (port.EmailSender, error) {
	// Decrypt host
	apiKey, err := cipherIns.Decrypt(conf.GetMailjetAPIKey())
	if err != nil {
		return nil, err
	}

	// Decrypt password
	apiSecret, err := cipherIns.Decrypt(conf.GetMailjetAPISecret())
	if err != nil {
		return nil, err
	}

	// Return new email sender instance
	return emailSender.New(apiKey, apiSecret, conf.GetMailjetFromEmail(), conf.GetMailjetFromName()), nil
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
