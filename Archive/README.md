README

-host

https://bsc-mainnet.nodereal.io/v1/{{Your-Api-Key-Here}}

-rpc

eth_call_balanceOf, eth_getBalance

-block

numbers of blocks, 100 means the driver will pick from latest-100 to latest in a random fashion

-token

erc20/bep20 token's hex address

-users

1> if omit, it will pull at least 10000 external addresses that have been used the above token

2> if present, it will load the external addresses from previously generated account-*-10000.csv

-qps

query per second

-dur

duration in sec

-l

log level: 1=debug; 2=Info; 3=Error, please do not use 1 for load test

build:

build on Mac: go build -ldflags="-s -w" -o build/loadtestM loadtest.go

build on Linux: go build -ldflags="-s -w" -o build/loadtestL loadtest.go

example:

query USDT balance test on BSC

build/loadtestM -host=https://bsc-mainnet.nodereal.io/v1/{{Your-Api-Key-Here}} -rpc=eth_call_balanceOf -blocks=10000 -token=0x55d398326f99059fF775485246999027B3197955 -users= -qps=5 -dur=60 -l=2

query BNB balance test on BSC

build/loadtestM -host=https://bsc-mainnet.nodereal.io/v1/{{Your-Api-Key-Here}} -rpc=eth_getBalance -blocks=10000 -token=0x55d398326f99059fF775485246999027B3197955 -users=account-19912525-10000.csv -qps=1 -dur=30 -l=1