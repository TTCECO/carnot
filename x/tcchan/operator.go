// Copyright 2019 The TTC Authors
// This file is part of the TTC library.
//
// The TTC library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The TTC library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the TTC library. If not, see <http://www.gnu.org/licenses/>.

package tcchan

//curl -H 'content-type:application/json;' -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}' http://47.111.177.215:8511
import (
	"context"
	"errors"
	"fmt"
	"github.com/TTCECO/gttc/accounts/abi/bind"
	"github.com/TTCECO/gttc/accounts/keystore"
	"github.com/TTCECO/gttc/common"
	"github.com/TTCECO/gttc/core/types"
	"github.com/TTCECO/gttc/ethclient"
	"github.com/TTCECO/gttc/rlp"
	"github.com/TTCECO/gttc/rpc"
	"github.com/TTCECO/ttc-cosmos-channal/x/tcchan/contract"
	"github.com/tendermint/tendermint/libs/log"
	"io/ioutil"
	"math/big"
	"strconv"
)

type Operator struct {
	key          *keystore.Key
	cl           *rpc.Client
	logger       log.Logger
	contractAddr common.Address
	chainID      *big.Int
	localNonce   uint64
	contract     *contract.Contract
	client       *ethclient.Client
}

func NewCrossChainOperator(logger log.Logger, keyfilepath string, password string) *Operator {
	operator := Operator{logger: logger,
		contractAddr: common.HexToAddress(contractAddress),
		chainID:      big.NewInt(defaultChainID),
	}
	// unlock account
	if keyJson, err := ioutil.ReadFile(keyfilepath); err == nil {
		if operator.key, err = keystore.DecryptKey(keyJson, password); err == nil {
			logger.Info("Account unlock success", "address", operator.key.Address.Hex())
		} else {
			logger.Error("Account unlock fail", "error", err)
		}
	} else {
		fmt.Println("Keystore load fail", "error", err)
	}
	// dial rpc
	if client, err := rpc.Dial(rpcUrl); err == nil {
		operator.cl = client
		logger.Info("Dial rpc success", "url", rpcUrl)
	} else {
		fmt.Println("Dial rpc fail", "error", err)
	}
	// update chain id
	operator.updateVersion()
	// update nonce for this account
	if nonce, err := operator.getNonce(); err != nil {
		operator.localNonce = nonce
	}
	// check balance
	if balance, err := operator.getBalance(); err != nil && balance.Cmp(big.NewInt(minBalanceValue)) > 0 {
		operator.logger.Error("Balance of this account is not enough", "balance", balance)
	}
	// init contract
	if err := operator.createContract(); err != nil {
		operator.logger.Error("Contract initialized fail", "error", err)
	}

	// go operator.tmpTestCallContract()
	if err := operator.tmpTestQueryTransaction(); err != nil {
		operator.logger.Error("Query cross chain order fail", "error", err)
	}
	return &operator
}

func (o *Operator) getNonce() (uint64, error) {
	var response string
	if err := o.cl.Call(&response, "eth_getTransactionCount", o.key.Address, "latest"); err != nil {
		return 0, err
	} else {
		nonce, err := strconv.ParseUint(response[2:], 16, 64)
		if err != nil {
			return 0, err
		}
		o.logger.Info("Current status", "nonce", nonce)
		return nonce, nil
	}
}

func (o *Operator) getBalance() (*big.Int, error) {
	var response string
	if err := o.cl.Call(&response, "eth_getBalance", o.key.Address, "latest"); err != nil {
		return nil, err
	} else {
		balance := big.NewInt(0)
		if err := balance.UnmarshalText([]byte(response)); err != nil {
			return nil, err
		}
		o.logger.Info("Current status", "balance", balance)
		return balance, nil
	}
}

func (o *Operator) updateVersion() {
	var response string
	if err := o.cl.Call(&response, "net_version"); err != nil {
	} else {
		if chainID, err := strconv.ParseUint(response, 10, 64); err != nil {
		} else {
			o.logger.Info("Current status", "chainID", chainID)
		}
	}
}

func (o *Operator) sendTransaction() error {
	if nonce, err := o.getNonce(); err == nil && nonce > o.localNonce {
		o.localNonce = nonce
	}
	tx := types.NewTransaction(o.localNonce, o.contractAddr, big.NewInt(1), uint64(100000), big.NewInt(21000000), []byte{})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(o.chainID), o.key.PrivateKey)
	if err != nil {
		o.logger.Error("Transaction sign fail", "error", err)
		return err
	}
	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		o.logger.Error("RLP Data fail", "error", err)
		return err
	}
	var response string
	if err := o.cl.Call(&response, "eth_sendRawTransaction", common.ToHex(data)); err != nil {
		o.logger.Error("Cross chain transaction Execute fail", "error", err)
		return err
	} else {
		o.logger.Info("Transaction", "result", response)
		o.localNonce += 1
	}
	return nil
}

