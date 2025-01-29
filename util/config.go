package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct{
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	// RefreshTokenDuration string `mapstructure:"REFRESH_TOKEN_DURATION"`
	// AccessTokenSecret string `mapstructure:"ACCESS_TOKEN_SECRET"`
	// RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
	// AccessTokenIssuer string `mapstructure:"ACCESS_TOKEN_ISSUER"`
	// RefreshTokenIssuer string `mapstructure:"REFRESH_TOKEN_ISSUER"`
	// DBDriver string `mapstructure:"DB_DRIVER"`
	// DBSource string `mapstructure:"DB_SOURCE"`
	// ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	// TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	// AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	// RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	// AccessTokenSecret string `mapstructure:"ACCESS_TOKEN_SECRET"`
	// RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
	// AccessTokenIssuer string `mapstructure:"ACCESS_TOKEN_ISSUER"`
	// RefreshTokenIssuer string `mapstructure:"REFRESH_TOKEN_ISSUER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}