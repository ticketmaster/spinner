// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
func QueryService() {

	manager, err := mgr.Connect()

	if err != nil {
		log.Fatal("Cannot connect to manager: ", err)
	}
	defer manager.Disconnect()

	for {

		service, err := manager.OpenService(serviceFlag)
		if err != nil {
			log.Fatal("service does not exist:", err)
		}
		defer service.Close()

		var status svc.Status
		status, err = service.Query()
		if err != nil {
			log.Fatal("failed to get service status: ", err)
		}

		stateError := testState(status.State)

		if err != nil {
			Kill(stateError)
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Watch a Windows Service",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service called")
		if debugFlag {
			fmt.Println("with debug")
		}
		if serviceFlag != "" {
			fmt.Println("service: ", serviceFlag)
		}
		go QueryService()
		TailLog()
	},
}

func init() {
	RootCmd.AddCommand(serviceCmd)

	serviceCmd.Flags().StringVarP(&serviceFlag, "name", "n", "", "Service name to Watch")
}
