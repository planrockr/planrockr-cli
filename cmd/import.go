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
	"fmt"

	"github.com/planrockr/planrockr-cli/pkg"
	"github.com/spf13/cobra"
)

var (
	importType     string
	importServer   string
	importUser     string
	importPassword string
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import a project to Planrockr",
	Long:  "You can import a project from Jira or Gitlab",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch importType {
		case "jira":
			err = pkg.JiraImport(importServer, importUser, importPassword)
		case "gitlab":
			fmt.Println("Not implemented yet")
		}
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&importType, "type", "t", "", "Import source")
	importCmd.Flags().StringVarP(&importServer, "server", "s", "", "Server address")
	importCmd.Flags().StringVarP(&importUser, "user", "u", "", "User e-mail")
	importCmd.Flags().StringVarP(&importPassword, "password", "p", "", "User password")
}
