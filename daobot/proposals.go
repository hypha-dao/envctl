package daobot

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/eoscanada/eos-go"
	dao "github.com/hypha-dao/dao-contracts/dao-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/hypha-dao/envctl/e"
	"github.com/spf13/viper"
)

// type updateDoc struct {
// 	Hash  eos.Checksum256    `json:"hash"`
// 	Group string             `json:"group"`
// 	Key   string             `json:"key"`
// 	Value docgraph.FlexValue `json:"value"`
// }

// type docGroups struct {
// 	ContentGroups []docgraph.ContentGroup `json:"content_groups"`
// }

func getSettings(ctx context.Context, api *eos.API, contract eos.AccountName) (docgraph.Document, error) {

	rootHash := string("52a7ff82bd6f53b31285e97d6806d886eefb650e79754784e9d923d3df347c91")
	if len(viper.GetString("rootHash")) > 0 {
		rootHash = viper.GetString("rootHash")
	}

	root, err := docgraph.LoadDocument(ctx, api, contract, rootHash)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("root document not found, required for default period %v", err)
	}

	edges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, contract, root, eos.Name("settings"))
	if err != nil || len(edges) <= 0 {
		return docgraph.Document{}, fmt.Errorf("error retrieving settings edge %v", err)
	}

	return docgraph.LoadDocument(ctx, api, contract, edges[0].ToNode.String())
}

// func getDefaultPeriod(api *eos.API, contract eos.AccountName) docgraph.Document {

// 	ctx := context.Background()

// 	root, err := docgraph.LoadDocument(ctx, api, contract, viper.GetString("rootHash"))
// 	if err != nil {
// 		panic("Root document not found, required for default period.")
// 	}

// 	edges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, contract, root, eos.Name("start"))
// 	if err != nil || len(edges) <= 0 {
// 		panic("Next document not found: " + edges[0].ToNode.String())
// 	}

// 	lastDocument, err := docgraph.LoadDocument(ctx, api, contract, edges[0].ToNode.String())
// 	if err != nil {
// 		panic("Next document not found: " + edges[0].ToNode.String())
// 	}

// 	index := 1
// 	for index < 6 {

// 		edges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, contract, lastDocument, eos.Name("next"))
// 		if err != nil || len(edges) <= 0 {
// 			panic("There are no next edges")
// 		}

// 		lastDocument, err = docgraph.LoadDocument(ctx, api, contract, edges[0].ToNode.String())
// 		if err != nil {
// 			panic("Next document not found: " + edges[0].ToNode.String())
// 		}

// 		index++
// 	}
// 	return lastDocument
// }

type proposal struct {
	Proposer      eos.AccountName         `json:"proposer"`
	ProposalType  eos.Name                `json:"proposal_type"`
	ContentGroups []docgraph.ContentGroup `json:"content_groups"`
}

