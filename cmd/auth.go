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
	authUser     string
	authPassword string
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth commands",
	Long:  "You need to use the e-mail and password that you use to log in http://planrockr.com",
	Run: func(cmd *cobra.Command, args []string) {
		err := doLogin(authUser, authPassword)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func doLogin(user string, password string) error {
	body := strings.NewReader("parameters%5Blogin%5D=" + user + "&parameters%5Bpassword%5D=" + password)
	req, err := http.NewRequest("POST", "https://app.planrockr.com/rpc/v1/authentication/login", body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "planrockr-cli")

	resp, err := getDefaultClient(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return errors.New("Invalid credentials")
	}
	buf, _ := ioutil.ReadAll(resp.Body)

	err = config.Init()
	if err != nil {
		return errors.New("Error reading config file")
	}

	type AuthData struct {
		Token         string
		Refresh_Token string
	}
	var authData AuthData
	err = json.Unmarshal(buf, &authData)
	if err != nil {
		return errors.New("Error parsing authorization data")
	}
	err = config.Set("auth.token", authData.Token)
	err = config.Set("auth.refreshtoken", authData.Refresh_Token)
	if err != nil {
		return errors.New("Error writing config file")
	}

	fmt.Println("Authorized")

	return nil
}

func init() {
	RootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(&authUser, "user", "u", "", "User e-mail")
	authCmd.Flags().StringVarP(&authPassword, "password", "p", "", "User password")
}
