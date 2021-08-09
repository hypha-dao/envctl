package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type Account struct {
	EOS       *service.EOS
	PublicKey *ecc.PublicKey
}

func NewAccount(eos *service.EOS, publicKey *ecc.PublicKey) *Account {
	return &Account{
		EOS:       eos,
		PublicKey: publicKey,
	}
}

func (m *Account) Handle(data map[interface{}]interface{}, config map[interface{}]interface{}, initOp InitializeOp) error {
	accountNames := m.getAccountNames(data)
	data["accounts"] = accountNames
	if create, ok := config["create"].(bool); !ok || create {
		return m.createAccounts(accountNames)
	}
	return nil
}

func (m *Account) createAccounts(accounts []string) error {
	for _, account := range accounts {
		_, err := m.EOS.CreateAccount(account, m.PublicKey, false)
		if err != nil {
			return err
		}
		fmt.Println("Created account: ", account)
	}
	return nil
}

func (m *Account) getAccountNames(account map[interface{}]interface{}) []string {
	accountNames := make([]string, 0)
	accountName := account["name"].(string)
	var total int
	if totalI, ok := account["total"]; ok {
		total = totalI.(int)
	} else {
		total = 1
	}
	if !strings.Contains(accountName, "0") && total != 1 {
		accountName += "0"
	}
	for i := 1; i <= total; i++ {
		accountNames = append(accountNames, strings.Replace(accountName, "0", strconv.Itoa(i), -1))
	}
	return accountNames
}
