module github.com/hypha-dao/envctl

go 1.16

require (
	github.com/bronze1man/go-yaml2json v0.0.0-20150129175009-f6f64b738964
	github.com/dfuse-io/logging v0.0.0-20210109005628-b97a57253f70
	github.com/digital-scarcity/eos-go-test v0.0.0-20210330215538-8bd2f83c706a
	github.com/eoscanada/eos-go v0.9.1-0.20200805141443-a9d5402a7bc5
	github.com/eoscanada/eosc v1.4.0
	github.com/hypha-dao/dao-contracts/dao-go v0.0.0-20210323131703-3cf67485fa79
	github.com/hypha-dao/daoctl v0.4.2-0.20210318181659-cc9a9152a1a5
	github.com/hypha-dao/document-graph/docgraph v0.0.0-20210301235139-24626f87a02a
	github.com/k0kubun/go-ansi v0.0.0-20180517002512-3bf9e2903213
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/schollz/progressbar/v3 v3.7.6
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/tidwall/gjson v1.7.4
	github.com/tidwall/sjson v1.1.5
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
)

replace github.com/hypha-dao/dao-contracts/dao-go => ../dao-contracts/dao-go

replace github.com/hypha-dao/document-graph/docgraph => ../dao-contracts/document-graph/docgraph

replace github.com/digital-scarcity/eos-go-test => ../../digital-scarcity/eos-go-test