func (o *Operator) createContract() error {
	// init contract
	o.client = ethclient.NewClient(o.cl)
	contract, err := contract.NewContract(common.HexToAddress(contractAddress), o.client)
	if err != nil {
		o.contract, o.client = nil, nil
		return err
	}
	o.contract = contract
	return nil
}

func (o *Operator) SendConfirmTx(orderID string, target string, coinName string, value *big.Int) error {
	ctx := context.Background()
	tx, err := o.contract.Confirm(bind.NewKeyedTransactor(o.key.PrivateKey), orderID, common.HexToAddress(target), coinName, value)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(ctx, o.client, tx)
	if err != nil {
		return err
	}
	o.logger.Info("Contract Confirm", "status", receipt.Status)
	return nil
}

func (o *Operator) tmpTestQueryTransaction() error {

	currentOrderID, err := o.contract.WithdrawOrderID(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if currentOrderID.Cmp(big.NewInt(0)) == 0 {
		return errors.New("order not exist")
	}

	res, err := o.contract.WithdrawRecords(&bind.CallOpts{}, currentOrderID)
	if err != nil {
		return err
	}

	o.logger.Info("last Order", "res", res)
	return nil
}

func (o *Operator) tmpTestCallContract() error {
	// test data
	testAddress := common.HexToAddress("t0c233eC8cB98133Bf202DcBAF07112C6Abb058B89")

	// init contract
	ctx := context.Background()
	o.logger.Info("Contract", "testAddress", testAddress)
	exist, err := o.contract.Validators(&bind.CallOpts{}, testAddress)
	if err != nil {
		return err
	}
	o.logger.Info("Contract Validators", "exist", exist)

	tx, err := o.contract.AddValidator(bind.NewKeyedTransactor(o.key.PrivateKey), testAddress)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(ctx, o.client, tx)
	if err != nil {
		return err
	}
	o.logger.Info("Contract AddValidator", "status", receipt.Status, "address", testAddress)

	tx, err = o.contract.AddValidator(bind.NewKeyedTransactor(o.key.PrivateKey), o.key.Address)
	if err != nil {
		return err
	}
	receipt, err = bind.WaitMined(ctx, o.client, tx)
	if err != nil {
		return err
	}
	o.logger.Info("Contract AddValidator", "status", receipt.Status, "address", o.key.Address)

	o.logger.Info("Contract", "testAddress", testAddress)
	exist, err = o.contract.Validators(&bind.CallOpts{}, testAddress)
	if err != nil {
		return err
	}
	o.logger.Info("Contract Validators", "exist", exist)

	confirmed, err := o.contract.GetConfirmStatus(&bind.CallOpts{}, "_id", o.key.Address)
	if err != nil {
		return err
	}
	o.logger.Info("Contract GetConfirmStatus", "confirmed", confirmed)

	tx, err = o.contract.Confirm(bind.NewKeyedTransactor(o.key.PrivateKey), "_id", o.key.Address, "acn", big.NewInt(1000000))
	if err != nil {
		return err
	}
	receipt, err = bind.WaitMined(ctx, o.client, tx)
	if err != nil {
		return err
	}
	o.logger.Info("Contract Confirm", "status", receipt.Status)

	confirmed, err = o.contract.GetConfirmStatus(&bind.CallOpts{}, "_id", o.key.Address)
	if err != nil {
		return err
	}
	o.logger.Info("Contract GetConfirmStatus", "confirmed", confirmed)

	o.getBalance()
	opts := bind.NewKeyedTransactor(o.key.PrivateKey)
	opts.Value = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(250))

	tx, err = o.contract.OwnerChargeFund(opts)
	receipt, err = bind.WaitMined(ctx, o.client, tx)
	if err != nil {
		return err
	}
	o.logger.Info("Contract Confirm", "status", receipt.Status)
	o.getBalance()
	opts = bind.NewKeyedTransactor(o.key.PrivateKey)
	tx, err = o.contract.OwnerWithdrawFund(opts)
	receipt, err = bind.WaitMined(ctx, o.client, tx)
	if err != nil {
		return err
	}
	o.logger.Info("Contract Confirm", "status", receipt.Status)
	o.getBalance()
	return nil
}
