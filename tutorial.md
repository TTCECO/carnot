# Building and running the carnot

## Building the `carnot` application

If you want to build the `carnot` application in this repo to see the functionalities, **Go 1.12.1+** is required .

Add some parameters to environment is necessary if you have never used the `go mod` before.

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
echo "export GO111MODULE=on" >> ~/.bash_profile
source ~/.bash_profile
```

Now You can deploy a new contract on TTC block chain.
Please follow the step on [Deploy Contract For Test](x/tcchan/contract/README.md)

When the contract is ready, you can install and run the application.

```bash
# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
carnot help
carnotcli help
```

## Running the live network and using the commands

To initialize configuration and a `genesis.json` file for your application and an account for the transactions, start by running:

> _*NOTE*_: In the below commands addresses are are pulled using terminal utilities. You can also just input the raw strings saved from creating keys, shown below. The commands require [`jq`](https://stedolan.github.io/jq/download/) to be installed on your machine.

> _*NOTE*_: If you have run the tutorial before, you can start from scratch with a `carnot unsafe-reset-all` or by deleting both of the home folders `rm -rf ~/.tc*`

> _*NOTE*_: If you have the Cosmos app for ledger and you want to use it, when you create the key with `carnotcli keys add jack` just add `--ledger` at the end. That's all you need. When you sign, `jack` will be recognized as a Ledger key and will require a device.

```bash
# Create validators
carnotcli keys add validator1 --home node1/carnotcli
carnotcli keys add validator2 --home node2/carnotcli
carnotcli keys add validator3 --home node3/carnotcli

# Initialize configuration files
carnot init  testing --chain-id=testing  --home node1/carnot --p2p.laddr tcp://0.0.0.0:26656 --rpc.laddr tcp://0.0.0.0:26657
carnot init  testing --chain-id=testing  --home node2/carnot --p2p.laddr tcp://0.0.0.0:26659 --rpc.laddr tcp://0.0.0.0:26660
carnot init  testing --chain-id=testing  --home node3/carnot --p2p.laddr tcp://0.0.0.0:26661 --rpc.laddr tcp://0.0.0.0:26662

# Add account, with coins to the genesis file
carnot add-genesis-account --home node1/carnot $(carnotcli keys show validator1 -a --home node1/carnotcli) 1000000000stake;
carnot add-genesis-account --home node1/carnot $(carnotcli keys show validator2 -a --home node2/carnotcli) 1000000000stake;
carnot add-genesis-account --home node1/carnot $(carnotcli keys show validator3 -a --home node3/carnotcli) 1000000000stake;
...

# Create transactions for set validators and write them into genesis
mkdir gentxs
carnot gentx --name validator1 --memo.port 26656 --home node1/carnot --home-client node1/carnotcli --output-document gentxs/node1.json
carnot gentx --name validator2 --memo.port 26659 --home node2/carnot --home-client node2/carnotcli --output-document gentxs/node2.json
carnot gentx --name validator3 --memo.port 26661 --home node3/carnot --home-client node3/carnotcli --output-document gentxs/node3.json
carnot collect-gentxs --home node1/carnot --gentx-dir gentxs/

# Copy the genesis file to all node dir
cp node1/carnot/config/genesis.json node2/carnot/config/
cp node1/carnot/config/genesis.json node3/carnot/config/

# Modify settings on config.toml
set allow_duplicate_ip to true // for test three nodes on one machine.
set persistent_peers by the memo in genesis file

# Start the service
carnot cc-start keyfile_1.json 1 validator1 11111111 --home node1/carnot --home-client node1/carnotcli
carnot cc-start keyfile_2.json 1 validator2 22222222 --home node2/carnot --home-client node2/carnotcli
carnot cc-start keyfile_3.json 1 validator3 33333333 --home node3/carnot --home-client node3/carnotcli

```

> _*NOTE*_: keyfile_1.json is the keyfile node used for send confirm transaction, you can use the keyfile in contract/testdata only for Test!!
> Or You can [create your TTC account](https://github.com/TTCECO/gttc/wiki/TRY-AS-SUPERNODE-ON-TESTNET#create-your-new-accountaddress-by-gttc) yourself.

```bash

# Configure your CLI to eliminate need for chain-id flag
carnotcli --home node1/carnotcli config chain-id testing
carnotcli --home node1/carnotcli config output json
carnotcli --home node1/carnotcli config indent true
carnotcli --home node1/carnotcli config trust-node true
...

# First check the accounts to ensure they have funds
carnotcli --home node1/carnotcli q account $(carnotcli keys show validator1 -a --home node1/carnotcli)

# Deposit using your coins from the genesis file
carnotcli --home node1/carnotcli tx tcchan deposit t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89 2ttc --from validator1

# Query the current status
carnotcli --home node1/carnotcli query tcchan current

# Query the order by id
carnotcli --home node1/carnotcli query tcchan order 1

```


### Congratulations, you have built a Cosmos SDK application! This tutorial is now complete. If you want to see how to run the same commands using the REST server [click here](run-rest.md).
