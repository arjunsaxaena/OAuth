package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	MessageCentral MessageCentralConfig `mapstructure:"message_central"`
}

type DatabaseConfig struct {
	DBURL string `mapstructure:"db_url"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type MessageCentralConfig struct {
	AuthURL    string `mapstructure:"auth_url"`
	ValidateURL string `mapstructure:"validate_url"`
	SendOTPURL string `mapstructure:"send_otp_url"`
	CID        string `mapstructure:"cid"`
	Key        string `mapstructure:"key"`
}

func LoadConfig() (*Config, error) {
    paths := []string{".", "./..", "./../.."}
    var envErr error
    for _, p := range paths {
        envPath := filepath.Join(p, ".env")
        if err := godotenv.Overload(envPath); err == nil {
            envErr = nil
            break
        } else {
            envErr = err
        }
    }
    if envErr != nil {
        log.Printf(".env not loaded from search paths: %v", envErr)
    }

    viper.AutomaticEnv()
    bindEnvVars()

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("unable to decode config into struct: %w", err)
    }

    if config.Database.DBURL == "" {
        config.Database.DBURL = viper.GetString("DB_URL")
    }

    return &config, nil
}

func bindEnvVars() {
    viper.BindEnv("database.db_url", "DB_URL")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("message_central.auth_url", "MESSAGE_CENTRAL_AUTH_URL")
	viper.BindEnv("message_central.validate_url", "MESSAGE_CENTRAL_VALIDATE_URL")
	viper.BindEnv("message_central.send_otp_url", "MESSAGE_CENTRAL_SEND_OTP_URL")
	viper.BindEnv("message_central.cid", "MESSAGE_CENTRAL_CID")
	viper.BindEnv("message_central.key", "MESSAGE_CENTRAL_KEY")
}

func (c *Config) GetConfigMap() map[string]interface{} {
	return map[string]interface{}{
		"database": map[string]interface{}{
			"db_url": c.Database.DBURL,
		},
		"jwt": map[string]interface{}{
			"secret": c.JWT.Secret,
		},
		"message_central": map[string]interface{}{
			"auth_url":     c.MessageCentral.AuthURL,
			"validate_url": c.MessageCentral.ValidateURL,
			"send_otp_url": c.MessageCentral.SendOTPURL,
			"cid":          c.MessageCentral.CID,
			"key":          c.MessageCentral.Key,
		},
	}
}

func (c *Config) GetEnvMap() map[string]string {
	return map[string]string{
		"DB_URL":                        c.Database.DBURL,
		"JWT_SECRET":                    c.JWT.Secret,
		"MESSAGE_CENTRAL_AUTH_URL":      c.MessageCentral.AuthURL,
		"MESSAGE_CENTRAL_VALIDATE_URL":  c.MessageCentral.ValidateURL,
		"MESSAGE_CENTRAL_SEND_OTP_URL":  c.MessageCentral.SendOTPURL,
		"MESSAGE_CENTRAL_CID":           c.MessageCentral.CID,
		"MESSAGE_CENTRAL_KEY":           c.MessageCentral.Key,
	}
}

