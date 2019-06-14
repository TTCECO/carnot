# Deploy Contract For Test

>>  Do not deploy like this for Online service.

## Step 1. Compile Contract

compile *.sol in tcchan to get tcchan.abi and bytecode

## Step 2. Gen ABI golang file

1. save abi to carnot/x/tcchan/contract/tcchan.abi
2. ```cd $GOPATH/src/github.com/TTCECO/carnot/x/tcchan/contract```
3. ```abigen --abi tcchan.abi --pkg contract --type Contract --out contract.go```

contract.go should be found in ```$GOPATH/src/github.com/TTCECO/carnot/x/tcchan/contract```

## Step 3. Deploy Contract

1. set byteCode in carnot/x/tcchan/contract/carnotcd/params.go
2. ```cd $GOPATH/src/github.com/TTCECO/carnot/```
3. ```make cd```
4. ```carnotcd deploy```
5. use the contract address show in 4 to set ContractAddress in carnot/x/tcchan/params.

or your can deploy the contract by any way you prefer.

## Step 4. Initialize the contract setting

1. ```cd $GOPATH/src/github.com/TTCECO/carnot/```
2. ```make install```           // rebuild carnot & carnotcli with new contract address
3. ```make cd```                // rebuild carnotcd with new contract address
4. ```carnotcd init```

## Step 5. More test

1. ```carnotcd help```          // find out more ...



