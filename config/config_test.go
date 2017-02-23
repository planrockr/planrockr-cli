package config

import (
	"testing"

	"os"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitWithConfigFile(t *testing.T) {
	viper.Reset()
	SetParameters(Params{"config", "./../data/test"})
	err := Init()
	if err != nil {
		t.Error(err)
	}
	c := Get()
	assert.EqualValues(t, "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI3NjJmMjliZjYyZWFmZThiZjBiNmNjNDAwZDIxNjg5NzUzNGVkZjU3Nzk5Yzg2MTRmMzBjMGMwNGViOWE5Y2JjIiwiaXNzIjoiaHR0cDpcL1wvcGxhbnJvY2tyLmNvbSIsImF1ZCI6Imh0dHA6XC9cL3BsYW5yb2Nrci5jb20iLCJpYXQiOjE0ODY1NjcxMDEsIm5iZiI6MTQ4NjU2NzE2MSwiZXhwIjoxNDg3MTcxOTAxLCJ1c2VySWQiOjZ9.lU48XrcS5_EO_wFyikaYdSa7-yrq8JkCYe1m3LTnN71", c.Auth.Token, "Config get wrong value for auth.token")
	assert.EqualValues(t, "$2y$10$cH1lgjajPQIXGH.XxWB2eeA0WRb3Y9MfE77Cx3vKjHxK.hW.sh0a", c.Auth.RefreshToken, "Config get wrong value for auth.refreshtoken")
	assert.EqualValues(t, "http://app.planrockr.dev", c.BaseUrl, "Config get wrong value for baseurl")
	assert.EqualValues(t, 1, c.Auth.Id, "Config get wrong value for auth id")
}

func TestSaveToken(t *testing.T) {
	viper.Reset()
	expectedToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI3NjJmMjliZjYyZWFmZThiZjBiNmNjNDAwZDIxNjg5NzUzNGVkZjU3Nzk5Yzg2MTRmMzBjMGMwNGViOWE5Y2JjIiwiaXNzIjoiaHR0cDpcL1wvcGxhbnJvY2tyLmNvbSIsImF1ZCI6Imh0dHA6XC9cL3BsYW5yb2Nrci5jb20iLCJpYXQiOjE0ODY1NjcxMDEsIm5iZiI6MTQ4NjU2NzE2MSwiZXhwIjoxNDg3MTcxOTAxLCJ1c2VySWQiOjZ9.lU48XrcS5_EO_wFyikaYdSa7-yrq8JkCYe1m3LTnN71"
	expectedRefreshToken := "$2y$10$cH1lgjajPQIXGH.XxWB2eeA0WRb3Y9MfE77Cx3vKjHxK.hW.sh0a"
	expectedId := 1
	SetParameters(Params{"config", "./../data/test"})
	err := Init()
	if err != nil {
		t.Error(err)
	}
	c := Get()
	assert.EqualValues(t, expectedToken, c.Auth.Token, "Config get wrong value for auth.token")
	assert.EqualValues(t, expectedRefreshToken, c.Auth.RefreshToken, "Config get wrong value for auth.refresh_token")
	assert.EqualValues(t, expectedId, c.Auth.Id, "Config get wrong value for auth.id")
	SetString("auth.token", "teste")
	err = Init()
	if err != nil {
		t.Error(err)
	}
	c = Get()
	assert.EqualValues(t, "teste", c.Auth.Token, "Config get wrong value for auth.token")
	SetString("auth.token", expectedToken)
}

func TestInitWithEnviromentVariables(t *testing.T) {
	viper.Reset()
	err := os.Setenv("PLANROCKR_AUTH_TOKEN", "the token")
	if err != nil {
		t.Error(err)
	}
	err = os.Setenv("PLANROCKR_AUTH_REFRESHTOKEN", "the refresh token")
	if err != nil {
		t.Error(err)
	}
	err = os.Setenv("PLANROCKR_BASEURL", "base url")
	if err != nil {
		t.Error(err)
	}
	SetParameters(Params{"config", "./../data/test"})
	err = Init()
	if err != nil {
		t.Error(err)
	}
	c := Get()
	assert.EqualValues(t, "the refresh token", c.Auth.RefreshToken, "Config get wrong value for auth.refresh_token")
	assert.EqualValues(t, "base url", c.BaseUrl, "Config get wrong value for baseurl")
	assert.EqualValues(t, "the token", c.Auth.Token, "Config get wrong value for auth.token")
	os.Clearenv()
}

func TestSaveWithOutConfigFile(t *testing.T) {
	viper.Reset()
	SetParameters(Params{"config_not_found", "/tmp"})
	err := Init() // Find and read the config file
	if err != nil {
		t.Error(err)
	}
	err = SetString("auth.token", "teste")
	if err != nil {
		t.Error(err)
	}
	err = Init() // Find and read the config file
	var config Config
	viper.Unmarshal(&config)
	assert.EqualValues(t, "teste", config.Auth.Token, "Config get wrong value for auth.token")
	os.Remove("/tmp/config_not_found.yaml")
}
