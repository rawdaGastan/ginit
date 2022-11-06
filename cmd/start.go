/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rawdaGastan/ginit/internal"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start execution of procfile",
	Long:  `Execute each proc in the procfile with order`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")

		procfile, err := cmd.Flags().GetString("procfile")
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		envfile, err := cmd.Flags().GetString("env")
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		logger := zerolog.New(os.Stdout).With().Logger()

		ginit := internal.NewGinitService(procfile, envfile, args, logger)

		err = ginit.Start(context.Background())

		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP("procfile", "f", "", "Enter your procfile path")
	startCmd.Flags().StringP("env", "e", "", "Enter your env path")
}
