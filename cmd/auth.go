// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/planrockr/planrockr-cli/config"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	user     string
	password string
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth commands",
	Long:  "You need to use the e-mail and password that you use to log in http://planrockr.com",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := doLogin(user, password)
		if err != nil {
			fmt.Println(err)
		}
		//@todo: save on ~/.planrockr-cli.yaml
		fmt.Printf("%s", token)
	},
}

func doLogin(user string, password string) (token string, err error) {
	body := strings.NewReader("parameters%5Blogin%5D=" + user + "&parameters%5Bpassword%5D=" + password)
	req, err := http.NewRequest("POST", "https://app.planrockr.com/rpc/v1/authentication/login", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "planrockr-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err != nil || resp.Status == "404 Not Found" {
		return "", errors.New("Invalid credentials")
	}
	buf, _ := ioutil.ReadAll(resp.Body)

	err = config.Init()
	if err != nil {
		return "", errors.New("Error reading config file")
	}

	type AuthData struct {
		Token         string
		Refresh_Token string
	}
	var authData AuthData
	err = json.Unmarshal(buf, &authData)
	if err != nil {
		return "", errors.New("Error parsing authorization data")
	}
	err = config.Set("auth.token", authData.Token)
	err = config.Set("auth.refreshtoken", authData.Refresh_Token)
	if err != nil {
		return "", errors.New("Error writing config file")
	}

	return "Authorized\n", nil
}

func init() {
	RootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(&user, "user", "u", "", "User e-mail")
	authCmd.Flags().StringVarP(&password, "password", "p", "", "User password")

}
