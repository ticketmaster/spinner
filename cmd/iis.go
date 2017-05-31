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
	"log"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func appcmd(c string) string {
	cmd := exec.Command(`cmd`)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: `/C %windir%\system32\inetsrv\appcmd list ` + c,
	}
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out[:]))
}

func iis() {
	s := appcmd("SITE")
	p := appcmd("APPPOOL")
	fmt.Println("Monitoring the following website(s):\n", s)
	fmt.Println("Monitoring the following application pool(s):\n", p)

	for {

		s := appcmd("SITE")
		match, _ := regexp.MatchString("(state:)(Started)", s)

		if match != true {
			fmt.Println("Website not in a running state")
			log.Fatal("Current state:", s)
		} else if debugFlag {
			fmt.Println("Current state:", s)
		}

		p := appcmd("APPPOOL")
		scanner := bufio.NewScanner(strings.NewReader(p))

		for scanner.Scan() {
			match, _ := regexp.MatchString("(state:)(Started)", scanner.Text())

			if match != true {
				fmt.Println("App pool not in a running state")
				log.Fatal("Current state:", scanner.Text())
			} else if debugFlag {
				fmt.Println("Current state:", scanner.Text())
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

// iisCmd represents the iis command
var iisCmd = &cobra.Command{
	Use:     "iis",
	Short:   "Monitors IIS",
	Example: "spinner.exe iis -t c:\\iislog\\W3SVC\\u_extend1.log",
	Long: `Will monitor IIS websites and application pools. Will
terminate if any sites or application pools are not found in a
'Started' state.

Use this as the entrypoint for a container to stop the container if
the given service stops.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debugFlag {
			fmt.Println("with debug")
		}
		if tailFile != "" {
			go TailLog()
		}

		iis()

	},
}

func init() {
	RootCmd.AddCommand(iisCmd)

}
