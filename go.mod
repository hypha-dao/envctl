module github.com/hypha-dao/envctl

go 1.16

require (
	github.com/bronze1man/go-yaml2json v0.0.0-20150129175009-f6f64b738964
	github.com/dfuse-io/logging v0.0.0-20210109005628-b97a57253f70
	github.com/digital-scarcity/eos-go-test v0.0.0-20201030135239-784ff05708c0
	github.com/eoscanada/eos-go v0.9.1-0.20200805141443-a9d5402a7bc5
	github.com/eoscanada/eosc v1.4.0
	github.com/hypha-dao/dao-contracts/dao-go v0.0.0-00010101000000-000000000000
	github.com/hypha-dao/document-graph/docgraph v0.0.0-20201229193929-e09f4b1c9e47
	github.com/k0kubun/go-ansi v0.0.0-20180517002512-3bf9e2903213
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/schollz/progressbar/v3 v3.7.4
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/tidwall/gjson v1.3.2
	github.com/tidwall/sjson v1.0.4
	go.uber.org/zap v1.14.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
)

replace github.com/hypha-dao/dao-contracts/dao-go => ../dao-contracts/dao-go

replace github.com/hypha-dao/document-graph/docgraph => ../dao-contracts/document-graph/docgraph
