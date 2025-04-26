package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type File struct {
	Name string
	Ext  string
}
type App interface {
	GetAppName() string
	GetLogger() Logger
	GetActivationTemplatePath() string
	GetPasswordResetTemplatePath() string
	GetKafka() Kafka
}

func StartConfig(path string, file File) (App, error) {
	var appConfig app

	var viperIns = viper.New()

	viperIns.AddConfigPath(path)
	viperIns.SetConfigName(file.Name)
	viperIns.AddConfigPath(".")
	viperIns.SetConfigType(file.Ext)

	if err := viperIns.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err := viperIns.Unmarshal(&appConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return appConfig, nil
}

func (a app) GetAppName() string {
	return a.Application.Name
}

func (a app) GetLogger() Logger {
	return a.Logger
}

func (a app) GetActivationTemplatePath() string {
	return a.Activation.TemplatePath
}

func (a app) GetPasswordResetTemplatePath() string {
	return a.PasswordReset.TemplatePath
}

func (a app) GetKafka() Kafka {
	return a.Kafka
}
