package daobot

import (
	"context"

	eostest "github.com/digital-scarcity/eos-go-test"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/hypha-dao/envctl/e"
	"github.com/hypha-dao/envctl/pretend"
)

type claimNext struct {
	AssignmentHash eos.Checksum256 `json:"assignment_hash"`
}

func ClaimNextPeriod(ctx context.Context, api *eos.API, contract, claimer eos.AccountName, assignment docgraph.Document) (string, error) {

	actions := []*eos.Action{{
		Account: contract,
		Name:    eos.ActN("claimnextper"),
		Authorization: []eos.PermissionLevel{
			{Actor: claimer, Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(claimNext{
			AssignmentHash: assignment.Hash,
		}),
	}}

	trxID, err := eostest.ExecTrx(ctx, api, actions)

	if err != nil {
		e.Pause(pretend.PayPeriodDuration(), "", "Waiting for a period to lapse")

		actions := []*eos.Action{{
			Account: contract,
			Name:    eos.ActN("claimnextper"),
			Authorization: []eos.PermissionLevel{
				{Actor: claimer, Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(claimNext{
				AssignmentHash: assignment.Hash,
			}),
		}}

		trxID, err = eostest.ExecTrx(ctx, api, actions)
	}

	return trxID, err
}
