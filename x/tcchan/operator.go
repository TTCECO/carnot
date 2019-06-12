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
	"github.com/TTCECO/carnot/x/tcchan/contract"
	"github.com/TTCECO/gttc/accounts/abi/bind"
	"github.com/TTCECO/gttc/accounts/keystore"
	"github.com/TTCECO/gttc/common"
	"github.com/TTCECO/gttc/core/types"
	"github.com/TTCECO/gttc/ethclient"
	"github.com/TTCECO/gttc/rlp"
	"github.com/TTCECO/gttc/rpc"
	"io/ioutil"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Operator struct {
	key          *keystore.Key
	cl           *rpc.Client
	contractAddr common.Address
	chainID      *big.Int
	localNonce   uint64
	contract     *contract.Contract
	client       *ethclient.Client
}

var (
	errTTCAccountMissing = errors.New("ttc account missing")
	errCallContractFail  = errors.New("call contract fail")
)

func NewCrossChainOperator(keyfilepath string, password string) *Operator {
	operator := Operator{
		contractAddr: common.HexToAddress(ContractAddress),
		chainID:      big.NewInt(defaultChainID),
	}
	// unlock account
	if keyJson, err := ioutil.ReadFile(keyfilepath); err == nil {
		if operator.key, err = keystore.DecryptKey(keyJson, password); err == nil {
			//logger.Info("Account unlock success", "address", operator.key.Address.Hex())
		} else {
			//logger.Error("Account unlock fail", "error", err)
		}
	} else {
		//fmt.Println("Keystore load fail", "error", err)
	}
	// dial rpc
	if client, err := rpc.Dial(RPCURL); err == nil {
		operator.cl = client
		//logger.Info("Dial rpc success", "url", RPCURL)
	} else {
		//fmt.Println("Dial rpc fail", "error", err)
	}
	// update chain id
	operator.updateVersion()

	if operator.key == nil {
		return &operator
	}
	// update nonce for this account
	if nonce, err := operator.getNonce(); err != nil {
		operator.localNonce = nonce
	}
	// check balance
	if balance, err := operator.getBalance(); err != nil && balance.Cmp(big.NewInt(minBalanceValue)) > 0 {
		//operator.logger.Error("Balance of this account is not enough", "balance", balance)
	}
	// init contract
	if err := operator.createContract(); err != nil {
		//operator.logger.Error("Contract initialized fail", "error", err)
	}

	return &operator
}

func (o *Operator) getNonce() (uint64, error) {
	if o.key == nil {
		return 0, errTTCAccountMissing
	}
	var response string
	if err := o.cl.Call(&response, "eth_getTransactionCount", o.key.Address, "latest"); err != nil {
		return 0, err
	} else {
		nonce, err := strconv.ParseUint(response[2:], 16, 64)
		if err != nil {
			return 0, err
		}
		//o.logger.Info("Current status", "nonce", nonce)
		return nonce, nil
	}
}

func (o *Operator) getBalance() (*big.Int, error) {
	if o.key == nil {
		return nil, errTTCAccountMissing
	}
	var response string
	if err := o.cl.Call(&response, "eth_getBalance", o.key.Address, "latest"); err != nil {
		return nil, err
	} else {
		balance := big.NewInt(0)
		if err := balance.UnmarshalText([]byte(response)); err != nil {
			return nil, err
		}
		//o.logger.Info("Current status", "balance", balance)
		return balance, nil
	}
}

func (o *Operator) updateVersion() {
	var response string
	if err := o.cl.Call(&response, "net_version"); err != nil {
	} else {
		if chainID, err := strconv.ParseUint(response, 10, 64); err != nil {
		} else {
			o.chainID.SetUint64(chainID)
			//o.logger.Info("Current status", "chainID", chainID)
		}
	}
}

