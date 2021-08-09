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

	"github.com/hypha-dao/envctl/domain"
	"github.com/hypha-dao/envctl/e"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// backendStopCmd represents the backend stop command
var backendStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops backend environment",
	Long:  "Stops backend environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		bkd := domain.NewBackend(viper.GetString("BackendConfigDir"), e.EOS)
		destroy, _ := cmd.Flags().GetBool("destroy")
		if destroy {
			fmt.Println("Destroying backend services...")
			err := bkd.Destroy()
			if err != nil {
				return err
			}
			zlog.Info("Backend services destroyed.")
		} else {
			fmt.Println("Stopping backend services...")
			err := bkd.Stop()
			if err != nil {
				return err
			}
			zlog.Info("Backend services stopped.")
		}

		return nil
	},
}

func init() {
	backendStopCmd.Flags().BoolP("destroy", "d", false, "Destroy backend and delete all data")
	backendCmd.AddCommand(backendStopCmd)
}
