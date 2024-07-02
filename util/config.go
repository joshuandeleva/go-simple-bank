package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	TokenSecret   string        `mapstructure:"TOKEN_SECRET_KEY"`
	TokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// laod config reads configuration from a config file if it exists

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // or viper.SetConfigType("YAML") // xml //json
	viper.AutomaticEnv()       // read in environment variables that match

	err = viper.ReadInConfig()

	if err != nil {
		return
	}
	err = viper.Unmarshal(&config) // unmarshal config -> converts this to a go struct

	return

}
