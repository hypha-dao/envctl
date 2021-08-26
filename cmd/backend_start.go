// The MIT License (MIT)

// Copyright (c) 2020, Digital Scarcity

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package cmd ...
package cmd

import (
	"fmt"
	"time"

	"github.com/hypha-dao/envctl/domain"
	"github.com/hypha-dao/envctl/e"
	"github.com/hypha-dao/envctl/initialize"
	"github.com/hypha-dao/envctl/initialize/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// backendStartCmd represents the backend start command
var backendStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts backend environment",
	Long:  "Starts backend environment, and deploys dhos contracts",
	RunE: func(cmd *cobra.Command, args []string) error {
		bkd := domain.NewBackend(viper.GetString("BackendConfigDir"), e.EOS)

		restart, _ := cmd.Flags().GetBool("restart")
		if restart {
			fmt.Println("Destroying backend services...")
			err := bkd.Destroy()
			if err != nil {
				return err
			}
			fmt.Println("Backend services have been destroyed. Restarting...")
		}
		fmt.Println("Starting backend services...")
		err := bkd.Start()
		if err != nil {
			return err
		}
		fmt.Println("Waiting for backend services to finish their setup...")
		time.Sleep(time.Minute)
		fmt.Println("Backend services started. Initializing...")
		initOp := handler.InitializeOp_Start
		if restart {
			initOp = handler.InitializeOp_Restart
		}
		err = initialize.Initialize(viper.Get("backend-init-settings").([]interface{}), e.EOS, initOp)
		if err != nil {
			return err
		}
		zlog.Info("Deployed contracts and created accounts.")
		return nil
	},
}

func init() {
	backendStartCmd.Flags().BoolP("restart", "r", false, "Starts a clean environment")
	backendCmd.AddCommand(backendStartCmd)
}
