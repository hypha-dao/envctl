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

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/hypha-dao/envctl/e"
	"github.com/spf13/cobra"
)

// eraseCmd represents the erase command
var eraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "erase everything from the test environment",
	Long:  "erase everything from the test environment",
	RunE: func(cmd *cobra.Command, args []string) error {

		eraseAllDocuments(e.E().X, e.E().A, e.E().Contract)
		eraseAllEdges(e.E().X, e.E().A, e.E().Contract)

		fmt.Println("Erased all documents and edges")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(eraseCmd)
}

type eraseDoc struct {
	Hash eos.Checksum256 `json:"hash"`
}

func eraseAllDocuments(ctx context.Context, api *eos.API, contract eos.AccountName) {

	documents, err := docgraph.GetAllDocuments(ctx, api, contract)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nErasing documents: " + strconv.Itoa(len(documents)))
	bar := e.DefaultProgressBar(len(documents), "erasing documents...")

	// docType := eos.Name("unknown")
	for _, document := range documents {

		// typeFV, err := document.GetContent("type")
		// if err != nil {
		// 	docType = eos.Name("unknown")
		// } else {
		// 	docType = typeFV.Impl.(eos.Name)
		// }

		// if docType == eos.Name("settings") ||
		// 	docType == eos.Name("dho") {
		// 	// do not erase
		// 	fmt.Println("\nSkipping document because type of " + string(docType) + " : " + document.Hash.String())
		// } else {
		actions := []*eos.Action{{
			Account: contract,
			Name:    eos.ActN("erasedoc"),
			Authorization: []eos.PermissionLevel{
				{Actor: contract, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(eraseDoc{
				Hash: document.Hash,
			}),
		}}

		_, err := Exec(ctx, api, actions)
		if err != nil {
			fmt.Println("\nFailed to erase : ", document.Hash.String())
			fmt.Println(err)
		} else {
			time.Sleep(e.E().Pause)
			// }
		}
		bar.Add(1)
	}
}

type killEdge struct {
	EdgeID uint64 `json:"id"`
}

func eraseAllEdges(ctx context.Context, api *eos.API, contract eos.AccountName) {

	edges, err := docgraph.GetAllEdges(ctx, api, contract)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nErasing edges: " + strconv.Itoa(len(edges)))
	bar := e.DefaultProgressBar(len(edges), "erasing edges...")

	for _, edge := range edges {

		actions := []*eos.Action{{
			Account: contract,
			Name:    eos.ActN("killedge"),
			Authorization: []eos.PermissionLevel{
				{Actor: contract, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(killEdge{
				EdgeID: edge.ID,
			}),
		}}

		_, err := Exec(ctx, api, actions)
		if err != nil {
			fmt.Println("\nFailed to erase : ", strconv.Itoa(int(edge.ID)))
			fmt.Println(err)
		} else {
			time.Sleep(e.E().Pause)
		}
	}
	bar.Add(1)

}
