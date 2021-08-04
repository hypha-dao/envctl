package pretend

import (
	"time"

	"github.com/dfuse-io/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var zlog *zap.Logger

func init() {
	logging.Register("github.com/hypha-dao/envctl/pretend", &zlog)
}

func PayPeriodDuration() time.Duration {
	if viper.InConfig("payperiodduration") {
		zlog.Debug("Found PayPeriodDuration in configuration: ", zap.String("pay-period-duration", viper.GetString("PayPeriodDuration")))
		return viper.GetDuration("PayPeriodDuration")
	}

	d, _ := time.ParseDuration("5m") // default to 5m
	zlog.Debug("Did NOT find PayPeriodDuration in configuration; using default: ", zap.String("pay-period-duration", d.String()))

	return d
}

func VotingPeriodDuration() time.Duration {
	if viper.InConfig("votingperiodduration") {
		zlog.Debug("Found VotingPeriodDuration in configuration: ", zap.String("voting-period-duration", viper.GetString("VotingPeriodDuration")))
		return viper.GetDuration("VotingPeriodDuration")
	}
	d, _ := time.ParseDuration("30s") // default to 30s
	zlog.Debug("Did NOT find VotingPeriodDuration in configuration; using default: ", zap.String("voting-period-duration", d.String()))

	return d
}
