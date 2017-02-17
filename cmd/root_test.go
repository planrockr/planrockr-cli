package cmd

import (
	"testing"

	"github.com/planrockr/planrockr-cli/config"
	"github.com/stretchr/testify/assert"
)

func TestWithToken(t *testing.T) {
	_, err := GetToken()
	assert.EqualValues(t, nil, err, "Error on existing token")
	// assert.EqualValues(t, "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI3NjJmMjliZjYyZWFmZThiZjBiNmNjNDAwZDIxNjg5NzUzNGVkZjU3Nzk5Yzg2MTRmMzBjMGMwNGViOWE5Y2JjIiwiaXNzIjoiaHR0cDpcL1wvcGxhbnJvY2tyLmNvbSIsImF1ZCI6Imh0dHA6XC9cL3BsYW5yb2Nrci5jb20iLCJpYXQiOjE0ODY1NjcxMDEsIm5iZiI6MTQ4NjU2NzE2MSwiZXhwIjoxNDg3MTcxOTAxLCJ1c2VySWQiOjZ9.lU48XrcS5_EO_wFyikaYdSa7-yrq8JkCYe1m3LTnN71", token, "Error on existing token")
}

func TestWithoutToken(t *testing.T) {
	config.SetParameters(config.Params{ConfigName: "config_not_found", ConfigPath: "/tmp"})
	//     _, err := GetToken()
	//     assert.EqualValues(t, "Missing token", err.Error(), "Missing token")
}
