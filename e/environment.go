package e

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	eostest "github.com/digital-scarcity/eos-go-test"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Environment struct {
	A             *eos.API
	X             context.Context
	AppName       string
	Contract      eos.AccountName
	TelosDecide   eos.AccountName
	User          eos.AccountName
	Pause         time.Duration
	DAO           eos.AccountName
	HusdToken     eos.AccountName
	HyphaToken    eos.AccountName
	HvoiceToken   eos.AccountName
	Bank          eos.AccountName
	Events        eos.AccountName
	Members       []eos.AccountName
	GenesisHVOICE int64
}

var once sync.Once
var Env *Environment

func E() *Environment {
	onceBody := func() {

		Env = &Environment{
			A:           eos.New(viper.GetString("EosioEndpoint")),
			X:           context.Background(),
			AppName:     viper.GetString("AppName"),
			Contract:    eos.AN(viper.GetString("Contract")),
			DAO:         eos.AN(viper.GetString("DAO")),
			TelosDecide: eos.AN("trailservice"),
			User:        eos.AN(viper.GetString("UserAccount")),
			Pause:       viper.GetDuration("Pause"),
		}

		keyBag := &eos.KeyBag{}
		keyBag.ImportPrivateKey(context.Background(), "5KCZ9VBJMMiLaAY24Ro66mhx4vU1VcJELZVGrJbkUBATyqxyYmj")
		keyBag.ImportPrivateKey(context.Background(), "5HwnoWBuuRmNdcqwBzd1LABFRKnTk2RY2kUMYKkZfF8tKodubtK")
		keyBag.ImportPrivateKey(context.Background(), eostest.DefaultKey())
		Env.A.SetSigner(keyBag)

		zap.S().Debug("Configured Environment object with sync.Once.Do")
	}
	once.Do(onceBody)
	return Env
}

func DefaultProgressBar(counter int, prefix string) *progressbar.ProgressBar {
	return progressbar.NewOptions(counter,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(90),
		// progressbar.OptionShowIts(),
		progressbar.OptionSetDescription("[cyan]"+fmt.Sprintf("%20v", prefix)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}

func Pause(seconds time.Duration, headline, prefix string) {
	if headline != "" {
		fmt.Println(headline)
	}

	bar := progressbar.NewOptions(100,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(90),
		// progressbar.OptionShowIts(),
		progressbar.OptionSetDescription("[cyan]"+fmt.Sprintf("%20v", prefix)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	chunk := seconds / 100
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(chunk)
	}
	fmt.Println()
	fmt.Println()
}

func DefaultPause(headline string) {
	if headline != "" {
		fmt.Println(headline)
	}

	bar := progressbar.NewOptions(100,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(90),
		// progressbar.OptionShowIts(),
		progressbar.OptionSetDescription("[cyan]"+fmt.Sprintf("%20v", "")),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	chunk := E().Pause / 100
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(chunk)
	}
	fmt.Println()
	fmt.Println()
}

func ExecWithRetry(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
	trxId, err := Exec(ctx, api, actions)

	if err != nil {
		if !strings.Contains(err.Error(), "deadline exceeded") {
			return string(""), err
		} else {
			attempts := 1
			for attempts < 3 {
				trxId, err = Exec(ctx, api, actions)
				if err == nil {
					return trxId, nil
				}
				attempts++
			}
		}
		return string(""), err
	}
	return trxId, nil
}

func Exec(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		return string(""), fmt.Errorf("error filling tx opts: %s", err)
	}

	tx := eos.NewTransaction(actions, txOpts)

	_, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		return string(""), fmt.Errorf("error signing transaction: %s", err)
	}

	response, err := api.PushTransaction(ctx, packedTx)
	if err != nil {
		return string(""), fmt.Errorf("error pushing transaction: %s", err)
	}
	trxID := hex.EncodeToString(response.Processed.ID)
	return trxID, nil
}
