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
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	s "strings"
	"time"

	"github.com/spf13/cobra"
)

var urlFlag string

func queryPage(u, cl, tte string) {
	var c int64
	var fullURL string

	if s.HasPrefix(u, "http://") || s.HasPrefix(u, "https://") {
		fullURL = u
	} else {
		fullURL = "http://" + u
	}

	log.Println("Full URL being monitored:", fullURL, "Counter limit:", cl)

	for {
		cl, err := strconv.ParseInt(cl, 10, 8)
		if err != nil {
			log.Fatal(err)
		}

		tte, err := strconv.ParseInt(tte, 10, 8)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get(fullURL)
		if err != nil {
			log.Fatal("An error occurred during the request:", err)
		}

		if resp.StatusCode >= 300 {
			c++
			if c >= cl {
				rd, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				// Convert rd to string and remove carriage returns from output
				rs := s.Replace(string(rd), "\r", "", -1)
				// Remove line feeds from output
				rs = s.Replace(rs, "\n", "", -1)

				// Output to file if path is set
				if outFile != "" {
					d, _ := filepath.Split(outFile)
					err := os.MkdirAll(d, os.ModePerm)
					if err != nil {
						panic(err)
					}

					f, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
					if err != nil {
						panic(err)
					}
					defer f.Close()

					ls := time.Now().Format("2006-01-02 15:04:05") + " [spinner-app] " + rs + "\n"

					b := bufio.NewWriterSize(f, 10000)
					_, err = b.WriteString(ls)
					if err != nil {
						log.Println(err)
					}

					err = b.Flush()
					if err != nil {
						log.Println(err)
					}
				}

				// Output to stdout
				log.Println("Status Code:", resp.StatusCode, " Body:", rs)
				// Sleep for n second(s) allowing output error to stdout to be picked up
				// by monitoring software/container (if needed)
				time.Sleep(time.Duration(tte) * time.Second)
				log.Fatalln("Spinner shutting down, status code was:", resp.StatusCode)
			} else {
				log.Println("Status Code:", resp.StatusCode, "Counter count:", c)
			}
		} else if debugFlag {
			c = 0
			log.Println("Status Code:", resp.StatusCode)
		}

		resp.Body.Close()

		time.Sleep(1 * time.Second)
	}
}

// siteCmd represents the site command
var siteCmd = &cobra.Command{
	Use:     "site [url] [counter limit] [time to exit]",
	Short:   "Watch a Site",
	Aliases: []string{"url", "address"},
	Example: "spinner.exe site http://localhost 5 5 -t c:\\iislog\\W3SVC\\u_extend1.log -o c:\\logs\\spinner.log",
	Long: `Poll Web Site by Get request and terminate this process if the a >300
status code is returned.

Counter Limit (default 1) is the number of times the site being monitored can
be down before spinner exits.

Time to Exit (default 1) is the time (in seconds) after the response body is
logged to stdout before spinner will shutdown. This can be useful
if the monitoring software does not catch the error quick enough.

OutFile is the location path to a file to write out the response body if
the response status is >= 300. This can be useful in troubleshooting why
the application is exiting.

Use this as the entrypoint for a container to stop the container if
the given service stops.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("site needs a url for the command")
			os.Exit(1)
		}
		if debugFlag {
			fmt.Println("with debug")
		}
		if tailFile != "" {
			go TailLog()
		}
		if len(args) == 1 {
			queryPage(args[0], "1", "1")
		} else if len(args) == 2 {
			queryPage(args[0], args[1], "1")
		} else {
			queryPage(args[0], args[1], args[2])
		}
	},
}

func init() {
	RootCmd.AddCommand(siteCmd)

	//siteCmd.Flags().StringVarP(&urlFlag, "url", "u", "http://localhost/", "Url to watch")

}
