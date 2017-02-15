package config

import (
	"strings"

	"errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	Auth struct {
		Token        string
		RefreshToken string
	}
}

type Params struct {
	ConfigName string
	ConfigPath string
}

var config Config

var defaultParams = Params{".planrockr-cli", ""}

// Init will initilize the Config struct using one file or enviroment variables.
func Init() error {
	if defaultParams.ConfigPath == "" {
		defaultParams.ConfigPath = os.Getenv("HOME")
	}
	var _, err = os.Stat(defaultParams.ConfigPath + "/" + defaultParams.ConfigName + ".yaml")
	if os.IsNotExist(err) {
		var file, err = os.Create(defaultParams.ConfigPath + "/" + defaultParams.ConfigName + ".yaml")
		if err != nil {
			return err
		}
		file.Close()
	}

	viper.SetConfigName(defaultParams.ConfigName)
	viper.AddConfigPath(defaultParams.ConfigPath)
	err = viper.ReadInConfig() // Find and read the config file
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

func SetParameters(params Params) {
	defaultParams = params
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
