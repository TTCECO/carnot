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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	OrderstatusProcess = 0
	OrderStatusSuccess = 1
	OrderStatusFail    = -1
)

// CCTxOrder is the struct that contains cross chain transaction
type CCTxOrder struct {
	OrderID     uint64         `json:"orderId"`
	BlockNumber uint64         `json:"blockNumber"`
	AccAddress  sdk.AccAddress `json:"accAddress"`
	TTCAddress  string         `json:"ttcAddress"`
	Value       sdk.Coin       `json:"value"`
	IsDeposit   bool           `json:"deposit"` // true = deposit(cosmos to ttc)
	Status      int            `json:"status"`  // -1: fail 0: processing 1: success
}

// PersonalOrderRecord is the struct that contains transaction related to accAddress
type PersonalOrderRecord struct {
	AccAddress       sdk.AccAddress `json:"accAddress"`
	DepositOrderIDs  []uint64       `json:"depositOrderIDs"`
	WithdrawOrderIDs []uint64       `json:"withdrawOrderIDs"`
}

// OrderExtra is the struct that contains extra order info during process
type OrderExtra struct {
	OrderID uint64 `json:"orderID"` // ID of CCTxOrder
	Step    int    `json:"step"`    // Process step
}

// CurrentOrderRecord is the struct that contains order not finish (CCtxOrder.Status==0)
type CurrentOrderRecord struct {
	MaxOrderNum uint64       `json:"maxOrderNumber"`
	Deposit     []OrderExtra `json:"currentDeposit"`
	Withdraw    []OrderExtra `json:"currentWithdraw"`
}
