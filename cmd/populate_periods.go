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
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/hypha-dao/envctl/e"
	"github.com/hypha-dao/envctl/pretend"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// populatePeriodsCmd populates the environment with the known Pretend environment
var populatePeriodsCmd = &cobra.Command{
	Use:   "periods",
	Short: "appends periods to the current calendar",
	Long:  "appends periods to the current calendar",
	RunE: func(cmd *cobra.Command, args []string) error {

		e.DefaultPause("Warming up...")

		rootDocument, err := docgraph.LoadDocument(e.E().X, e.E().A, e.E().Contract, viper.GetString("RootHash"))
		if err != nil {
			return fmt.Errorf("cannot load root document: %v", err)
		}

		startEdges, err := docgraph.GetEdgesFromDocumentWithEdge(e.E().X, e.E().A, e.E().Contract, rootDocument, eos.Name("start"))
		if err != nil {
			return fmt.Errorf("error while retrieving start edge: %v", err)
		}
		if len(startEdges) == 0 {
			return fmt.Errorf("no start edge from the root node exists: %v", err)
		}

		startPeriodDoc, err := docgraph.LoadDocument(e.E().X, e.E().A, e.E().Contract, startEdges[0].ToNode.String())
		if err != nil {
			return fmt.Errorf("error loading the start period document: %v", err)
		}

		period, err := models.NewPeriod(e.E().X, e.E().A, e.E().Contract, startPeriodDoc)
		if err != nil {
			return fmt.Errorf("cannot convert document to period type: %v", err)
		}

		for period.Next != nil {
			period = *period.Next
		}

		zlog.Info("Adding periods", zap.Int("pay-period-count", viper.GetInt("PeriodCount")), zap.Duration("pay-period-duration", pretend.PayPeriodDuration()))
		_, err = addPeriods(e.E().X, e.E().A, e.E().Contract, period, viper.GetInt("PeriodCount"), pretend.PayPeriodDuration())
		if err != nil {
			return fmt.Errorf("cannot add periods: %v ", err)
		}

		return nil
	},
}

type addPeriod struct {
	Predecessor eos.Checksum256 `json:"predecessor"`
	StartTime   eos.TimePoint   `json:"start_time"`
	Label       string          `json:"label"`
}

func addPeriods(ctx context.Context, api *eos.API, daoContract eos.AccountName,
	predecessor models.Period,
	numPeriods int,
	periodDuration time.Duration) ([]models.Period, error) {

	periods := make([]models.Period, numPeriods)

	zlog.Info("\nAdding periods: " + strconv.Itoa(len(periods)))
	bar := e.DefaultProgressBar(len(periods), "")

	for i := 0; i < numPeriods; i++ {
		startTime := eos.TimePoint(predecessor.StartTime.Add(periodDuration).UnixNano()/1000 + 1)
		addPeriodAction := eos.Action{
			Account: daoContract,
			Name:    eos.ActN("addperiod"),
			Authorization: []eos.PermissionLevel{
				{Actor: daoContract, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(addPeriod{
				Predecessor: predecessor.Document.Hash,
				StartTime:   startTime,
				Label:       "period #" + strconv.Itoa(i+1) + " of " + strconv.Itoa(numPeriods),
			}),
		}

		//startTime = eos.TimePoint(predecessor.StartTime.Add(periodDuration).UnixNano() / 1000)
		// marker = marker.Add(periodDuration).Add(time.Millisecond)

		_, err := eostest.ExecWithRetry(ctx, api, []*eos.Action{&addPeriodAction})
		if err != nil {
			return periods, fmt.Errorf("cannot add period: %v", err)
		}

		lastDoc, _ := docgraph.GetLastDocument(ctx, api, daoContract)

		periods[i], err = models.NewPeriod(e.E().X, e.E().A, e.E().Contract, lastDoc)
		if err != nil {
			return periods, fmt.Errorf("cannot convert document to period type: %v", err)
		}

		predecessor = periods[i]
		time.Sleep(e.Env.Pause)
		bar.Add(1)
	}

	return periods, nil
}

func init() {
	populateCmd.AddCommand(populatePeriodsCmd)
}
