# Requirements

- docker
- docker-compose

# Quick Start

```
git clone https://github.com/hypha-dao/envctl
cd envctl
go build
./envctl -h
```

### Create a vault file
```
./envctl vault create --import
```
Paste your private key and set your encryption password

### Start backend services, deploy contracts and create test accounts specified in the config file
```
./envctl start
```

### Destroy and restart backend services, deploy contracts and create test accounts specified in the config file
```
./envctl start -r
```

### Stop backend services
```
./envctl stop
```

### Destroy backend services
```
./envctl stop -d
```

### Erase all documents on testnet
```
./envctl erase
```

### Create the pretend environment on the testnet
```
./envctl populate pretend
```

## Example envctl config file (init-settings specifies the contracts and accounts that the start command will create)

```
Contract: dao.hypha
DAO: dao.hypha
HusdToken: husd.hypha
HyphaToken: token.hypha
HvoiceToken: voice.hypha
Bank: bank.hypha
Events: publsh.hypha
Pause: 1s
VotingPeriodDuration: 16m
PayPeriodDuration: 1h
RootHash: 52a7ff82bd6f53b31285e97d6806d886eefb650e79754784e9d923d3df347c91
EosioEndpoint: http://localhost:8888
DAOHome: /home/vsc-workspace/hypha-dao-contracts
BackendConfigDir: /home/vsc-workspace/envctl/dho-backend-env
init-settings:
  deploy:
    base-path: /home/vsc-workspace
    contracts:
      - path: hypha-dao-contracts/build/dao/
        file-name: dao
        account: dao.hypha
      - path: hypha-dao-contracts/dao-go/artifacts/treasury
        file-name: treasury
        account: bank.hypha
      - path: hypha-dao-contracts/dao-go/artifacts/monitor
        file-name: monitor
        account: publsh.hypha
      - path: hypha-dao-contracts/dao-go/artifacts/decide
        file-name: decide
        account: trailservice
      - path: hypha-dao-contracts/dao-go/artifacts/token
        file-name: token
        account: husd.hypha
        supply: 1000000000.00 HUSD
      - path: hypha-dao-contracts/dao-go/artifacts/token
        file-name: token
        account: token.hypha
        supply: 1000000000.00 HYPHA
      - path: voice-token/build/voice
        file-name: voice
        account: voice.hypha
  accounts:
    - name: mem0.hypha
      total: 5
    - name: johnnyhypha0
      total: 1
```