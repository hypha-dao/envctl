package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	eosc "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
)

var EOSIOKey = "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"

var retries = 3

type EOS struct {
	API *eosc.API
}

func NewEOSFromUrl(url string) *EOS {
	api := eosc.New(url)
	api.SetSigner(&eosc.KeyBag{})
	return NewEOS(api)
}

func NewEOS(api *eosc.API) *EOS {
	return &EOS{
		API: api,
	}
}

func (m *EOS) AddKey(privateKey string) (*ecc.PublicKey, error) {
	// logger.Infof("PKey: %v", pkey)
	key, err := ecc.NewPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("privateKey parameter is not a valid format: %s", err)
	}
	err = m.API.Signer.ImportPrivateKey(context.Background(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed importing key, error: %v", err)
	}
	publicKey := key.PublicKey()
	return &publicKey, nil
}

func (m *EOS) AddEOSIOKey() (*ecc.PublicKey, error) {
	return m.AddKey(EOSIOKey)
}

func (m *EOS) Trx(retries int, actions ...*eosc.Action) (*eosc.PushTransactionFullResp, error) {
	// for _, action := range actions {
	// 	logger.Infof("Trx Account: %v Name: %v, Authorization: %v, Data: %v", action.Account, action.Name, action.Authorization, action.ActionData)

	// }
	resp, err := m.API.SignPushActions(context.Background(), actions...)
	if err != nil {
		if retries > 0 {
			if isRetryableError(err) {
				time.Sleep(3 * time.Second)
				return m.Trx(retries-1, actions...)
			}
		}
		return nil, fmt.Errorf("failed to push trx: %v, error: %v", actions, err)
	}
	return resp, nil
}

func isRetryableError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "connection reset by peer") ||
		strings.Contains(errMsg, "Transaction took too long")
}

func (m *EOS) SimpleTrx(contract, action, actor string, data interface{}) (*eosc.PushTransactionFullResp, error) {
	return m.Trx(retries, buildAction(contract, action, actor, data))
}

func (m *EOS) DebugTrx(contract, action, actor string, data interface{}) (*eosc.PushTransactionFullResp, error) {
	txOpts := &eosc.TxOptions{}
	if err := txOpts.FillFromChain(context.Background(), m.API); err != nil {
		return nil, err
	}

	tx := eosc.NewTransaction([]*eosc.Action{buildAction(contract, action, actor, data)}, txOpts)
	signedTx, packedTx, err := m.API.SignTransaction(context.Background(), tx, txOpts.ChainID, eosc.CompressionNone)
	if err != nil {
		return nil, err
	}

	content, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		return nil, err
	}

	fmt.Println(string(content))
	return m.API.PushTransaction(context.Background(), packedTx)
}

func (m *EOS) CreateAccount(accountName string, publicKey *ecc.PublicKey, failIfExists bool) (eosc.AccountName, error) {
	account := eosc.AccountName(accountName)
	_, err := m.Trx(retries, system.NewNewAccount("eosio", account, *publicKey))
	if err != nil {
		if failIfExists || !strings.Contains(err.Error(), "name is already taken") {
			return "", err
		}
	}
	return account, nil
}

//SetContract sets contract, if publicKey is not nil it creates the account
func (m *EOS) SetContract(accountName, wasmFile, abiFile string, publicKey *ecc.PublicKey) (*eosc.PushTransactionFullResp, error) {

	if publicKey != nil {
		_, err := m.CreateAccount(accountName, publicKey, false)
		if err != nil {
			return nil, err
		}
	}
	account := eosc.AccountName(accountName)
	setCodeAction, err := system.NewSetCode(account, wasmFile)
	if err != nil {
		return nil, fmt.Errorf("unable construct set_code action: %v", err)
	}

	setAbiAction, err := system.NewSetABI(account, abiFile)
	if err != nil {
		return nil, fmt.Errorf("unable construct set_abi action: %v", err)
	}
	resp, err := m.Trx(retries, setCodeAction, setAbiAction)
	if err != nil {
		if !strings.Contains(err.Error(), "contract is already running this version of code") {
			return nil, err
		}
	}
	return resp, nil
}

func (m *EOS) SetEOSIOCode(accountName string, publicKey *ecc.PublicKey) error {
	acct := eosc.AN(accountName)
	codePermissionAction := system.NewUpdateAuth(acct,
		"active",
		"owner",
		eosc.Authority{
			Threshold: 1,
			Keys: []eosc.KeyWeight{{
				PublicKey: *publicKey,
				Weight:    1,
			}},
			Accounts: []eosc.PermissionLevelWeight{{
				Permission: eosc.PermissionLevel{
					Actor:      acct,
					Permission: "eosio.code",
				},
				Weight: 1,
			}},
			Waits: []eosc.WaitWeight{},
		}, "owner")

	_, err := m.Trx(retries, codePermissionAction)
	if err != nil {
		return fmt.Errorf("error setting eosio.code permission for account: %v, error: %v", accountName, err)
	}
	return nil
}

func (m *EOS) GetTableRows(request eosc.GetTableRowsRequest, rows interface{}) error {

	request.JSON = true
	response, err := m.API.GetTableRows(context.Background(), request)
	if err != nil {
		return fmt.Errorf("get table rows %v", err)
	}

	err = response.JSONToStructs(rows)
	if err != nil {
		return fmt.Errorf("json to structs %v", err)
	}
	return nil
}

func (m *EOS) GetBalance(account eosc.AccountName, symbol eosc.Symbol, contract eosc.AccountName) (*eosc.Asset, error) {

	assets, err := m.API.GetCurrencyBalance(context.Background(), account, symbol.MustSymbolCode().String(), contract)
	if err != nil {
		return nil, fmt.Errorf("failed to get currency balance, error: %v", err)
	}
	if len(assets) > 0 {
		return &assets[0], nil
	}
	return nil, nil
}

func buildAction(contract, action, actor string, data interface{}) *eosc.Action {
	return &eosc.Action{
		Account: eosc.AN(contract),
		Name:    eosc.ActN(action),
		Authorization: []eosc.PermissionLevel{
			{
				Actor:      eosc.AccountName(actor),
				Permission: eosc.PN("active"),
			},
		},
		ActionData: eosc.NewActionData(data),
	}
}