func proposeAndPass(ctx context.Context, api *eos.API,
	contract, telosDecide, proposer eos.AccountName, proposal proposal) (docgraph.Document, error) {
	action := eos.ActN("propose")
	actions := []*eos.Action{{
		Account: contract,
		Name:    action,
		Authorization: []eos.PermissionLevel{
			{Actor: proposer, Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(proposal)}}

	trxID, err := e.ExecWithRetry(ctx, api, actions)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("error proposeAndPass: %v", err)
	}
	fmt.Println("Proposed. Transaction ID: " + trxID)
	e.DefaultPause("Building a block...")

	return closeLastProposal(ctx, api, contract, telosDecide, proposer)
}

func closeLastProposal(ctx context.Context, api *eos.API, contract, telosDecide, member eos.AccountName) (docgraph.Document, error) {

	// retrieve the last proposal
	proposal, err := docgraph.GetLastDocumentOfEdge(ctx, api, contract, eos.Name("proposal"))
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("error retrieving proposal document %v", err)
	}
	fmt.Println("Retrieved proposal document to close: " + proposal.Hash.String())

	_, err = dao.VotePass(ctx, api, contract, telosDecide, member, &proposal)
	if err == nil {
		fmt.Println("Member voted : " + string(member))
	}
	e.DefaultPause("Building a block...")

	_, err = dao.VotePass(ctx, api, contract, telosDecide, eos.AN("alice"), &proposal)
	if err == nil {
		fmt.Println("Member voted : alice")
	}
	e.DefaultPause("Building a block...")

	_, err = dao.VotePass(ctx, api, contract, telosDecide, eos.AN("johnnyhypha1"), &proposal)
	if err == nil {
		fmt.Println("Member voted : johnnyhypha1")
	}
	e.DefaultPause("Building a block...")

	index := 1
	for index < 5 {

		memberNameIn := "mem" + strconv.Itoa(index) + ".hypha"
		//memberNameIn := "member" + strconv.Itoa(index)

		_, err = dao.VotePass(ctx, api, contract, telosDecide, eos.AN(memberNameIn), &proposal)
		if err != nil {
			return docgraph.Document{}, fmt.Errorf("error voting %v", err)
		}
		e.DefaultPause("Building a block...")
		fmt.Println("Member voted : " + string(memberNameIn))
		index++
	}

	settings, err := getSettings(ctx, api, contract)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("cannot retrieve settings document %v", err)
	}

	votingPeriodDuration, err := settings.GetContent("voting_duration_sec")
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("cannot retrieve voting_duration_sec setting %v", err)
	}

	votingPause := time.Duration((5 + votingPeriodDuration.Impl.(int64)) * 1000000000)
	e.Pause(votingPause, "Waiting on voting period to lapse: "+strconv.Itoa(int(5+votingPeriodDuration.Impl.(int64)))+" seconds", "")

	fmt.Println("Closing proposal: " + proposal.Hash.String())
	_, err = dao.CloseProposal(ctx, api, contract, member, proposal.Hash)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("cannot close proposal %v", err)
	}

	e.DefaultPause("Building a block...")
	passedProposal, err := docgraph.GetLastDocumentOfEdge(ctx, api, contract, eos.Name("passedprops"))
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("error retrieving passed proposal document %v", err)
	}
	fmt.Println("Retrieved passed proposal document to close: " + passedProposal.Hash.String())

	return passedProposal, nil
}

func createParent(ctx context.Context, api *eos.API, contract, telosDecide, member eos.AccountName, parentType eos.Name, data []byte) (docgraph.Document, error) {
	var doc docgraph.Document
	err := json.Unmarshal([]byte(data), &doc)
	if err != nil {
		panic(err)
	}

	return proposeAndPass(ctx, api, contract, telosDecide, member, proposal{
		Proposer:      member,
		ProposalType:  parentType,
		ContentGroups: doc.ContentGroups,
	})
}

// CreateRole ...
func CreateRole(ctx context.Context, api *eos.API, contract, telosDecide, member eos.AccountName, data []byte) (docgraph.Document, error) {
	return createParent(ctx, api, contract, telosDecide, member, eos.Name("role"), data)
}

// CreateBadge ...
func CreateBadge(ctx context.Context, api *eos.API, contract, telosDecide, member eos.AccountName, data []byte) (docgraph.Document, error) {
	return createParent(ctx, api, contract, telosDecide, member, eos.Name("badge"), data)
}

