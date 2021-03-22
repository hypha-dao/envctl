package e

import (
	"context"
	"fmt"
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
	A           *eos.API
	X           context.Context
	AppName     string
	Contract    eos.AccountName
	TelosDecide eos.AccountName
	User        eos.AccountName
	Pause       time.Duration
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
