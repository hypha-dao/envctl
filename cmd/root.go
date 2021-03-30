// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/hypha-dao/envctl/e"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Version represents the cmd command version
var Version string
var cfgFile string

var yamlDefault = []byte(`
#EosioEndpoint: https://test.telos.kitchen
AssetsAsFloat: true
Contract: dao.hypha
UserAccount: johnnyhypha1
Pause: 1s
VotingPeriodDuration: 30s
PayPeriodDuration: 5m
global-expiration: 10
RootHash: 52a7ff82bd6f53b31285e97d6806d886eefb650e79754784e9d923d3df347c91
PrivateKey: xxx
`)

const mainnetChainId = "4667b205c6838ef70ff7988f6e8257e8be0e1284a2f59699054a018f743b1d11"
const testnetChainId = "1eaa0824707c8c16bd25145493bf062aecddfeb56c736f6ba6397f3195f33c9f"
const appName = "envctl"

// RootCmd represents the cli command
var RootCmd = &cobra.Command{
	Use:   "envctl",
	Short: "envctl is for managing DAO environments",
}

// Execute executes the configured RootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		zap.S().DPanic(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Not implemnted
	//RootCmd.PersistentFlags().BoolP("debug", "", false, "Enables verbose API debug messages")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./envctl.yaml)")
}

func networkWarning() {
	colorRed := "\033[31m"
	colorReset := "\033[0m"
	info, err := e.E().A.GetInfo(context.Background())
	if err != nil {
		zap.S().Fatal(string(colorRed) + "ERROR: Unable to get " + e.E().AppName + " Blockchain Node info. Please check the EosioEndpoint configuration.")
	}

	if hex.EncodeToString(info.ChainID) == mainnetChainId {
		fmt.Println(string(colorRed) + "\nERROR: Endpoint is connected to the Telos mainnet - cannot run envctl there. Please change your EOSIO endpoint configuration.")
		fmt.Println(string(colorReset))
		os.Exit(1)
	} else if hex.EncodeToString(info.ChainID) == testnetChainId {
		zap.S().Info("\nNETWORK: Connecting to the Test Network")
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			zap.S().Warn("Cannot find home directory looking for config file", zap.Error(err))
		}
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(home)
		viper.SetConfigName(appName)
	}

	viper.SetEnvPrefix("ENVCTL")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	recurseViperCommands(RootCmd, nil)

	if err := viper.ReadInConfig(); err == nil {
		zap.S().Debug("Using config file", zap.String("config-file", viper.ConfigFileUsed()))
	} else {
		err := viper.ReadConfig(bytes.NewBuffer(yamlDefault))
		if err != nil {
			zap.S().Fatal("No configuration file and error reading default config", zap.Error(err))
		}
	}

	e := e.E()
	if e == nil {
		zap.S().Fatal("unable to configure environment - E() is nil")
	}
	networkWarning()
	SetupLogger()
}

func recurseViperCommands(root *cobra.Command, segments []string) {
	// Stolen from: github.com/abourget/viperbind
	var segmentPrefix string
	if len(segments) > 0 {
		segmentPrefix = strings.Join(segments, "-") + "-"
	}

	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + "global-" + f.Name
		err := viper.BindPFlag(newVar, f)
		if err != nil {
			zap.S().Error("Cannot bind PFlags to variables", zap.Error(err))
		}
	})
	root.Flags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + "cmd-" + f.Name
		err := viper.BindPFlag(newVar, f)
		if err != nil {
			zap.S().Error("Cannot bind PFlags to variables", zap.Error(err))
		}
	})

	for _, cmd := range root.Commands() {
		recurseViperCommands(cmd, append(segments, cmd.Name()))
	}
}
