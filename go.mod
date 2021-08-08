module github.com/hypha-dao/envctl

go 1.16

require (
	cloud.google.com/go v0.60.0 // indirect
	github.com/bronze1man/go-yaml2json v0.0.0-20150129175009-f6f64b738964
	github.com/dfuse-io/logging v0.0.0-20210109005628-b97a57253f70
	github.com/digital-scarcity/eos-go-test v0.0.0-20210402180605-db8bc7b54896
	github.com/eoscanada/eos-go v0.9.1-0.20200805141443-a9d5402a7bc5
	github.com/eoscanada/eosc v1.4.0
	github.com/go-git/go-git/v5 v5.4.2
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/hypha-dao/dao-contracts/dao-go v0.0.0-20210408021016-878ebc6d5d0d
	github.com/hypha-dao/daoctl v0.4.2-0.20210408023659-31ca5bcb62c4
	github.com/hypha-dao/document-graph/docgraph v0.0.0-20210408001022-43385207b5d0
	github.com/k0kubun/go-ansi v0.0.0-20180517002512-3bf9e2903213
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/schollz/progressbar/v3 v3.7.6
	github.com/shirou/gopsutil v3.21.7+incompatible
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/tidwall/gjson v1.7.4
	github.com/tidwall/sjson v1.1.5
	github.com/tklauser/go-sysconf v0.3.7 // indirect
	go.opencensus.io v0.22.4 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	golang.org/x/tools v0.0.0-20200806022845-90696ccdc692 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.29.0 // indirect
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/grpc v1.32.0 // indirect
	gotest.tools v2.2.0+incompatible
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

// replace github.com/hypha-dao/dao-contracts/dao-go => ../dao-contracts/dao-go

// replace github.com/hypha-dao/document-graph/docgraph => ../dao-contracts/document-graph/docgraph

// replace github.com/digital-scarcity/eos-go-test => ../../digital-scarcity/eos-go-test
