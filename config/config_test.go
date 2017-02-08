package config

import (
	"testing"

	"os"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitWithConfigFile(t *testing.T) {
	viper.AddConfigPath("./../data/test/")
	err := Init()
	if err != nil {
		t.Error(err)
	}
	c := Get()
	assert.EqualValues(t, "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI3NjJmMjliZjYyZWFmZThiZjBiNmNjNDAwZDIxNjg5NzUzNGVkZjU3Nzk5Yzg2MTRmMzBjMGMwNGViOWE5Y2JjIiwiaXNzIjoiaHR0cDpcL1wvcGxhbnJvY2tyLmNvbSIsImF1ZCI6Imh0dHA6XC9cL3BsYW5yb2Nrci5jb20iLCJpYXQiOjE0ODY1NjcxMDEsIm5iZiI6MTQ4NjU2NzE2MSwiZXhwIjoxNDg3MTcxOTAxLCJ1c2VySWQiOjZ9.lU48XrcS5_EO_wFyikaYdSa7-yrq8JkCYe1m3LTnN71", c.Auth.Token, "Config get wrong value for auth.token")
	assert.EqualValues(t, "$2y$10$cH1lgjajPQIXGH.XxWB2eeA0WRb3Y9MfE77Cx3vKjHxK.hW.sh0a", c.Auth.RefreshToken, "Config get wrong value for auth.refresh_token")
}

func TestInitWithEnviromentVariables(t *testing.T) {
	err := os.Setenv("PLANROCKR_AUTH_TOKEN", "the token")
	if err != nil {
		t.Error(err)
	}
	err = os.Setenv("PLANROCKR_AUTH_REFRESHTOKEN", "the refresh token")
	if err != nil {
		t.Error(err)
	}
	viper.AddConfigPath("./../data/")
	err = Init()
	if err != nil {
		t.Error(err)
	}
	c := Get()
	assert.EqualValues(t, "the token", c.Auth.Token, "Config get wrong value for auth.token")
	assert.EqualValues(t, "the refresh token", c.Auth.RefreshToken, "Config get wrong value for auth.refresh_token")

}
