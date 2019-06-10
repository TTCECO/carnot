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

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"math/big"
)

var (
	// errUnknownOrder is returned when order is not exist
	errUnknownOrder = errors.New("unknown order")

	// errUndefinedStatus is returned when try to set a undefined status to order
	errUndefinedStatus = errors.New("undefined status")

	// errUndefinedPrefix is returned when the prefix is undefined
	errUndefinedPrefix = errors.New("undefined prefix")
)

// TCChanKeeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type TCChanKeeper struct {
	coinKeeper    bank.Keeper
	tcchanKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc           *codec.Codec // The wire codec for binary encoding/decoding.
	operator      *Operator
	logger        log.Logger
	cliCtx        context.CLIContext
	validatorName string // validator name
	validatorPass string // validator pass
	validator     sdk.AccAddress
}

// NewTCChanKeeper creates new instances of the tcchan Keeper
func NewTCChanKeeper(logger log.Logger, coinKeeper bank.Keeper, tcchanKey sdk.StoreKey, cdc *codec.Codec,
	keyfilepath string, password string,
	validatorName string, validatorPass string, RPCPort int, keyPath string) TCChanKeeper {

	// new cli context
	cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
	cliCtx.BroadcastMode = client.BroadcastSync
	rpc := rpcclient.NewHTTP(fmt.Sprintf("tcp://127.0.0.1:%d", RPCPort), "/websocket")
	cliCtx = cliCtx.WithClient(rpc)

	// set validator
	kb, err := keys.NewKeyBaseFromHomeFlag()
	if keyPath != "" {
		kb, err = keys.NewKeyBaseFromDir(keyPath)
	}
	info, err := kb.Get(validatorName)
	if err != nil {
		panic(err)
	}
	keeper := TCChanKeeper{
		logger:        logger,
		coinKeeper:    coinKeeper,
		tcchanKey:     tcchanKey,
		cdc:           cdc,
		operator:      NewCrossChainOperator(logger, keyfilepath, password),
		cliCtx:        cliCtx,
		validatorName: validatorName,
		validatorPass: validatorPass,
		validator:     info.GetAddress(),
	}
	return keeper
}

func buildKey(input interface{}, prefix string) ([]byte, error) {
	var key []byte
	switch prefix {
	case prefixDeposit:
		key = []byte(fmt.Sprintf("%s-%s", prefixDeposit, input))
	case prefixPerson:
		key = []byte(fmt.Sprintf("%s-%s", prefixPerson, input))
	case prefixConfirm:
		key = []byte(fmt.Sprintf("%s-%s", prefixConfirm, input))
	case prefixCurrent:
		key = []byte(fmt.Sprintf("%s", prefixCurrent))
	default:
		return key, errUndefinedPrefix

	}
	return key, nil
}

// Gets the entire CCTxOrder metadata struct by OrderID
func (k TCChanKeeper) GetOrder(ctx sdk.Context, id uint64) (DepositOrder, error) {
	tmpKey, err := buildKey(id, prefixDeposit)
	if err != nil {
		return DepositOrder{}, err
	}
	store := ctx.KVStore(k.tcchanKey)
	if !store.Has(tmpKey) {
		return DepositOrder{}, errUnknownOrder
	}
	bz := store.Get(tmpKey)
	var order DepositOrder
	k.cdc.MustUnmarshalBinaryBare(bz, &order)
	return order, nil
}

// Sets the entire CCTxOrder metadata struct
func (k TCChanKeeper) SetOrder(ctx sdk.Context, order DepositOrder) error {
	tmpKey, err := buildKey(order.OrderID, prefixDeposit)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.tcchanKey)
	store.Set(tmpKey, k.cdc.MustMarshalBinaryBare(order))
	return nil
}

// Gets the entire CCTxOrder metadata struct by OrderID
func (k TCChanKeeper) GetConfirm(ctx sdk.Context, id uint64) (WithdrawConfirm, error) {
	tmpKey, err := buildKey(id, prefixConfirm)
	if err != nil {
		return WithdrawConfirm{}, err
	}
	store := ctx.KVStore(k.tcchanKey)
	if !store.Has(tmpKey) {
		return WithdrawConfirm{}, errUnknownOrder
	}
	bz := store.Get(tmpKey)
	var order WithdrawConfirm
	k.cdc.MustUnmarshalBinaryBare(bz, &order)
	return order, nil
}

