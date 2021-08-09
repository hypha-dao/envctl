package initialize

import (
	"fmt"

	"github.com/hypha-dao/envctl/initialize/handler"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

func Initialize(inits []interface{}, eos *service.EOS, initOp handler.InitializeOp) error {
	publicKey, err := eos.AddEOSIOKey()
	if err != nil {
		return err
	}
	for _, initI := range inits {
		init := initI.(map[interface{}]interface{})
		name := init["name"].(string)
		handlers := handler.GetHandlers(name, eos, publicKey)
		config, _ := init["config"].(map[interface{}]interface{})
		data, _ := init["data"].([]interface{})
		if data == nil {
			return fmt.Errorf("no data for init with name: %v", name)
		}
		for _, dI := range data {
			d := dI.(map[interface{}]interface{})
			for _, handler := range handlers {
				err = handler.Handle(d, config, initOp)
				if err != nil {
					return err
				}
			}
		}
		fmt.Printf("Setup of %v done.\n", name)
	}
	return nil
}