func (o *Operator) sendTransaction() error {

	if o.key == nil {
		return errTTCAccountMissing
	}
	if nonce, err := o.getNonce(); err == nil && nonce > o.localNonce {
		o.localNonce = nonce
	}
	tx := types.NewTransaction(o.localNonce, o.contractAddr, big.NewInt(1), uint64(100000), big.NewInt(21000000), []byte{})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(o.chainID), o.key.PrivateKey)
	if err != nil {
		//o.logger.Error("Transaction sign fail", "error", err)
		return err
	}
	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		//o.logger.Error("RLP Data fail", "error", err)
		return err
	}
	var response string
	if err := o.cl.Call(&response, "eth_sendRawTransaction", common.ToHex(data)); err != nil {
		//o.logger.Error("Cross chain transaction Execute fail", "error", err)
		return err
	} else {
		//o.logger.Info("Transaction", "result", response)
		o.localNonce += 1
	}
	return nil
}

func (o *Operator) createContract() error {

	if o.key == nil {
		return errTTCAccountMissing
	}
	// init contract
	o.client = ethclient.NewClient(o.cl)
	contract, err := contract.NewContract(common.HexToAddress(ContractAddress), o.client)
	if err != nil {
		o.contract, o.client = nil, nil
		return err
	}
	o.contract = contract
	return nil
}

// SendConfirmTx send tx to TTC contract to confirm this validator confirm deposit transaction on cosmos
func (o *Operator) SendConfirmTx(orderID uint64, target string, coinName string, value *big.Int) error {

	if o.key == nil {
		return errTTCAccountMissing
	}

	ctx := context.Background()
	tx, err := o.contract.ConfirmDeposit(bind.NewKeyedTransactor(o.key.PrivateKey), new(big.Int).SetUint64(orderID), common.HexToAddress(target), coinName, value)
	if err != nil {
		return err
	}
	receipt, err := bind.WaitMined(ctx, o.client, tx)
	if err != nil {
		return err
	}
	// o.logger.Info("Contract Confirm", "status", receipt.Status)
	if receipt.Status != 1 {
		return errCallContractFail
	}
	return nil
}

// GetBlockNumber return the current block number
func (o *Operator) GetBlockNumber() (*big.Int, error) {
	var response string
	if err := o.cl.Call(&response, "eth_blockNumber"); err != nil {
		return nil, errors.New("block number query fail")
	} else {
		// o.logger.Info("Contract Confirm", "status", response)
		blockNumber := big.NewInt(0)
		err = blockNumber.UnmarshalText([]byte(response))
		return blockNumber, err
	}
}

// GetContractWithdrawRecords
func (o *Operator) GetContractWithdrawRecords(ctx sdk.Context, lastID uint64, blockDelay uint64, validator sdk.AccAddress) ([]MsgWithdrawConfirm, error) {
	logger := ctx.Logger().With("module", "x/tcchan")
	currentOrderID, err := o.contract.WithdrawOrderID(&bind.CallOpts{})
	var resultMsg []MsgWithdrawConfirm
	if err != nil {
		return nil, err
	}
	if currentOrderID.Cmp(big.NewInt(0)) == 0 || currentOrderID.Cmp(new(big.Int).SetUint64(lastID)) <= 0 {
		return resultMsg, nil
	}
	currentRemoteHeight, err := o.GetBlockNumber()
	if err != nil {
		return resultMsg, nil
	}
	maxNeedConfirmHeight := currentRemoteHeight.Sub(currentRemoteHeight, new(big.Int).SetUint64(blockDelay))
	for id := lastID + 1; id <= currentOrderID.Uint64(); id++ {
		record, err := o.contract.WithdrawRecords(&bind.CallOpts{}, new(big.Int).SetUint64(id))
		if err != nil {
			continue
		}
		if record.Height.Cmp(maxNeedConfirmHeight) < 0 && record.Value.Cmp(big.NewInt(0)) > 0 {
			to, err := sdk.AccAddressFromBech32(record.Target)
			if err != nil {
				continue
			}

			amount := new(big.Int).Div(record.Value, big.NewInt(1e+18))
			logger.Info("Contract ", "Value", amount)
			resultMsg = append(resultMsg, MsgWithdrawConfirm{
				OrderID:   record.OrderID.String(),
				From:      record.Source.String(),
				To:        to,
				Value:     sdk.NewCoin(CoinTTC, sdk.NewIntFromBigInt(amount)),
				Validator: validator,
			})
		}
	}
	return resultMsg, nil
}
