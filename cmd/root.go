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
	"os"

	"github.com/hpcloud/tail"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var tailFile string
var debugFlag bool

//func Kill(err error) {
//	log.Fatal(err)
//}

// TailLog will echo the contents of a given file to stdout.
func TailLog() {
	t, err := tail.TailFile(tailFile, tail.Config{ReOpen: true, Follow: true, MustExist: true, Poll: true})
	if err != nil {
		fmt.Printf("Failed to tail \"%s\": %s\n", tailFile, err)
		os.Exit(4)
		//log.Fatal("Failed to tail \"", tailFile, "\" ", err)

	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "spinner",
	Short: "Spinner is a Service Monitor for Windows that also tails a log file and echoes to Stdout.",
	Long: `Spinner is meant to be uses as the entrypoint/cmd for Windows containers
where the process can't be invoked directly. This ensures that if the service stops
the container will terminate and your container orchestration can take necessary
steps to restart the application.`,
	Example: `spinner.exe service W3SVC -t c:\\iislog\\W3SVC\\u_extend1.log
spinner.exe site http://localhost -t c:\\iislog\\W3SVC\\u_extend1.log
spinner.exe iis -t c:\\iislog\\W3SVC\\u_extend1.log`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spinner.yaml)")

	RootCmd.PersistentFlags().StringVarP(&tailFile, "tail", "t", "", "Path to file to tail and pipe to STDOUT.")
	RootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Print debug logging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".spinner" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".spinner")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
