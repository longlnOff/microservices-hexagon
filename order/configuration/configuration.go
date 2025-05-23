package configuration

import (
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	GRPCServer 		GRPCServerConfiguration
	HTTPServer   	HTTPServerConfiguration
	Database 		DatabaseConfiguration
	Cache    		CacheConfiguration
	Payment			PaymentConfiguration
	Mail     		MailConfiguration
	Auth     		AuthConfiguration
}

type PaymentConfiguration struct {
	PAYMENT_SERVICE_ADDRESS string `mapstructure:"PAYMENT_SERVICE_ADDRESS"`
}

type CacheConfiguration struct {
	CACHE_ADDRESS  string `mapstructure:"CACHE_ADDRESS"`
	CACHE_PASSWORD string `mapstructure:"CACHE_PASSWORD"`
	CACHE_DATABASE int    `mapstructure:"CACHE_DATABASE"`
	CACHE_ENABLED  bool   `mapstructure:"CACHE_ENABLED"`
}

type AuthConfiguration struct {
	Basic BasicAuthencation
	Token TokenAuthentication
}

type TokenAuthentication struct {
	AUTH_TOKEN_SECRET string        `mapstructure:"AUTH_TOKEN_SECRET"`
	AUTH_TOKEN_EXP    time.Duration `mapstructure:"AUTH_TOKEN_EXP"`
	AUTH_TOKEN_ISS    string        `mapstructure:"AUTH_TOKEN_ISS"`
}

type BasicAuthencation struct {
	AUTH_BASIC_USER     string `mapstructure:"AUTH_BASIC_USER"`
	AUTH_BASIC_PASSWORD string `mapstructure:"AUTH_BASIC_PASSWORD"`
}

type MailConfiguration struct {
	EXP        time.Duration `mapstructure:"MAIL_EXP"`
	FROM_EMAIL string        `mapstructure:"FROM_EMAIL"`
	SendGrid   SendGridConfiguration
	MailTrap   MailTrapConfiguration
}

type SendGridConfiguration struct {
	API_KEY string `mapstructure:"SEND_GRID_API_KEY"`
}

type MailTrapConfiguration struct {
	API_KEY string `mapstructure:"MAIL_TRAP_API_KEY"`
}

type HTTPServerConfiguration struct {
	SERVER_ADDRESS   string `mapstructure:"SERVER_ADDRESS"`
	SERVER_PORT      string `mapstructure:"SERVER_PORT"`
	ENVIRONMENT      string `mapstructure:"ENVIRONMENT"`
	VERSION          string `mapstructure:"VERSION"`
	EXTERNAL_ADDRESS string `mapstructure:"EXTERNAL_ADDRESS"`
	EXTERNAL_PORT    string `mapstructure:"EXTERNAL_PORT"`
	FRONTEND_URL     string `mapstructure:"FRONTEND_URL"`
}

type GRPCServerConfiguration struct {
	GRPC_SERVER_ADDRESS   string `mapstructure:"GRPC_SERVER_ADDRESS"`
	GRPC_SERVER_PORT      int	 `mapstructure:"GRPC_SERVER_PORT"`
	GRPC_ENVIRONMENT      string `mapstructure:"GRPC_ENVIRONMENT"`
	GRPC_VERSION          string `mapstructure:"GRPC_VERSION"`
	GRPC_EXTERNAL_ADDRESS string `mapstructure:"GRPC_EXTERNAL_ADDRESS"`
	GRPC_EXTERNAL_PORT    string `mapstructure:"GRPC_EXTERNAL_PORT"`
}

type DatabaseConfiguration struct {
	ENGINE            string        `mapstructure:"DB_ENGINE"`
	HOST              string        `mapstructure:"DB_HOST"`
	PORT              string        `mapstructure:"DB_PORT"`
	USER              string        `mapstructure:"DB_USER"`
	PASSWORD          string        `mapstructure:"DB_PASSWORD"`
	DB_NAME           string        `mapstructure:"DB_NAME"`
	DB_MAX_OPEN_CONNS int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DB_MAX_IDLE_CONNS int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DB_MAX_IDLE_TIME  time.Duration `mapstructure:"DB_MAX_IDLE_TIME"`
	DB_SSL_MODE       string        `mapstructure:"DB_SSL_MODE"`
	MIGRATE_UP        bool          `mapstructure:"MIGRATE_UP"`
	MIGRATE_VERSION	  int           `mapstructure:"MIGRATE_VERSION"`
}

