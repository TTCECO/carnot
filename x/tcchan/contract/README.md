## Step 1. Compile Contract

compile *.sol in tcchan to get tcchan.abi and bytecode

## Step 2. Deploy Contract

1. copy bytecode into carnot/x/tcchan/contract/carnotcd/deploy.go
2. ```cd $GOPATH/src/github.com/TTCECO/carnot/
3. ```make cd```
4. ```carnotcd```
5. set contractAddress in carnot/x/tcchan/params by contract address.

or your can deploy the contract by any way you prefer.

## Step 3. Initialize the contract setting

1. contract.addValidator(_address)
2. contract.setMinConfirmNum(_num)
3. contract.ownerChargeFund()           // for test

## Step 4. Gen ABI golang file

```abigen --abi tcchan.abi --pkg contract --type Contract --out contract.go```