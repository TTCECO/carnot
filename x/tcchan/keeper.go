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
	"github.com/TTCECO/gttc/accounts/keystore"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"io/ioutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	coinKeeper bank.Keeper
	tcchanKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
	key        *keystore.Key
}

// NewTCChanKeeper creates new instances of the tcchan Keeper
func NewTCChanKeeper(coinKeeper bank.Keeper, tcchanKey sdk.StoreKey, cdc *codec.Codec, keyfilepath string, password string) TCChanKeeper {
	keeper := TCChanKeeper{
		coinKeeper: coinKeeper,
		tcchanKey:  tcchanKey,
		cdc:        cdc,
		key:        nil,
	}
	// unlock account
	if keyJson, err := ioutil.ReadFile(keyfilepath); err == nil {
		if keeper.key, err = keystore.DecryptKey(keyJson, password); err == nil {
			fmt.Println("Address  unlock success : ", keeper.key.Address.Hex())
		}
	}

	return keeper
}

func buildKey(input interface{}, prefix string) ([]byte, error) {
	var key []byte
	switch prefix {
	case prefixOrder:
		key = []byte(fmt.Sprintf("%s-%s", prefixOrder, input))
	case prefixPerson:
		key = []byte(fmt.Sprintf("%s-%s", prefixPerson, input))
	case prefixCurrent:
		key = []byte(fmt.Sprintf("%s", prefixCurrent))
	default:
		return key, errUndefinedPrefix

	}
	return key, nil
}

// Gets the entire CCTxOrder metadata struct by OrderID
func (k TCChanKeeper) GetOrder(ctx sdk.Context, id uint64) (CCTxOrder, error) {
	tmpKey, err := buildKey(id, prefixOrder)
	if err != nil {
		return CCTxOrder{}, err
	}
	store := ctx.KVStore(k.tcchanKey)
	if !store.Has(tmpKey) {
		return CCTxOrder{}, errUnknownOrder
	}
	bz := store.Get(tmpKey)
	var order CCTxOrder
	k.cdc.MustUnmarshalBinaryBare(bz, &order)
	return order, nil
}

// Sets the entire CCTxOrder metadata struct
func (k TCChanKeeper) SetOrder(ctx sdk.Context, order CCTxOrder) error {
	tmpKey, err := buildKey(order.OrderID, prefixOrder)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.tcchanKey)
	store.Set(tmpKey, k.cdc.MustMarshalBinaryBare(order))
	return nil
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
		return CurrentOrderRecord{MaxOrderNum: 0, Deposit: []OrderExtra{}, Withdraw: []OrderExtra{}}, nil
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
