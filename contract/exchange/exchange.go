package exchange

import (
	"context"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type SeedsExchConfigTable struct {
	SeedsPerUsd   eos.Asset `json:"seeds_per_usd"`
	TlosPerUsd    eos.Asset `json:"tlos_per_usd"`
	CitizenLimit  eos.Asset `json:"citizen_limit"`
	ResidentLimit eos.Asset `json:"resident_limit"`
	VisitorLimit  eos.Asset `json:"visitor_limit"`
}

func (m *SeedsExchConfigTable) String() string {
	return fmt.Sprintf(
		`
		SeedsExchConfigTable {
			SeedsPerUsd: %v,
			TlosPerUsd: %v,
			CitizenLimit: %v,
			ResidentLimit: %v,
			VisitorLimit: %v,	
		}
		`,
		m.SeedsPerUsd,
		m.TlosPerUsd,
		m.CitizenLimit,
		m.ResidentLimit,
		m.VisitorLimit,
	)
}

// SeedsPriceHistory ...
type SeedsPriceHistory struct {
	ID       uint64        `json:"id"`
	SeedsUSD eos.Asset     `json:"seeds_usd"`
	Date     eos.TimePoint `json:"date"`
}

func (m *SeedsPriceHistory) String() string {
	return fmt.Sprintf(
		`
		SeedsPriceHistory {
			ID: %v,
			SeedsUSD: %v,
			Date: %v,
		}
		`,
		m.ID,
		m.SeedsUSD,
		m.Date,
	)
}

type SeedsExchange struct {
	*contract.Contract
}

func NewSeedsExchange(eos *service.EOS, contractName string) *SeedsExchange {
	return &SeedsExchange{
		&contract.Contract{
			EOS:          eos,
			ContractName: contractName,
		},
	}
}

func (m *SeedsExchange) Reset() (*eos.PushTransactionFullResp, error) {
	resp, err := m.ExecAction(m.ContractName, "reset", nil)
	if err != nil {
		return nil, fmt.Errorf("failed calling reset action of contract: %v, error: %v", m.ContractName, err)
	}
	return resp, nil
}

func (m *SeedsExchange) LoadSeedsTablesFromProd(prodEndpoint string) error {
	prodApi := *eos.New(prodEndpoint)

	var config []SeedsExchConfigTable
	var request eos.GetTableRowsRequest
	request.Code = "tlosto.seeds"
	request.Scope = "tlosto.seeds"
	request.Table = "config"
	request.Limit = 1
	request.JSON = true
	response, _ := prodApi.GetTableRows(context.Background(), request)
	response.JSONToStructs(&config)

	fmt.Println("Copying tlosto.seeds configuration table from production")
	// fmt.Println("Config: ", config[0])
	_, err := m.ExecAction(m.ContractName, "updateconfig", config[0])
	if err != nil {
		return fmt.Errorf("failed copying tlosto.seeds configuration table from production, error: %v", err)
	}

	var priceHistory []SeedsPriceHistory
	var request2 eos.GetTableRowsRequest
	request2.Code = "tlosto.seeds"
	request2.Scope = "tlosto.seeds"
	request2.Table = "pricehistory"
	request2.Limit = 1000
	request2.JSON = true
	response2, _ := prodApi.GetTableRows(context.Background(), request2)
	response2.JSONToStructs(&priceHistory)

	fmt.Println("Copying tlosto.seeds price history records from production")
	for _, record := range priceHistory {
		// fmt.Println("Price: ", record)
		_, err = m.ExecAction(m.ContractName, "inshistory", record)
		if err != nil {
			return fmt.Errorf("failed copying tlosto.seeds price history table from production, error: %v", err)
		}
	}
	return nil
}
