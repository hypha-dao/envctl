package daobot

import (
	"context"
	"fmt"
	"strconv"

	dao "github.com/hypha-dao/dao-contracts/dao-go"
	"go.uber.org/zap"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
)

// EnrollMembers ...
func EnrollMembers(ctx context.Context, api *eos.API, contract eos.AccountName) error {

	// re-enroll members
	index := 1
	for index < 6 {

		memberNameIn := "mem" + strconv.Itoa(index) + ".hypha"
		//memberNameIn := "member" + strconv.Itoa(index)

		newMember, err := enrollMember(ctx, api, contract, eos.AN(memberNameIn))
		if err != nil {
			return fmt.Errorf("unable to enroll member : "+string(memberNameIn)+": %v ", err)
		}
		zlog.Info("Member enrolled : " + string(memberNameIn) + " with hash: " + newMember.Hash.String())
		index++
	}

	johnnyhypha, err := enrollMember(ctx, api, contract, eos.AN("johnnyhypha1"))
	if err != nil {
		return fmt.Errorf("unable to enroll member johnnyhypha1 : %v ", err)
	}

	zlog.Info("Member enrolled : johnnyhypha1 with hash: " + johnnyhypha.Hash.String())
	return nil
}

func enrollMember(ctx context.Context, api *eos.API, contract, member eos.AccountName) (docgraph.Document, error) {
	zlog.Info("Enrolling account", zap.String("account-name", string(member)))

	trxID, err := dao.Apply(ctx, api, contract, member, "apply to DAO")
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("error applying %v", err)
	}
	zlog.Info("Completed the apply transaction: " + trxID)

	// e.DefaultPause("Building block...")

	_, err = dao.Enroll(ctx, api, contract, contract, member)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("error enrolling %v", err)
	}

	// e.DefaultPause("Building block...")
	memberDoc, err := docgraph.GetLastDocumentOfEdge(ctx, api, contract, "member")
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("error enrolling %v", err)
	}

	return memberDoc, nil
}
