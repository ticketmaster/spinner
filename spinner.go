package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"errors"

	"github.com/hpcloud/tail"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const usage = `Spinner, A Service Monitor and file tail executable for Windows Containers.
Usage:
  spinner.exe flags
The flags are:
  --service <name>    Service name to Watch
  --path <path>       Path to log file to tail
  --usage             Print usage
  --debug             Print debug logging
  --version           Print version of this binary
Examples:
  # Watch IIS and tail the access log:
  .\spinner.exe -service W3SVC -path c:\iislog\W3SVC\u_extend1.log
Note:
  Spinner will exit if the file does not already exist. Make sure to generate an event
  and flush the buffer immediately or wait for the file to generate before starting
  Spinner.

  CMD Start-Service W3SVC; ` + "`" + `
    Invoke-WebRequest http://localhost -UseBasicParsing | Out-Null; ` + "`" + `
    netsh http flush logbuffer | Out-Null; ` + "`" + `
    spinner.exe -service W3SVC -path c:\iislog\W3SVC\u_extend1.log
`

var version = "0.0.0"
var serviceFlag = flag.String("service", "", "Service name to Watch")
var pathFlag = flag.String("path", "", "Path to log file to tail")
var debugFlag = flag.Bool("debug", false, "Enabled Debug Output")
var versionFlag = flag.Bool("version", false, "Show version info")
var usageFlag = flag.Bool("usage", false, "Print usage")

func usageExit(rc int) {
	fmt.Println(usage)
	os.Exit(rc)
}

func kill(err error) {
	log.Fatal(err)
}

func testState(s svc.State) error {

	switch s {
	case svc.Running:
		if *debugFlag {
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

		service, err := manager.OpenService(*serviceFlag)
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
			kill(stateError)
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

// TailLog will echo the contents of a given file to stdout.
func TailLog() {
	t, err := tail.TailFile(*pathFlag, tail.Config{ReOpen: true, Follow: true, MustExist: true, Poll: true})
	if err != nil {
		log.Fatal("Failed to tail \"", *pathFlag, "\" ", err)
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

func main() {
	flag.Usage = func() { usageExit(0) }
	flag.Parse()

	switch {
	case *usageFlag:
		flag.Usage()
	case *versionFlag:
		fmt.Printf("spinner v%s\n", version)
		return
	case *serviceFlag == "":
		fmt.Println("--service is required. See \"spinner -h\" for usage.")
		return
	case *pathFlag == "":
		fmt.Println("--path is required. See \"spinner -h\" for usage.")
		return
	default:
		go QueryService()
		TailLog()
	}

}
