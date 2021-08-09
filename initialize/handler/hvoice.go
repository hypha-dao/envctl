package handler

import (
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type tokenCreate struct {
	Issuer         eos.AccountName
	MaxSupply      eos.Asset
	DecayPeriod    eos.Uint64
	DecayPerPeriod eos.Uint64
}

type HVoice struct {
	EOS           *service.EOS
	TokenContract *contract.TokenContract
}

func NewHVoice(eos *service.EOS) *HVoice {
	return &HVoice{
		EOS:           eos,
		TokenContract: contract.NewTokenContract(eos),
	}
}

func (m *HVoice) Handle(data map[interface{}]interface{}, config map[interface{}]interface{}, initOp InitializeOp) error {
	contract := data["contract"].(string)
	issuer := data["issuer"].(string)
	maxSupply := data["max-supply"].(string)
	decayPeriod := data["decay-period"].(int)
	decayPerPeriod := data["decay-per-period"].(int)
	hvoiceMaxSupply, err := eos.NewAssetFromString(maxSupply)
	if err != nil {
		return fmt.Errorf("failed to parse hvoice max supply: %v, error: %v", maxSupply, err)
	}
	fmt.Printf("Creating HVoice contract: %v, issuer: %v, max supply: %v, decay period: %v, decay per period: %v\n", contract, issuer, maxSupply, decayPeriod, decayPerPeriod)
	_, err = m.TokenContract.CreateTokenBase(contract, tokenCreate{
		Issuer:         eos.AN(issuer),
		MaxSupply:      hvoiceMaxSupply,
		DecayPeriod:    eos.Uint64(decayPeriod),
		DecayPerPeriod: eos.Uint64(decayPerPeriod),
	}, false)
	if err != nil {
		return fmt.Errorf("failed to create hvoidce token, error %v", err)
	}
	fmt.Printf("Created token: %v, supply: %v, decayPeriod: %v, decayPerPeriod: %v\n", contract, maxSupply, decayPeriod, decayPerPeriod)
	return nil
}
