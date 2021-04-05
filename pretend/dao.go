package pretend

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func PayPeriodDuration() time.Duration {
	if viper.InConfig("payperiodduration") {
		zap.S().Debug("Found PayPeriodDuration in configuration: ", viper.GetString("PayPeriodDuration"))
		return viper.GetDuration("PayPeriodDuration")
	}
	zap.S().Debug("Did NOT find PayPeriodDuration in configuration; using default of 5m")
	d, _ := time.ParseDuration("5m") // default to 5m
	return d
}

func VotingPeriodDuration() time.Duration {
	if viper.InConfig("votingperiodduration") {
		zap.S().Debug("Found VotingPeriodDuration in configuration: ", viper.GetString("VotingPeriodDuration"))
		return viper.GetDuration("VotingPeriodDuration")
	}
	zap.S().Debug("Did NOT find VotingPeriodDuration in configuration; using default of 30s")
	d, _ := time.ParseDuration("30s") // default to 30s
	return d
}
