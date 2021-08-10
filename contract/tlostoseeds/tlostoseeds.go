package tlostoseeds

import (
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type TlosToSeedsContract struct {
	*contract.Contract
}

func NewTlosToSeedsContract(eos *service.EOS, contractName string) *TlosToSeedsContract {
	return &TlosToSeedsContract{
		&contract.Contract{
			EOS:          eos,
			ContractName: contractName,
		},
	}
}

func (m *TlosToSeedsContract) Reset() (*eos.PushTransactionFullResp, error) {
	resp, err := m.ExecAction(m.ContractName, "reset", nil)
	if err != nil {
		return nil, fmt.Errorf("failed calling reset action of contract: %v, error: %v", m.ContractName, err)
	}
	return resp, nil
}
