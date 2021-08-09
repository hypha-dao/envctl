package handler

import (
	"fmt"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type InitializeOp string

const (
	InitializeOp_Start   InitializeOp = "start"
	InitializeOp_Restart InitializeOp = "restart"
	InitializeOp_Stop    InitializeOp = "stop"
	InitializeOp_Destroy InitializeOp = "start"
)

type Handler interface {
	Handle(data map[interface{}]interface{}, config map[interface{}]interface{}, initOp InitializeOp) error
}

func GetHandlers(name string, eos *service.EOS, publicKey *ecc.PublicKey) []Handler {
	switch name {
	case "accounts":
		return []Handler{
			NewAccount(eos, publicKey),
		}
	case "build-contracts":
		return []Handler{NewBuildContract()}
	case "checkout-repos":
		return []Handler{NewCheckoutRepo()}
	case "deploy":
		return []Handler{NewDeploy(eos, publicKey)}
	case "hvoice":
		return []Handler{NewHVoice(eos)}
	case "run-apps":
		return []Handler{NewRunApp("runBy=envctl")}
	default:
		panic(fmt.Sprintf("No handlers exist for name: %v", name))
	}
}