// Sets the entire CCTxOrder metadata struct
func (k TCChanKeeper) SetConfirm(ctx sdk.Context, confirm WithdrawConfirm) error {
	tmpKey, err := buildKey(confirm.OrderID, prefixConfirm)
	if err != nil {
		return err
	}
	record, err := k.GetConfirm(ctx, confirm.OrderID)
	if err == nil && len(record.Confirms) > 0 && sameConfirm(record, confirm) && len(confirm.Confirms) == 1 {
		newConfirm := true
		for _, v := range record.Confirms {
			if v.Equals(confirm.Confirms[0]) {
				newConfirm = false
				break
			}
		}

		if newConfirm {
			record.Confirms = append(record.Confirms, confirm.Confirms[0])
			store := ctx.KVStore(k.tcchanKey)
			store.Set(tmpKey, k.cdc.MustMarshalBinaryBare(record))
		}
	} else {
		// todo : if some of the validator got wrong withdraw transaction , the wrong will replace the old one here
		store := ctx.KVStore(k.tcchanKey)
		store.Set(tmpKey, k.cdc.MustMarshalBinaryBare(confirm))
	}
	return nil
}

//
func (k TCChanKeeper) setConfirmSuccess(ctx sdk.Context, orderID uint64) error {
	tmpKey, err := buildKey(orderID, prefixConfirm)
	if err != nil {
		return err
	}
	record, err := k.GetConfirm(ctx, orderID)
	if err != nil {
		return err
	}
	record.Status = 1
	store := ctx.KVStore(k.tcchanKey)
	store.Set(tmpKey, k.cdc.MustMarshalBinaryBare(record))
	return nil
}

func (k TCChanKeeper) CalculateConfirm(ctx sdk.Context) error {
	iterator := k.GetRecordsIterator(ctx, prefixConfirm)
	for ; iterator.Valid(); iterator.Next() {
		var order WithdrawConfirm
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &order)
		// todo: min validator count, should related to validator count in genesis.json
		if len(order.Confirms) >= minValidatorCount && order.Status == 0 {
			if _, _, err := k.coinKeeper.AddCoins(ctx, order.AccAddress, sdk.NewCoins(order.Value)); err != nil {
				return err
			} else {
				// update the status the order to 1
				k.setConfirmSuccess(ctx, order.OrderID)
				// add the balance of target
				k.coinKeeper.AddCoins(ctx, order.AccAddress, sdk.Coins{order.Value})
			}
		}
	}
	iterator.Close()
	return nil
}

func (k TCChanKeeper) CatchWithdrawOrder(ctx sdk.Context) error {
	// todo : need find the lastID this validator already confirm for withdraw order.
	current, err := k.GetCurrent(ctx)
	if err != nil {
		return err
	}
	// get with mags from contract on ttc mainnet
	msgs, err := k.operator.GetContractWithdrawRecords(current.MaxWithdraw, blockDelay, k.validator)
	if err != nil {
		return err
	}
	if len(msgs) > 0 {
		go k.sendConfirmWith(ctx, msgs)
	} else {
		k.logger.Info("CatchWithdrawOrder", "msgs", "no new order")
	}
	return nil
}

func (k TCChanKeeper) sendConfirmWith(ctx sdk.Context, msgs []MsgWithdrawConfirm) error {
	txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(k.cdc))
	txBldr = txBldr.WithChainID(ctx.ChainID())
	if path := viper.GetString("home-client"); path != "" {
		kb, err := keys.NewKeyBaseFromDir(path)
		if err != nil {
			return err
		}
		txBldr = txBldr.WithKeybase(kb)
	}
	// get sequence of validator account
	accSeq, err := k.cliCtx.GetAccountSequence(k.validator)
	if err != nil {
		return err
	} else {
		txBldr = txBldr.WithSequence(accSeq)
	}

	// get account number
	accNum, err := k.cliCtx.GetAccountNumber(k.validator)
	if err != nil {
		return err
	} else {
		txBldr = txBldr.WithAccountNumber(accNum)
	}

	// build targetMsg for sign
	var targetMsg []sdk.Msg
	for _, msg := range msgs {
		orderID, err := strconv.Atoi(msg.OrderID)
		if err != nil {
			continue
		}
		if confirm, err := k.GetConfirm(ctx, uint64(orderID)); err != nil || confirm.Status == 1 {
			continue
		}
		if err := msg.ValidateBasic(); err != nil || !msg.Validator.Equals(k.validator) {
			continue
		}
		targetMsg = append(targetMsg, msg)
	}
	if len(targetMsg) == 0 {
		return nil
	}

	txBytes, err := txBldr.BuildAndSign(k.validatorName, k.validatorPass, targetMsg)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := k.cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}
	k.logger.Info("Confirm Withdraw Order tx", "result", res)
	return nil
}

