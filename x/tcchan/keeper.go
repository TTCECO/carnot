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
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// errUnknownOrder is returned when order is not exist
	errUnknownOrder = errors.New("unknown order")

	// errUndefinedStatus is returned when try to set a undefined status to order
	errUndefinedStatus = errors.New("undefined status")
)

// TCChanKeeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type TCChanKeeper struct {
	coinKeeper bank.Keeper
	orderKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	personKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	extraKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewTCChanKeeper creates new instances of the tcchan Keeper
func NewTCChanKeeper(coinKeeper bank.Keeper, orderKey sdk.StoreKey, personKey sdk.StoreKey, extraKey sdk.StoreKey, cdc *codec.Codec) TCChanKeeper {
	return TCChanKeeper{
		coinKeeper: coinKeeper,
		orderKey:   orderKey,
		personKey:  personKey,
		extraKey:   extraKey,
		cdc:        cdc,
	}
}

// Gets the entire CCTxOrder metadata struct by OrderID
func (k TCChanKeeper) GetOrder(ctx sdk.Context, id uint64) (CCTxOrder, error) {
	tmpKey := []byte(string(id))
	store := ctx.KVStore(k.orderKey)
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
	tmpKey := []byte(string(order.OrderID))
	store := ctx.KVStore(k.orderKey)
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

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k TCChanKeeper) GetOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.orderKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
