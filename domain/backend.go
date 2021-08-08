package domain

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dfuse-io/logging"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/hypha-dao/envctl/contract"
	"github.com/hypha-dao/envctl/service"
	"go.uber.org/zap"
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

var zlog *zap.Logger

func init() {
	logging.Register("github.com/hypha-dao/envctl/domain", &zlog)
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

func (m *Backend) Init(initSettings map[string]interface{}, restart bool) error {
	publicKey, err := m.EOS.AddEOSIOKey()
	if err != nil {
		return err
	}
	err = m.checkoutRepos(initSettings["checkout-repos"].(map[string]interface{}))
	if err != nil {
		return err
	}
	err = m.buildContracts(initSettings["build-contracts"].(map[string]interface{}), restart)
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

func (m *Backend) checkoutRepos(checkout map[string]interface{}) error {
	basePath := checkout["base-path"].(string)
	repos := checkout["repos"].([]interface{})
	for _, repoI := range repos {
		repo := repoI.(map[interface{}]interface{})
		url := repo["url"].(string)
		branch := repo["branch"].(string)
		fmt.Printf("Checking out repo: %v, branch: %v\n", url, branch)
		err := service.CheckoutRepo(basePath, url, branch)
		if err != nil {
			return fmt.Errorf("failed to check out repo: %v, branch: %v, error %v", url, branch, err)
		}
	}
	return nil
}

func (m *Backend) buildContracts(build map[string]interface{}, restart bool) error {
	basePath := build["base-path"].(string)
	repos := build["repos"].([]interface{})
	for _, repoI := range repos {
		repo := repoI.(map[interface{}]interface{})
		name := repo["name"].(string)
		repoPath := path.Join(basePath, name)
		fmt.Printf("Building repo: %v\n", repoPath)
		buildPath := path.Join(repoPath, "build")
		err := os.Mkdir(buildPath, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
		isDirEmpty, err := service.IsDirEmpty(buildPath)
		if err != nil {
			return fmt.Errorf("failed to check if build path: %v is empty, error: %v", buildPath, err)
		}
		if restart || isDirEmpty {
			err = service.ExecCmd(&service.CmdArgs{
				Name: "cmake",
				Args: []string{".."},
				Dir:  buildPath,
			})
			if err != nil {
				return fmt.Errorf("failed running cmake for repo: %v, build path: %v, error %v", repoPath, buildPath, err)
			}

			err = service.ExecCmd(&service.CmdArgs{
				Name: "make",
				Dir:  buildPath,
			})
			if err != nil {
				return fmt.Errorf("failed running make for repo: %v, build path: %v, error %v", repoPath, buildPath, err)
			}
		}
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
		fmt.Printf("Deploying contract: %v, to account: %v\n", fullPath, account)
		_, err := m.EOS.SetContract(account, fmt.Sprintf("%v.wasm", fullPath), fmt.Sprintf("%v.abi", fullPath), publicKey)
		if err != nil {
			return fmt.Errorf("failed to deploy contract %v, error %v", account, err)
		}
		err = m.EOS.SetEOSIOCode(account, publicKey)
		if err != nil {
			return err
		}

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
		zlog.Info("Created account", zap.String("account-name", account))
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
	fmt.Println("Backend config dir: ", m.ConfigDir)
	err := service.ExecCmd(&service.CmdArgs{
		Name: "docker-compose",
		Args: args,
		Dir:  m.ConfigDir,
	})
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
