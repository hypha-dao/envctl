package domain

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/hypha-dao/envctl/contract"
	"github.com/hypha-dao/envctl/service"
)

type Backend struct {
	ConfigDir     string
	EOS           *service.EOS
	TokenContract *contract.TokenContract
}

func NewBackend(configDir string, eos *service.EOS) *Backend {

	return &Backend{
		ConfigDir:     configDir,
		EOS:           eos,
		TokenContract: contract.NewTokenContract(eos),
	}
}

func (m *Backend) Start() error {
	return m.dockerCmd("up", "-d")
}

func (m *Backend) Stop() error {
	return m.dockerCmd("stop")
}

func (m *Backend) Destroy() error {
	return m.dockerCmd("down", "-v")
}

func (m *Backend) Init(initSettings map[string]interface{}) error {
	publicKey, err := m.EOS.AddEOSIOKey()
	if err != nil {
		return err
	}
	err = m.deployContracts(initSettings["deploy"].(map[string]interface{}), publicKey)
	if err != nil {
		return err
	}
	hvoiceMaxSupply, _ := eos.NewAssetFromString("-1.00 HVOICE")
	err = m.createHVoiceToken("voice.hypha", "dao.hypha", hvoiceMaxSupply, 1, 100000)
	if err != nil {
		return err
	}
	accounts := m.getAccountNames(initSettings["accounts"].([]interface{}))
	err = m.createAccounts(accounts, publicKey)
	if err != nil {
		return err
	}
	return nil
}

func (m *Backend) deployContracts(deploy map[string]interface{}, publicKey *ecc.PublicKey) error {
	basePath := deploy["base-path"].(string)
	contracts := deploy["contracts"].([]interface{})
	for _, contractI := range contracts {
		contract := contractI.(map[interface{}]interface{})
		path := contract["path"].(string)
		fileName := contract["file-name"].(string)
		account := contract["account"].(string)
		fullPath := filepath.Join(basePath, path, fileName)
		_, err := m.EOS.SetContract(account, fmt.Sprintf("%v.wasm", fullPath), fmt.Sprintf("%v.abi", fullPath), publicKey)
		if err != nil {
			return fmt.Errorf("failed to deploy contract %v, error %v", account, err)
		}
		err = m.EOS.SetEOSIOCode(account, publicKey)
		if err != nil {
			return err
		}
		fmt.Println("Deployed contract: ", account)

		if supplyI, ok := contract["supply"]; ok {
			supply := supplyI.(string)
			asset, err := eos.NewAssetFromString(supply)
			if err != nil {
				return fmt.Errorf("failed to parse asset: %v error %v", supply, err)
			}
			_, err = m.TokenContract.CreateToken(account, account, asset, false)
			if err != nil {
				return fmt.Errorf("failed to create token: %v, supply: %v error %v", account, supply, err)
			}
			fmt.Printf("Created token: %v, supply: %v\n", account, supply)
		}

	}
	return nil
}

func (m *Backend) createAccounts(accounts []string, publicKey *ecc.PublicKey) error {
	for _, account := range accounts {
		_, err := m.EOS.CreateAccount(account, publicKey, false)
		if err != nil {
			return err
		}
		fmt.Println("Created account: ", account)
	}
	return nil
}

func (m *Backend) getAccountNames(accounts []interface{}) []string {
	accountNames := make([]string, 0)
	for _, accountI := range accounts {
		account := accountI.(map[interface{}]interface{})
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
	}
	return accountNames
}

func (m *Backend) createHVoiceToken(contract, issuer string,
	maxSupply eos.Asset, decayPeriod eos.Uint64, decayPerPeriod eos.Uint64) error {
	type tokenCreate struct {
		Issuer         eos.AccountName
		MaxSupply      eos.Asset
		DecayPeriod    eos.Uint64
		DecayPerPeriod eos.Uint64
	}
	_, err := m.TokenContract.CreateTokenBase(contract, tokenCreate{
		Issuer:         eos.AN(issuer),
		MaxSupply:      maxSupply,
		DecayPeriod:    decayPeriod,
		DecayPerPeriod: decayPerPeriod,
	}, false)
	if err != nil {
		return fmt.Errorf("failed to create hvoidce token, error %v", err)
	}
	fmt.Printf("Created token: %v, supply: %v, decayPeriod: %v, decayPerPeriod: %v\n", contract, maxSupply, decayPeriod, decayPerPeriod)
	return nil
}

// func (m *Backend) createDataDirs() error {
// 	err := m.createDir(m.DfuseDataDir)
// 	if err != nil {
// 		return err
// 	}
// 	return m.createDir(m.DGraphDataDir)
// }

// func (m *Backend) deleteDataDirs() error {
// 	err := m.deleteDir(m.DfuseDataDir)
// 	if err != nil {
// 		return err
// 	}
// 	return m.deleteDir(m.DGraphDataDir)
// }

// func (m *Backend) createDir(dir string) error {
// 	err := os.MkdirAll(dir, 0777)
// 	if err != nil {
// 		return fmt.Errorf("Error creating dir: %v, error: %v", dir, err)
// 	}
// 	return nil
// }

// func (m *Backend) deleteDir(dir string) error {
// 	err := os.RemoveAll(dir)
// 	if err != nil {
// 		return fmt.Errorf("Error deleting dir: %v, error: %v", dir, err)
// 	}
// 	return nil
// }

// func (m *Backend) dirExist(dir string) (bool, error) {
// 	_, err := os.Stat(dir)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return true, nil
// }

func (m *Backend) dockerCmd(args ...string) error {
	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = m.ConfigDir
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	fmt.Println("Docker output: ", stdout.String())
	fmt.Println("Docker error output: ", stderr.String())
	if err != nil {
		return fmt.Errorf("error running docker-compose command with args: %v, error: %v", args, err)
	}
	return nil
}

// func (m *Backend) updateEnvVars() error {

// 	envPath := filepath.Join(m.ConfigDir, ".env")
// 	envVars, err := godotenv.Read(envPath)
// 	if err != nil {
// 		return fmt.Errorf("Error reading dho backend .env file, error: %v", err)
// 	}
// 	envVars["DFUSE_DATA_DIR"] = m.DfuseDataDir
// 	envVars["DGRAPH_DATA_DIR"] = m.DGraphDataDir

// 	err = godotenv.Write(envVars, envPath)
// 	if err != nil {
// 		return fmt.Errorf("Error writing dho backend .env file, error: %v", err)
// 	}
// 	return nil
// }

// func (m *Backend) getEnvVar(key string) (string, error) {
// 	value, ok := m.EnvVars[key]
// 	if !ok || strings.Trim(value, " ") == "" {
// 		return "", fmt.Errorf("%v var not set in backend .env file")
// 	}
// 	return value, nil
// }
