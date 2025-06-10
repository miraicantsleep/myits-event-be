package config

import (
	"github.com/spf13/viper"
)

type EmailConfig struct {
	Host         string `mapstructure:"SMTP_HOST"`
	Port         int    `mapstructure:"SMTP_PORT"`
	SenderName   string `mapstructure:"SMTP_SENDER_NAME"`
	AuthUsername string `mapstructure:"SMTP_AUTH_EMAIL"` // still your AKIA… string
	AuthPassword string `mapstructure:"SMTP_AUTH_PASSWORD"`
	ApiBaseUrl   string `mapstructure:"API_BASE_URL"` // New field
}

func NewEmailConfig() (*EmailConfig, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	var config EmailConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
