package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Auth       struct {
		Token     string
		RefreshToken     string
	}
}

var config Config

// Init will initilize the Config struct using one file or enviroment variables.
func Init() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/planrockr-cli/")
	viper.AddConfigPath("$HOME/.planrockr-cli/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return err
	}
	//Override configs with ENV variables
	viper.SetEnvPrefix("planrockr")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err = viper.Unmarshal(&config)
	return err
}

// Get will return the config initialized.
func Get() Config {
	return config
}
