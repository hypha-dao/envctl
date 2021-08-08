package service_test

import (
	"testing"

	"github.com/hypha-dao/envctl/service"
	"gotest.tools/assert"
)

func TestCheckout(t *testing.T) {
	err := service.CheckoutRepo("/home/sebastian/test-checkout", "https://github.com/hypha-dao/dao-contracts", "master")
	assert.NilError(t, err)
}
