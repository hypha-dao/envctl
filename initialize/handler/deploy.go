package handler

import (
	"fmt"
	"path/filepath"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type Deploy struct {
	EOS           *service.EOS
	TokenContract *contract.TokenContract
	PublicKey     *ecc.PublicKey
}

func NewDeploy(eos *service.EOS, publicKey *ecc.PublicKey) *Deploy {
	return &Deploy{
		EOS:           eos,
		TokenContract: contract.NewTokenContract(eos),
		PublicKey:     publicKey,
	}
}

func (m *Deploy) Handle(data map[interface{}]interface{}, config map[interface{}]interface{}, initOp InitializeOp) error {
	basePath := config["base-path"].(string)
	path := data["path"].(string)
	fileName := data["file-name"].(string)
	account := data["account"].(string)
	fullPath := filepath.Join(basePath, path, fileName)
	fmt.Printf("Deploying contract: %v, to account: %v\n", fullPath, account)
	_, err := m.EOS.SetContract(account, fmt.Sprintf("%v.wasm", fullPath), fmt.Sprintf("%v.abi", fullPath), m.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to deploy contract %v, error %v", account, err)
	}
	err = m.EOS.SetEOSIOCode(account, m.PublicKey)
	if err != nil {
		return err
	}

	if supplyI, ok := data["supply"]; ok {
		supply := supplyI.(string)
		asset, err := eos.NewAssetFromString(supply)
		if err != nil {
			return fmt.Errorf("failed to parse asset: %v error %v", supply, err)
		}
		issuer := account
		if issuerI, ok := data["issuer"]; ok {
			issuer = issuerI.(string)
		}
		_, err = m.TokenContract.CreateToken(account, issuer, asset, false)
		if err != nil {
			return fmt.Errorf("failed to create token: %v, supply: %v error %v", account, supply, err)
		}
		fmt.Printf("Created token: %v, supply: %v\n", account, supply)

		if issue, ok := data["issue"]; ok {
			_, err = m.TokenContract.Issue(account, issuer, issue.(string), "Initial issue")
			if err != nil {
				return fmt.Errorf("failed to issue token: %v, issue amount: %v error %v", account, supply, err)
			}
			fmt.Printf("Issueed token: %v, issue amount: %v\n", account, issue)
		}

	}
	return nil
}
