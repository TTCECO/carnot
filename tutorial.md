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

Now, you can install and run the application.

```bash
# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
tcd help
tccli help
```

## Running the live network and using the commands

To initialize configuration and a `genesis.json` file for your application and an account for the transactions, start by running:

> _*NOTE*_: In the below commands addresses are are pulled using terminal utilities. You can also just input the raw strings saved from creating keys, shown below. The commands require [`jq`](https://stedolan.github.io/jq/download/) to be installed on your machine.

> _*NOTE*_: If you have run the tutorial before, you can start from scratch with a `tcd unsafe-reset-all` or by deleting both of the home folders `rm -rf ~/.tc*`

> _*NOTE*_: If you have the Cosmos app for ledger and you want to use it, when you create the key with `tccli keys add jack` just add `--ledger` at the end. That's all you need. When you sign, `jack` will be recognized as a Ledger key and will require a device.

```bash
# Initialize configuration files and genesis file
tcd init --chain-id tctestchain

# Copy the `Address` output here and save it for later use
# [optional] add "--ledger" at the end to use a Ledger Nano S
tccli keys add jack

# Add account, with coins to the genesis file
tcd add-genesis-account $(tccli keys show jack -a) 1000cttc

# Configure your CLI to eliminate need for chain-id flag
tccli config chain-id tctestchain
tccli config output json
tccli config indent true
tccli config trust-node true
```

You can now start `tcd` by calling `tcd start`. You will see logs begin streaming that represent blocks being produced, this will take a couple of seconds.

Open another terminal to run commands against the network you have just created:

```bash
# First check the accounts to ensure they have funds
tccli query account $(tccli keys show jack -a)

# Deposit using your coins from the genesis file
tccli tx tcchan deposit t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89 50cttc --from jack

# Query the order by id
tccli query tcchan order 1
# > ...


```

### Congratulations, you have built a Cosmos SDK application! This tutorial is now complete. If you want to see how to run the same commands using the REST server [click here](run-rest.md).
