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
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper TCChanKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case prefixOrder:
			return queryOrder(ctx, path[1:], req, keeper)
		case prefixPerson:
			return queryPerson(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown tcchan query endpoint")
		}
	}
}

// nolint: unparam
func queryOrder(ctx sdk.Context, path []string, req abci.RequestQuery, keeper TCChanKeeper) ([]byte, sdk.Error) {

	orderID, err := strconv.Atoi(path[0])
	if err != nil || orderID < 0 {
		panic("order ID is not int")
	}
	order, err := keeper.GetOrder(ctx, uint64(orderID))
	if err != nil {
		panic("could not get order from local")
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, order)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

// implement fmt.Stringer
func (o CCTxOrder) String() string {
	return strings.TrimSpace(fmt.Sprintf(`OrderID: %s || AccAddress: %s || TccAddress: %s || Deposit: %t || Status: %d`, o.OrderID, o.AccAddress, o.TTCAddress, o.IsDeposit, o.Status))
}

// nolint: unparam
func queryPerson(ctx sdk.Context, path []string, req abci.RequestQuery, keeper TCChanKeeper) ([]byte, sdk.Error) {
	person, err := keeper.GetPerson(ctx, []byte(path[0]))
	if err != nil {
		panic("could not get person from local")
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, person)
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return bz, nil
}

// implement fmt.Stringer
func (p PersonalOrderRecord) String() string {
	return strings.TrimSpace(fmt.Sprintf(` AccAddress: %s || DepositIDs: %v || WithdrawIDs: %v`, p.AccAddress, p.DepositOrderIDs, p.WithdrawOrderIDs))
}
