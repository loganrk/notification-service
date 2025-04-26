package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/loganrk/message-service/config"
	"github.com/loganrk/message-service/internal/core/port"
	userUsecase "github.com/loganrk/message-service/internal/core/usecase/user"

	aesCipher "github.com/loganrk/message-service/internal/adapters/cipher/aes"
	"github.com/loganrk/message-service/internal/adapters/handler"
	zapLogger "github.com/loganrk/message-service/internal/adapters/logger/zapLogger"
	kafkaMessage "github.com/loganrk/message-service/internal/adapters/message/kafka"
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
	logger, err := initLogger(appConfig.GetLogger())
	if err != nil {
		log.Println("failed to initialize logger:", err)
		return
	}

	// Initialize kafka
	kafkaIns, err := initKafka(appConfig.GetKafka())
	if err != nil {
		logger.Errorw(context.Background(), "failed to initialize kafka", "error", err)
		return
	}

	// Initialize user service
	userService := userUsecase.New(logger)
	services := port.SvrList{User: userService}

	handlerIns := handler.New(logger, services)

	kafkaIns.ConsumeActivation(context.Background(), handlerIns.Activation, handlerIns.ActivationError)
	kafkaIns.ConsumePasswordReset(context.Background(), handlerIns.PasswordReset, handlerIns.PasswordResetError)

}

// initLogger creates a new zap-based logger with the given config.
func initLogger(conf config.Logger) (port.Logger, error) {
	loggerConf := zapLogger.Config{
		Level:          conf.GetLoggerLevel(),
		Encoding:       conf.GetLoggerEncodingMethod(),
		EncodingCaller: conf.GetLoggerEncodingCaller(),
		OutputPath:     conf.GetLoggerPath(),
	}
	return zapLogger.New(loggerConf)
}

// initKafka creates a kafka instance with decrypted brokers.
func initKafka(conf config.Kafka) (port.Messager, error) {
	cipherKey := os.Getenv("CIPHER_CRYPTO_KEY")
	cipher := aesCipher.New(cipherKey)
	var brokers []string

	for _, brokerEnc := range conf.GetBrokers() {
		broker, err := cipher.Decrypt(brokerEnc)
		if err != nil {
			return nil, err
		}
		brokers = append(brokers, broker)
	}

	return kafkaMessage.New(brokers, conf)
}
