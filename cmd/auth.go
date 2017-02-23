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
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/planrockr/planrockr-cli/config"
	"github.com/spf13/cobra"
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
	errConfig := config.Init()
	if errConfig != nil {
		return errors.New("Error reading config file")
	}
	conf := config.Get()
	q, err := url.ParseQuery("parameters[login]=" + user + "&parameters[password]=" + password)
	if err != nil {
		return err
	}
	body := strings.NewReader(q.Encode())
	req, err := http.NewRequest("POST", conf.BaseUrl+"/rpc/v1/authentication/login", body)
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

	type AuthData struct {
		Token         string
		Refresh_Token string
		Id            int
	}
	var authData AuthData
	err = json.Unmarshal(buf, &authData)
	if err != nil {
		return errors.New("Error parsing authorization data")
	}
	err = config.SetString("auth.token", authData.Token)
	if err != nil {
		return errors.New("Error writing config file")
	}
	err = config.SetString("auth.refreshtoken", authData.Refresh_Token)
	if err != nil {
		return errors.New("Error writing config file")
	}
	err = config.SetInt("auth.id", authData.Id)
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
