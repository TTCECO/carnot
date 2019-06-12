## Step 1. Compile Contract

compile *.sol in tcchan to get tcchan.abi and bytecode

## Step 2. Gen ABI golang file

1. save abi into carnot/x/tcchan/contract/tcchan.abi
2. ```cd $GOPATH/src/github.com/TTCECO/carnot/x/tcchan/contract```
3. ```abigen --abi tcchan.abi --pkg contract --type Contract --out contract.go```

## Step 3. Deploy Contract

1. copy bytecode into carnot/x/tcchan/contract/carnotcd/deploy.go
2. ```cd $GOPATH/src/github.com/TTCECO/carnot/```
3. ```make cd```
4. ```carnotcd```
5. set contractAddress in carnot/x/tcchan/params by contract address.

or your can deploy the contract by any way you prefer.

## Step 4. Initialize the contract setting

1. contract.addValidator(_address)
2. contract.setMinConfirmNum(_num)
3. contract.ownerChargeFund()           // for test


