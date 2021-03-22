
# Quick Start

```
git clone https://github.com/hypha-dao/envctl
cd envctl
go build
./envctl -h
```

# Create a vault file
```
./envctl vault create --import
```
Paste your private key and set your encryption password

# Erase all documents on testnet
```
./envctl erase
```

# Create the pretend environment on the testnet
```
./envctl populate pretend
```