# Carnot

[![Go Report Card](https://goreportcard.com/badge/github.com/TTCECO/carnot)](https://goreportcard.com/report/github.com/TTCECO/carnot)
[![Travis](https://travis-ci.org/TTCECO/carnot.svg?branch=master)](https://travis-ci.org/TTCECO/carnot)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)

TTC Foundation introduces Carnot, a cross-chain solution, which will allow the TTC on-chain assets (TTC and TST-20 tokens) to transact with Binance Chain on-chain assets (BEP-2 tokens) by interlinking two blockchains. Carnot is designed to be both (i) operable independently by a single project (or an alliance) and (ii) integratable to Binance Chain as a module.

Any blockchain platforms that support smart contract can achieve interoperability with the Binance Chain by customizing Carnot to meet its needs. We believe that Carnot will empower more blockchain platforms to join the Binance DEX.

We encourage members of global developer communities to join in and improve the functionality and stability of this open source solution.
 

## Main components of Carnot
This project is composed of the following two components:

Carnot : The side chain developed based on the Cosmos SDK (which the Binance DEX grounds itself) will interlink the TTC blockchain and Binance Chain. It also enables TTC platform to interact with the Binance Chain through the designated RPC interface and the smart contract address.

Smart Contract : Triggered by the request from Carnot, the smart contract created on the TTC platform records all the cross-chain transactions and its state on the Carnot.


## Step-by-step process of a cross-chain transaction
TTC Foundation will issue BTTC, a Binance Chain based BEP-2 token, where 1 BTTC is pegged to 1 TTC on the TTC platform. The detailed process of a cross-chain transaction from TTC to BTTC is as following:
1. Address “A” of the TTC blockchain will send TTC and a parameter pack (containing the destination address “B” of the Binance Chain) to a smart contract on TTC blockchain
2. Through an RPC interface, Carnot checks and broadcasts the deposit transaction and its state on TTC blockchain.
3. Carnot confirms whether the block height of TTC blockchain satisfied the finality conditions. Then, it changes the status of the transactions as ‘ Transfer ’ and initiates the cross-chain transaction.
4. Carnot sends BTTC to the destination address “B” of the Binance Chain.
5. Carnot confirms that the block height of Binance Chain satisfied the conditions of finality and marks the transaction status within TTC blockchain smart contract as ‘Completed’ .

vice-versa, in case of a cross-chain transaction from BTTC to TTC.

## Key considerations for the development
Here are key considerations we tackled during the development of this technology:
1. All Carnot data has to be checked along every step of the process and the balance of coin and BEP-2 token has to be re-confirmed. The cross-chain transactions will be considered ‘completed’ only when these conditions are satisfied.
2. Any transaction that fails to process within the designated block height will be considered incomplete and be reversed.
3. The cross-chain control must be designed to be capable of handling multiple accounts with different authorities at the same time. It also must be able to perform 2-way confirmation and verification between the Carnot and smart contract both ways.
4. The smart contract which carries out the cross-chain transactions would be able to withdraw or deposit coins (and/or tokens), only when certain limited conditions are met.
5. As both the Binance Chain (Tendermint) and TTC blockchain (Multi-tier BFT-DPoS) satisfy finality, this cross-chain transaction also satisfies finality.


## Building and running the example
[Click here](tutorial.md) for instructions on how to build and run the code.

