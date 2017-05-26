// +build windows

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
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var serviceFlag string

func testState(s svc.State) error {

	switch s {
	case svc.Running:
		if debugFlag {
			log.Println("Running")
		}
		return nil
	case svc.Stopped:
		return errors.New("Stopped")
	case svc.StartPending:
		return errors.New("StartPending")
	case svc.StopPending:
		return errors.New("StopPending")
	case svc.ContinuePending:
		return errors.New("ContinuePending")
	case svc.PausePending:
		return errors.New("PausePending")
	case svc.Paused:
		return errors.New("Paused")
	}

	return nil
}

// QueryService returns the the status of the given service
// as a uint32. Stopped=1; StartPending=2; StopPending=3; Running=4
// ContinuePending=5; PausePending=6; Paused=7
func QueryService(s string) {

	manager, err := mgr.Connect()

	if err != nil {
		log.Fatal("Cannot connect to manager: ", err)
	}
	defer manager.Disconnect()

	for {

		service, err := manager.OpenService(s)
		if err != nil {
			fmt.Printf("service does not exist: %s\n", err)
			os.Exit(2)
		}
		defer service.Close()

		var status svc.Status
		status, err = service.Query()
		if err != nil {
			fmt.Printf("failed to get service status: %s\n", err)
			os.Exit(3)
		}

		stateError := testState(status.State)

		if stateError != nil {
			fmt.Printf("service status: %s\n", stateError)
			os.Exit(10)
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:     "service [name]",
	Short:   "Watch a Windows Service",
	Example: "spinner.exe service W3SVC -t c:\\iislog\\W3SVC\\u_extend1.log",
	Long: `Poll the state of a Windows Service and terminate this process if
State does not equal "Running".

Use this as the entrypoint for a container to stop the container if
the given service stops.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("service needs a name for the command")
			os.Exit(1)
		}
		if tailFile != "" {
			go TailLog()
		}
		if debugFlag {
			fmt.Println("with debug")
		}

		QueryService(args[0])

	},
}

func init() {
	RootCmd.AddCommand(serviceCmd)

}