func sameConfirm(origin, new WithdrawConfirm) bool {
	if origin.OrderID == new.OrderID && origin.Value.IsEqual(new.Value) && origin.AccAddress.Equals(new.AccAddress) {
		return true
	}
	return false
}

// GetOrderStatus - gets the order status by order id
func (k TCChanKeeper) GetOrderStatus(ctx sdk.Context, id uint64) (int, error) {
	order, err := k.GetOrder(ctx, id)
	if err != nil {
		return 0, err
	}
	return order.Status, nil

}

// SetOrderStatus - sets the current status by order id
func (k TCChanKeeper) SetOrderStatus(ctx sdk.Context, id uint64, status int) error {
	if status != OrderstatusProcess && status != OrderStatusSuccess && status != OrderStatusFail {
		return errUndefinedStatus
	}
	order, err := k.GetOrder(ctx, id)
	if err != nil {
		return err
	}
	order.Status = status
	k.SetOrder(ctx, order)
	return nil
}

// Gets the entire GetPerson metadata struct by address
func (k TCChanKeeper) GetPerson(ctx sdk.Context, address string) (PersonalOrderRecord, error) {
	tmpKey, err := buildKey(address, prefixPerson)
	if err != nil {
		return PersonalOrderRecord{}, err
	}
	store := ctx.KVStore(k.tcchanKey)
	if !store.Has(tmpKey) {
		tmpAddress, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return PersonalOrderRecord{}, err
		}
		return PersonalOrderRecord{AccAddress: tmpAddress, DepositOrderIDs: []uint64{}, WithdrawOrderIDs: []uint64{}}, nil
	}
	bz := store.Get(tmpKey)
	var person PersonalOrderRecord
	k.cdc.MustUnmarshalBinaryBare(bz, &person)
	return person, nil
}

// Sets the entire PersonalOrderRecord metadata struct
func (k TCChanKeeper) SetPerson(ctx sdk.Context, person PersonalOrderRecord) error {
	tmpKey, err := buildKey(person.AccAddress.String(), prefixPerson)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.tcchanKey)
	store.Set(tmpKey, k.cdc.MustMarshalBinaryBare(person))
	return nil
}

// Gets the entire GetPerson metadata struct by address
func (k TCChanKeeper) GetCurrent(ctx sdk.Context) (CurrentOrderRecord, error) {
	tmpKey, err := buildKey(nil, prefixCurrent)
	if err != nil {
		return CurrentOrderRecord{}, err
	}
	store := ctx.KVStore(k.tcchanKey)
	if !store.Has(tmpKey) {
		return CurrentOrderRecord{MaxDeposit: 0, MaxWithdraw: 0, Deposit: []OrderExtra{}, Withdraw: []OrderExtra{}}, nil
	}
	bz := store.Get(tmpKey)
	var current CurrentOrderRecord
	k.cdc.MustUnmarshalBinaryBare(bz, &current)
	return current, nil
}

// Sets the entire PersonalOrderRecord metadata struct
func (k TCChanKeeper) SetCurrent(ctx sdk.Context, current CurrentOrderRecord) error {
	tmpKey, err := buildKey(nil, prefixCurrent)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.tcchanKey)
	store.Set(tmpKey, k.cdc.MustMarshalBinaryBare(current))
	return nil
}

// Get an iterator over all data in which the keys by prefix
func (k TCChanKeeper) GetRecordsIterator(ctx sdk.Context, prefix string) sdk.Iterator {
	store := ctx.KVStore(k.tcchanKey)
	return sdk.KVStorePrefixIterator(store, []byte(prefix))
}

// SendConfirmTx send confirm tx to ttc
func (k TCChanKeeper) SendConfirmTx(orderID string, target string, coinName string, value *big.Int) error {
	return k.operator.SendConfirmTx(orderID, target, coinName, value)
}
