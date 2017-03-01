package config

import (
	"strings"

	"errors"
	"io"
	"os"

	"runtime"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BaseUrl string
	Auth    struct {
		Token        string
		RefreshToken string
		Id           int
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
		defaultParams.ConfigPath = getUserHomeDir()
	}
	var _, err = os.Stat(defaultParams.ConfigPath + "/" + defaultParams.ConfigName + ".yaml")
	if os.IsNotExist(err) {
		var file, err = os.Create(defaultParams.ConfigPath + "/" + defaultParams.ConfigName + ".yaml")
		if err != nil {
			return err
		}
		file.Close()
		SetString("baseurl", "https://app.planrockr.com")
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

func getUserHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}
	return os.Getenv(env)
}

func SetParameters(params Params) {
	defaultParams = params
}

// Get will return the config initialized.
func Get() Config {
	return config
}

func SetString(key string, value string) error {
	viper.Set(key, value)
	return writeConfig()
}

func SetInt(key string, value int) error {
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
