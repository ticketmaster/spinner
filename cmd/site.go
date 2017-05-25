// Copyright © 2017 Ticketmaster
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
	"log"
	"net/http"
	s "strings"
	"time"

	"github.com/spf13/cobra"
)

var urlFlag string

func queryPage() {

	var fullURL string

	if s.HasPrefix(urlFlag, "http://") || s.HasPrefix(urlFlag, "https://") {
		fullURL = urlFlag
	} else {
		fullURL = "http://" + urlFlag
	}

	fmt.Println("Full URL being monitored:", fullURL)
	for {

		resp, err := http.Get(fullURL)

		if err != nil {
			log.Fatal("An error occurred during the request:", err)
		}

		if resp.StatusCode != 200 {
			log.Fatal("Status code != 200")
		} else if debugFlag {
			log.Println("Status Code:", resp.StatusCode)
		}

		resp.Body.Close()

		time.Sleep(1000 * time.Millisecond)
	}
}

// siteCmd represents the site command
var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Watch a Site",
	Long: `Poll Web Site by Get request and terminate this process if
the a >300 status code is returned.

Use this as the entrypoint for a container to stop the container if
the given service stops.`,
	SuggestFor: []string{"url", "address"},
	Example:    "spinner.exe site http://localhost",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("site called ")
		if debugFlag {
			fmt.Println("with debug")
		}
		if urlFlag != "" {
			fmt.Println("url: ", urlFlag)
		}
		go queryPage()
		TailLog()
	},
}

func init() {
	RootCmd.AddCommand(siteCmd)

	siteCmd.Flags().StringVarP(&urlFlag, "url", "u", "http://localhost/", "Url to watch")

}
