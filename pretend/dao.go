package pretend

import (
	"time"

	"github.com/spf13/viper"
)

func PayPeriodDuration() time.Duration {
	if viper.InConfig("PayPeriodDuration") {
		return viper.GetDuration("PayPeriodDuration")
	}
	d, _ := time.ParseDuration("5m") // default to 5h
	return d
}

func VotingPeriodDuration() time.Duration {
	if viper.InConfig("VotingPeriodDuration") {
		return viper.GetDuration("VotingPeriodDuration")
	}
	d, _ := time.ParseDuration("30s") // default to 30s
	return d
}
