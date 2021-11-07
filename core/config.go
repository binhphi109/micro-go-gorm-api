package core

import (
	"github.com/spf13/viper"
)

type Config struct {
	MYSQL_CONNECTION  string
	TOKEN_SIGN_METHOD string
	TOKEN_SIGN_KEY    string
	TOKEN_VERIFY_KEY  string
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		panic(err)
	}

	return &Config{
		MYSQL_CONNECTION:  viper.GetString("MYSQL_CONNECTION"),
		TOKEN_SIGN_METHOD: viper.GetString("TOKEN_SIGN_METHOD"),
		TOKEN_SIGN_KEY:    viper.GetString("TOKEN_SIGN_KEY"),
		TOKEN_VERIFY_KEY:  viper.GetString("TOKEN_VERIFY_KEY"),
	}
}
