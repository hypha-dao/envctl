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
	"context"
	"fmt"
	"strconv"

	eostest "github.com/digital-scarcity/eos-go-test"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/hypha-dao/envctl/e"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// initCmd represents the erase command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new nodeos instance for population",
	Long:  "initialize a new nodeos instance for population",
	RunE: func(cmd *cobra.Command, args []string) error {

		zlog.Info("init command assumes local instance for now...")

		restartCmd, err := eostest.RestartNodeos(true)
		if err != nil {
			return fmt.Errorf("unable to restart nodeos: %v", err)
		}

		zlog.Info("(Re)started node", zap.String("nodeos-pid", strconv.Itoa(restartCmd.Process.Pid)))

		var daoHome = viper.GetString("DAOHome")
		var daoPrefix = daoHome + "/build/dao/dao."
		artifactsHome := daoHome + "/dao-go/artifacts"
		treasuryPrefix := artifactsHome + "/treasury/treasury."
		monitorPrefix := artifactsHome + "/monitor/monitor."
		voicePrefix := daoHome + "/../voice-token/build/voice/voice."
		// voicePrefix := artifactsHome + "/voice/voice."

		e.Env.DAO, err = eostest.CreateAccountFromString(e.E().X, e.E().A, viper.GetString("DAO"), eostest.DefaultKey())
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", viper.GetString("DAO"), err)
		}
		zlog.Info("Created account", zap.String("account-name", string(e.Env.DAO)))

		e.Env.HusdToken, err = eostest.CreateAccountFromString(e.Env.X, e.Env.A, viper.GetString("HusdToken"), eostest.DefaultKey())
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", viper.GetString("HusdToken"), err)
		}
		zlog.Info("Created account", zap.String("account-name", string(e.Env.HusdToken)))

		e.Env.HyphaToken, err = eostest.CreateAccountFromString(e.Env.X, e.Env.A, viper.GetString("HyphaToken"), eostest.DefaultKey())
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", viper.GetString("HyphaToken"), err)
		}
		zlog.Info("Created account", zap.String("account-name", string(e.Env.HyphaToken)))

		e.Env.HvoiceToken, err = eostest.CreateAccountFromString(e.Env.X, e.Env.A, viper.GetString("HvoiceToken"), eostest.DefaultKey())
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", viper.GetString("HvoiceToken"), err)
		}
		zlog.Info("Created account", zap.String("account-name", string(e.Env.HvoiceToken)))

		e.Env.Bank, err = eostest.CreateAccountFromString(e.Env.X, e.Env.A, viper.GetString("Bank"), eostest.DefaultKey())
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", viper.GetString("Bank"), err)
		}
		zlog.Info("Created account", zap.String("account-name", string(e.Env.Bank)))

		e.Env.Events, err = eostest.CreateAccountFromString(e.Env.X, e.Env.A, viper.GetString("Events"), eostest.DefaultKey())
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", viper.GetString("Events"), err)
		}
		zlog.Info("Created account", zap.String("account-name", string(e.Env.Events)))

		// e.Env.TelosDecide, err = eostest.CreateAccountFromString(e.Env.X, e.Env.A, viper.GetString("TelosDecide"), eostest.DefaultKey())
		// if err != nil {
		// 	return fmt.Errorf("unable to create account from string: %v %v", viper.GetString("TelosDecide"), err)
		// }
		// zlog.Info("Created account		: ", e.Env.TelosDecide)

		bankPublicKey, err := toPublic(eostest.DefaultKey())
		if err != nil {
			return fmt.Errorf("unable to derive public key: %v %v", eostest.DefaultKey(), err)
		}

		bankPermissionActions := []*eos.Action{system.NewUpdateAuth(e.Env.Bank,
			"active",
			"owner",
			eos.Authority{
				Threshold: 1,
				Keys: []eos.KeyWeight{{
					PublicKey: bankPublicKey,
					Weight:    1,
				}},
				Accounts: []eos.PermissionLevelWeight{
					{
						Permission: eos.PermissionLevel{
							Actor:      e.Env.Bank,
							Permission: "eosio.code",
						},
						Weight: 1,
					},
					{
						Permission: eos.PermissionLevel{
							Actor:      e.Env.DAO,
							Permission: "eosio.code",
						},
						Weight: 1,
					}},
				Waits: []eos.WaitWeight{},
			}, "owner")}

		trxId, err := e.ExecWithRetry(e.Env.X, e.Env.A, bankPermissionActions)
		if err != nil {
			return fmt.Errorf("unable to update bank account permissions: %v %v", viper.GetString("Bank"), err)
		}
		zlog.Info("Updated permissions", zap.String("account-name", string(e.Env.Bank)), zap.String("trx-id", trxId))

		trxId, err = eostest.SetContract(e.Env.X, e.Env.A, e.Env.DAO, daoPrefix+"wasm", daoPrefix+"abi")
		if err != nil {
			return fmt.Errorf("unable to set contract on DAO: %v %v", e.Env.DAO, err)
		}
		zlog.Info("Deployed contract", zap.String("account-name", string(e.Env.DAO)), zap.String("trx-id", trxId), zap.String("wasm", daoPrefix+"wasm"), zap.String("abi", daoPrefix+"abi"))

		trxId, err = eostest.SetContract(e.Env.X, e.Env.A, e.Env.Bank, treasuryPrefix+"wasm", treasuryPrefix+"abi")
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", e.Env.Bank, err)
		}
		zlog.Info("Deployed contract", zap.String("account-name", string(e.Env.Bank)), zap.String("trx-id", trxId), zap.String("wasm", treasuryPrefix+"wasm"), zap.String("abi", treasuryPrefix+"abi"))

		trxId, err = eostest.SetContract(e.Env.X, e.Env.A, e.Env.HvoiceToken, voicePrefix+"wasm", voicePrefix+"abi")
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", e.Env.HvoiceToken, err)
		}
		zlog.Info("Deployed contract", zap.String("account-name", string(e.Env.HvoiceToken)), zap.String("trx-id", trxId), zap.String("wasm", voicePrefix+"wasm"), zap.String("abi", voicePrefix+"abi"))

		// trxId, err = eostest.SetContract(e.Env.X, e.Env.A, e.Env.TelosDecide, decidePrefix+"wasm", decidePrefix+"abi")
		// if err != nil {
		// 	return fmt.Errorf("unable to create account from string: %v %v", e.Env.TelosDecide, err)
		// }
		// zlog.Info("Deployed Telos Decide contract to 				: ", e.Env.TelosDecide, "	;  TrxID: ", trxId)

		trxId, err = eostest.SetContract(e.Env.X, e.Env.A, e.Env.Events, monitorPrefix+"wasm", monitorPrefix+"abi")
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", e.Env.Events, err)
		}
		zlog.Info("Deployed contract", zap.String("account-name", string(e.Env.Events)), zap.String("trx-id", trxId), zap.String("wasm", monitorPrefix+"wasm"), zap.String("abi", monitorPrefix+"abi"))

		husdMaxSupply, _ := eos.NewAssetFromString("1000000000.00 HUSD")
		trxId, err = deployAndCreateToken(e.Env.X, e.Env.A, artifactsHome, e.Env.HusdToken, e.Env.Bank, husdMaxSupply)
		if err != nil {
			return fmt.Errorf("unable to deploy and create HUSD token: %v", err)
		}
		zlog.Info("Created Token", zap.String("supply", husdMaxSupply.String()), zap.String("trx-id", trxId), zap.String("issuer", string(e.Env.Bank)), zap.String("token-contract", string(e.Env.HusdToken)))

		hyphaMaxSupply, _ := eos.NewAssetFromString("1000000000.00 HYPHA")
		trxId, err = deployAndCreateToken(e.Env.X, e.Env.A, artifactsHome, e.Env.HyphaToken, e.Env.DAO, hyphaMaxSupply)
		if err != nil {
			return fmt.Errorf("unable to deploy and create HYPHA token: %v", err)
		}
		zlog.Info("Created Token", zap.String("supply", hyphaMaxSupply.String()), zap.String("trx-id", trxId), zap.String("issuer", string(e.Env.DAO)), zap.String("token-contract", string(e.Env.HyphaToken)))

		// Hvoice doesn't have any limit (max supply should be -1)
		hvoiceMaxSupply, _ := eos.NewAssetFromString("-1.00 HVOICE")
		trxId, err = createHVoiceToken(e.Env.X, e.Env.A, e.Env.HvoiceToken, e.Env.DAO, hvoiceMaxSupply, 1, 100000)
		if err != nil {
			return fmt.Errorf("unable to deploy and create HVOICE token: %v", err)
		}
		zlog.Info("Created Token", zap.String("supply", hvoiceMaxSupply.String()), zap.String("trx-id", trxId), zap.String("issuer", string(e.Env.DAO)), zap.String("token-contract", string(e.Env.HvoiceToken)))

		index := 1
		for index < 6 {

			memberName := "mem" + strconv.Itoa(index) + ".hypha"

			member, err := eostest.CreateAccountFromString(e.Env.X, e.Env.A, memberName, eostest.DefaultKey())
			if err != nil {
				return fmt.Errorf("unable to create account from string: %v %v", memberName, err)
			}
			zlog.Info("Created account", zap.String("account-name", memberName))

			e.Env.Members = append(e.Env.Members, member)
			index++
		}

		johnnyhypha, err := eostest.CreateAccountFromString(e.Env.X, e.Env.A, "johnnyhypha1", eostest.DefaultKey())
		if err != nil {
			return fmt.Errorf("unable to create account from string: %v %v", "johnnyhypha1", err)
		}

		e.Env.Members = append(e.Env.Members, johnnyhypha)

		zlog.Debug("Nodeos initialization complete")
		return nil
	},
}

