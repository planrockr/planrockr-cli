package config

import (
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"errors"
	"io"
	"os"
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
	// _, err = configFileWriter()
	// if err != nil {
		// panic(err)
	// }

	return err
}

// Get will return the config initialized.
func Get() Config {
	return config
}

func Set(key string, value string) error {
	viper.Set(key, value)
	return writeConfig()
}

func writeConfig() error {
	f, err := configFileWriter()
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		return errors.New("unable to encode configuration to YAML format")
	}

	_, err = f.Write(b)
	if err != nil {
		return errors.New("unable to write configuration")
	}

	return nil
}

func configFileWriter() (io.WriteCloser, error) {
	cfgFile := viper.ConfigFileUsed()
	f, err := os.Create(cfgFile)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(cfgFile, 0600); err != nil {
		return nil, err
	}

	return f, nil
}