// CreateAssignment ...
func CreateAssignment(ctx context.Context, api *eos.API, contract, telosDecide, member eos.AccountName, parentType, assignmentType eos.Name, data []byte) (docgraph.Document, error) {
	var proposalDoc docgraph.Document
	err := json.Unmarshal([]byte(data), &proposalDoc)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("cannot unmarshal error: %v ", err)
	}

	e.DefaultPause("Building block...")

	// e.g. a "role" is parent to a "role assignment"
	// e.g. a "badge" is parent to a "badge assignment"
	var parent docgraph.Document
	parent, err = docgraph.GetLastDocumentOfEdge(ctx, api, contract, parentType)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("cannot retrieve last document of edge: "+string(parentType)+"  - error: %v ", err)
	}

	// inject the parent hash in the first content group of the document
	// TODO: use content_group_label to find details group instead of just the first one
	proposalDoc.ContentGroups[0] = append(proposalDoc.ContentGroups[0], docgraph.ContentItem{
		Label: string(parentType),
		Value: &docgraph.FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: docgraph.GetVariants().TypeID("checksum256"),
				Impl:   parent.Hash,
			}},
	})

	// inject the assignee in the first content group of the document
	// TODO: use content_group_label to find details group instead of just the first one
	proposalDoc.ContentGroups[0] = append(proposalDoc.ContentGroups[0], docgraph.ContentItem{
		Label: "assignee",
		Value: &docgraph.FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: docgraph.GetVariants().TypeID("name"),
				Impl:   member,
			}},
	})

	return proposeAndPass(ctx, api, contract, telosDecide, member, proposal{
		Proposer:      member,
		ProposalType:  assignmentType,
		ContentGroups: proposalDoc.ContentGroups,
	})
}

// //Creates an assignment without the INIT_TIME_SHARE, CURRENT_TIME_SHARE & LAST_TIME_SHARE nodes
// func CreateOldAssignment(t *testing.T, ctx context.Context, api *eos.API, contract, member eos.AccountName, memberDocHash, roleDocHash, startPeriodHash eos.Checksum256, assignment string) (docgraph.Document, error) {

// 	_, err := dao.ProposeAssignment(ctx, api, contract, member, member, roleDocHash, startPeriodHash, assignment)
// 	if err != nil {
// 		return docgraph.Document{}, err
// 	}

// 	// retrieve the document we just created
// 	proposal, err := docgraph.GetLastDocumentOfEdge(ctx, api, contract, eos.Name("proposal"))
// 	if err != nil {
// 		return docgraph.Document{}, err
// 	}

// 	return proposal, nil
// }

// CreatePayout ...
func CreatePayout(ctx context.Context, api *eos.API,
	contract, telosDecide, proposer, recipient eos.AccountName,
	usdAmount eos.Asset, deferred int64, data []byte) (docgraph.Document, error) {

	var payoutDoc docgraph.Document
	err := json.Unmarshal([]byte(data), &payoutDoc)
	if err != nil {
		return docgraph.Document{}, fmt.Errorf("cannot unmarshal error: %v ", err)
	}

	// inject the recipient in the first content group of the document
	// TODO: use content_group_label to find details group instead of just the first one
	payoutDoc.ContentGroups[0] = append(payoutDoc.ContentGroups[0], docgraph.ContentItem{
		Label: "recipient",
		Value: &docgraph.FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: docgraph.GetVariants().TypeID("name"),
				Impl:   recipient,
			}},
	})

	// TODO: use content_group_label to find details group instead of just the first one
	payoutDoc.ContentGroups[0] = append(payoutDoc.ContentGroups[0], docgraph.ContentItem{
		Label: "usd_amount",
		Value: &docgraph.FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: docgraph.GetVariants().TypeID("asset"),
				Impl:   usdAmount,
			}},
	})

	// TODO: use content_group_label to find details group instead of just the first one
	payoutDoc.ContentGroups[0] = append(payoutDoc.ContentGroups[0], docgraph.ContentItem{
		Label: "deferred_perc_x100",
		Value: &docgraph.FlexValue{
			BaseVariant: eos.BaseVariant{
				TypeID: docgraph.GetVariants().TypeID("int64"),
				Impl:   deferred,
			}},
	})

	return proposeAndPass(ctx, api, contract, telosDecide, proposer, proposal{
		Proposer:      proposer,
		ProposalType:  eos.Name("payout"),
		ContentGroups: payoutDoc.ContentGroups,
	})
}