func createHVoiceToken(ctx context.Context, api *eos.API, contract, issuer eos.AccountName,
	maxSupply eos.Asset, decayPeriod eos.Uint64, decayPerPeriod eos.Uint64) (string, error) {
	type tokenCreate struct {
		Issuer         eos.AccountName
		MaxSupply      eos.Asset
		DecayPeriod    eos.Uint64
		DecayPerPeriod eos.Uint64
	}

	actions := []*eos.Action{{
		Account: contract,
		Name:    eos.ActN("create"),
		Authorization: []eos.PermissionLevel{
			{Actor: contract, Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(tokenCreate{
			Issuer:         issuer,
			MaxSupply:      maxSupply,
			DecayPeriod:    decayPeriod,
			DecayPerPeriod: decayPerPeriod,
		}),
	}}

	return e.ExecWithRetry(ctx, api, actions)
}

type tokenCreate struct {
	Issuer    eos.AccountName
	MaxSupply eos.Asset
}

func deployAndCreateToken(ctx context.Context, api *eos.API, tokenHome string,
	contract, issuer eos.AccountName, maxSupply eos.Asset) (string, error) {

	tokenWasm := tokenHome + "/token/token.wasm"
	tokenAbi := tokenHome + "/token/token.abi"

	trxId, err := eostest.SetContract(ctx, api, contract, tokenWasm, tokenAbi)
	if err != nil {
		return "", fmt.Errorf("cannot set contract from: %v %v", tokenHome, err)
	}
	zlog.Info("Deployed contract", zap.String("account-name", string(e.Env.Events)), zap.String("trx-id", trxId), zap.String("wasm", tokenWasm), zap.String("abi", tokenAbi))

	actions := []*eos.Action{{
		Account: contract,
		Name:    eos.ActN("create"),
		Authorization: []eos.PermissionLevel{
			{Actor: contract, Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(tokenCreate{
			Issuer:    issuer,
			MaxSupply: maxSupply,
		}),
	}}

	return e.ExecWithRetry(ctx, api, actions)
}

func toPublic(privateKey string) (ecc.PublicKey, error) {

	key, err := ecc.NewPrivateKey(privateKey)
	if err != nil {
		return ecc.PublicKey{}, fmt.Errorf("privateKey parameter is not a valid format: %s", err)
	}

	return key.PublicKey(), nil
}

func init() {
	RootCmd.AddCommand(initCmd)
}
