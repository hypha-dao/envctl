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
	"time"

	eostest "github.com/digital-scarcity/eos-go-test"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
	"github.com/hypha-dao/dao-contracts/dao-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/hypha-dao/envctl/contract/exchange"
	"github.com/hypha-dao/envctl/daobot"
	"github.com/hypha-dao/envctl/e"
	"github.com/hypha-dao/envctl/pretend"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// populatePretendCmd populates the environment with the known Pretend environment
var populatePretendCmd = &cobra.Command{
	Use:   "pretend",
	Short: "populates with the known Pretend environment",
	Long:  "populates with the known Pretend environment",
	RunE: func(cmd *cobra.Command, args []string) error {

		// e.DefaultPause("Warming up...")
		fmt.Printf("Setting bank permissions...\n")
		err := setBankPermissions()
		if err != nil {
			return fmt.Errorf("failed setting bank permissions, error: %v", err)
		}

		fmt.Printf("Initializing %v contract...\n", e.E().Exchange)
		exchange := exchange.NewSeedsExchange(e.EOS, string(e.E().Exchange))
		err = exchange.LoadSeedsTablesFromProd(viper.GetString("prod-eosio-endpoint"))
		if err != nil {
			return fmt.Errorf("failed initializing %v contract, error: %v", e.E().Exchange, err)
		}
		_, err = dao.CreateRoot(e.E().X, e.E().A, e.E().Contract)
		if err != nil {
			return fmt.Errorf("cannot create root document for pretend environment: %v ", err)
		}

		root, err := docgraph.LoadDocument(e.E().X, e.E().A, e.E().Contract, "52a7ff82bd6f53b31285e97d6806d886eefb650e79754784e9d923d3df347c91")
		if err != nil {
			return fmt.Errorf("cannot load root document for pretend environment: %v ", err)
		}

		_, err = dao.SetIntSetting(e.E().X, e.E().A, e.E().Contract, "voting_duration_sec", int64(pretend.VotingPeriodDuration().Round(time.Second))/1000000000)
		if err != nil {
			return fmt.Errorf("cannot set int setting for pretend environment: %v ", err)
		}
		fmt.Println("Set setting: voting_duration_sec	: ", int64(pretend.VotingPeriodDuration().Round(time.Second))/1000000000)

		settings, err := pretend.DefaultSettings()
		if err != nil {
			return fmt.Errorf("cannot retrieve default settings for pretend environment: %v ", err)
		}
		for _, setting := range settings {
			fmt.Println("Setting a setting. Key: " + setting.Label + " Value: " + setting.Value.String())
			_, err = dao.SetSetting(e.E().X, e.E().A, e.E().Contract, setting.Label, setting.Value)
			if err != nil {
				return fmt.Errorf("cannot set setting: %v ", err)
			}
		}

		fmt.Println("Adding "+strconv.Itoa(viper.GetInt("PeriodCount"))+" periods with duration 		: ", pretend.PayPeriodDuration())
		_, err = dao.AddPeriods(e.E().X, e.E().A, e.E().Contract, root.Hash, viper.GetInt("PeriodCount"), pretend.PayPeriodDuration())
		if err != nil {
			return fmt.Errorf("cannot add periods: %v ", err)
		}

		// fmt.Println("Adding start edge to root document...")
		// _, err = docgraph.CreateEdge(e.E().X, e.E().A, e.E().Contract, e.E().Contract, root.Hash, periods[0].Hash, eos.Name("start"))
		// if err != nil {
		// 	return fmt.Errorf("failed adding start edge: %v,  to root document: %v, error: %v", root.Hash, periods[0].Hash, err)
		// }
		fmt.Println("Enrolling members...")
		err = daobot.EnrollMembers(e.E().X, e.E().A, e.E().Contract)
		if err != nil {
			return fmt.Errorf("failed enrolling members: %v ", err)
		}

		d, err := createPretend(e.E().X, e.E().A, e.E().Contract, e.E().TelosDecide, e.E().User)
		if err != nil {
			return fmt.Errorf("cannot create pretend environment: %v ", err)
		}
		fmt.Println("Pretend environment successfully created; assignment document is	: ", d.Hash.String())
		return nil
	},
}

func init() {
	populateCmd.AddCommand(populatePretendCmd)
}

// createPretend returns the assignment document
func createPretend(ctx context.Context, api *eos.API, contract, telosDecide, member eos.AccountName) (docgraph.Document, error) {

	fmt.Println("Creating Role...")
	role, err := daobot.CreateRole(ctx, api, contract, telosDecide, member, []byte(pretend.Role))
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("unable to create role: %v", err)
	}
	fmt.Println("Role document successfully created	: ", role.Hash.String())
	// e.DefaultPause("Building a block...")

	fmt.Println("Creating assignment...")
	roleAssignment, err := daobot.CreateAssignment(ctx, api, contract, telosDecide, member, eos.Name("role"), eos.Name("assignment"), []byte(pretend.Assignment))
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("unable to create assignment: %v", err)
	}
	fmt.Println("Created role assignment document	: ", roleAssignment.Hash.String(), "Waiting for period to lapse: ", pretend.PayPeriodDuration())
	time.Sleep(pretend.PayPeriodDuration())
	// e.Pause(pretend.PayPeriodDuration()+e.E().Pause, "", "Waiting for a period to lapse")

	fmt.Println("Claiming period...")
	_, err = daobot.ClaimNextPeriod(ctx, api, contract, member, roleAssignment)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("unable to claim pay on assignment: %v %v", roleAssignment.Hash.String(), err)
	}
	fmt.Println("Claimed pay on the assignment 		: ", roleAssignment.Hash.String())
	// e.DefaultPause("Building a block...")

	fmt.Println("Creating Payout...")
	payAmt, _ := eos.NewAssetFromString("1000.00 USD")
	payout, err := daobot.CreatePayout(ctx, api, contract, telosDecide, member, member, payAmt, 50, []byte(pretend.Payout))
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("unable to create payout: %v", err)
	}
	fmt.Println("Created payout document	: ", payout.Hash.String())
	// e.DefaultPause("Building a block...")

	fmt.Println("Creating badge...")
	badge, err := daobot.CreateBadge(ctx, api, contract, telosDecide, member, []byte(pretend.Badge))
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("unable to create badge: %v", err)
	}
	fmt.Println("Created badge document	: ", badge.Hash.String())
	// e.DefaultPause("Building a block...")

	fmt.Println("Creating badge assignment...")
	badgeAssignment, err := daobot.CreateAssignment(ctx, api, contract, telosDecide, member, eos.Name("badge"), eos.Name("assignbadge"), []byte(pretend.BadgeAssignment))
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("unable to create badge assignment: %v", err)
	}
	fmt.Println("Created badge assignment document	: ", badgeAssignment.Hash.String())

	return roleAssignment, nil
}

func setBankPermissions() error {

	bankPublicKey, err := toPublic(eostest.DefaultKey())
	if err != nil {
		return fmt.Errorf("unable to derive public key: %v %v", eostest.DefaultKey(), err)
	}

	bankPermissionActions := []*eos.Action{system.NewUpdateAuth(e.E().Bank,
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
						Actor:      e.E().DAO,
						Permission: "eosio.code",
					},
					Weight: 1,
				}},
			Waits: []eos.WaitWeight{},
		}, "owner")}

	_, err = e.ExecWithRetry(e.E().X, e.E().A, bankPermissionActions)
	if err != nil {
		return fmt.Errorf("unable to update bank account permissions: %v %v", viper.GetString("Bank"), err)
	}
	return nil
}
