package handler

import (
	"fmt"

	"github.com/hypha-dao/envctl/service"
)

type CheckoutRepo struct {
}

func NewCheckoutRepo() *CheckoutRepo {
	return &CheckoutRepo{}
}

func (m *CheckoutRepo) Handle(data map[interface{}]interface{}, config map[interface{}]interface{}, initOp InitializeOp) error {
	if initOp == InitializeOp_Start || initOp == InitializeOp_Restart {
		basePath := config["base-path"].(string)
		url := data["url"].(string)
		branch := data["branch"].(string)
		fmt.Printf("Checking out repo: %v, branch: %v\n", url, branch)
		err := service.CheckoutRepo(basePath, url, branch)
		if err != nil {
			return fmt.Errorf("failed to check out repo: %v, branch: %v, error %v", url, branch, err)
		}
	}
	return nil
}