func LoadConfig(path string, useAbsPath bool) (cfg Configuration, err error) {
	if useAbsPath {
		viper.AddConfigPath(path)
		viper.SetConfigType("env")
		viper.SetConfigFile(".env")
	}

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return Configuration{}, err
	}

	basicAuth := BasicAuthencation{
		AUTH_BASIC_USER:     viper.GetString("AUTH_BASIC_USER"),
		AUTH_BASIC_PASSWORD: viper.GetString("AUTH_BASIC_PASSWORD"),
	}

	tokenAuth := TokenAuthentication{
		AUTH_TOKEN_SECRET: viper.GetString(("AUTH_TOKEN_SECRET")),
		AUTH_TOKEN_EXP:    viper.GetDuration("AUTH_TOKEN_EXP"),
		AUTH_TOKEN_ISS:    viper.GetString("AUTH_TOKEN_ISS"),
	}

	mail_cfg := MailConfiguration{
		EXP:        viper.GetDuration("MAIL_EXP"),
		FROM_EMAIL: viper.GetString("FROM_EMAIL"),
		SendGrid: SendGridConfiguration{
			API_KEY: viper.GetString("SENDGRID_API_KEY"),
		},
		MailTrap: MailTrapConfiguration{
			API_KEY: viper.GetString("MAILTRAP_API_KEY"),
		},
	}

	http_server_cfg := HTTPServerConfiguration{
		SERVER_ADDRESS:   viper.GetString("SERVER_ADDRESS"),
		SERVER_PORT:      viper.GetString("SERVER_PORT"),
		ENVIRONMENT:      viper.GetString("ENVIRONMENT"),
		VERSION:          viper.GetString("VERSION"),
		EXTERNAL_ADDRESS: viper.GetString("EXTERNAL_ADDRESS"),
		EXTERNAL_PORT:    viper.GetString("EXTERNAL_PORT"),
		FRONTEND_URL:     viper.GetString("FRONTEND_URL"),
	}

	database_cfg := DatabaseConfiguration{
		ENGINE:            viper.GetString("DB_ENGINE"),
		HOST:              viper.GetString("DB_HOST"),
		PORT:              viper.GetString("DB_PORT"),
		USER:              viper.GetString("DB_USER"),
		PASSWORD:          viper.GetString("DB_PASSWORD"),
		DB_NAME:           viper.GetString("DB_NAME"),
		DB_MAX_OPEN_CONNS: viper.GetInt("DB_MAX_OPEN_CONNS"),
		DB_MAX_IDLE_CONNS: viper.GetInt("DB_MAX_IDLE_CONNS"),
		DB_MAX_IDLE_TIME:  viper.GetDuration("DB_MAX_IDLE_TIME"),
		DB_SSL_MODE:       viper.GetString("DB_SSL_MODE"),
		MIGRATE_UP:        viper.GetBool("MIGRATE_UP"),
		MIGRATE_VERSION:   viper.GetInt("MIGRATE_VERSION"),
	}

	cache_cfg := CacheConfiguration{
		CACHE_ADDRESS:  viper.GetString("CACHE_ADDRESS"),
		CACHE_PASSWORD: viper.GetString("CACHE_PASSWORD"),
		CACHE_DATABASE: viper.GetInt("CACHE_DATABASE"),
		CACHE_ENABLED:  viper.GetBool("CACHE_ENABLED"),
	}

	grpc_server_cfg := GRPCServerConfiguration{
		GRPC_SERVER_ADDRESS:   viper.GetString("GRPC_SERVER_ADDRESS"),
		GRPC_SERVER_PORT:      viper.GetInt("GRPC_SERVER_PORT"),
		GRPC_ENVIRONMENT:      viper.GetString("GRPC_ENVIRONMENT"),
		GRPC_VERSION:          viper.GetString("GRPC_VERSION"),
		GRPC_EXTERNAL_ADDRESS: viper.GetString("GRPC_EXTERNAL_ADDRESS"),
		GRPC_EXTERNAL_PORT:    viper.GetString("GRPC_EXTERNAL_PORT"),
	}

	payment_service_cfg := PaymentConfiguration{
		PAYMENT_SERVICE_ADDRESS: viper.GetString("PAYMENT_SERVICE_ADDRESS"),
	}

	return Configuration{
		GRPCServer:   grpc_server_cfg,
		HTTPServer:   http_server_cfg,
		Database: database_cfg,
		Cache:    cache_cfg,
		Payment:  payment_service_cfg,
		Mail:     mail_cfg,
		Auth:     AuthConfiguration{Basic: basicAuth, Token: tokenAuth},
	}, nil
}
