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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

const (
	OrderstatusProcess = 0
	OrderStatusSuccess = 1
	OrderStatusFail    = -1
)

// DepositOrder is the struct that contains cross chain transaction
type DepositOrder struct {
	OrderID     uint64         `json:"orderId"`
	BlockNumber uint64         `json:"blockNumber"`
	AccAddress  sdk.AccAddress `json:"accAddress"`
	TTCAddress  string         `json:"ttcAddress"`
	Value       sdk.Coin       `json:"value"`
	Status      int            `json:"status"` // -1: fail 0: processing 1: success
}

// implement fmt.Stringer
func (o DepositOrder) String() string {
	return strings.TrimSpace(fmt.Sprintf(`OrderID: %d || AccAddress: %s || TccAddress: %s ||  Status: %d`, o.OrderID, o.AccAddress, o.TTCAddress, o.Status))
}

type WithdrawConfirm struct {
	OrderID     uint64           `json:"orderId"`
	BlockNumber uint64           `json:"blockNumber"`
	AccAddress  sdk.AccAddress   `json:"accAddress"`
	TTCAddress  string           `json:"ttcAddress"`
	Value       sdk.Coin         `json:"value"`
	Status      int              `json:"status"` // -1: fail 0: processing 1: success
	Confirms    []sdk.AccAddress `json:"confirms"`
}

func (w WithdrawConfirm) String() string {
	return strings.TrimSpace(fmt.Sprintf(`OrderID: %d || AccAddress: %s || TccAddress: %s ||  Status: %d`, w.OrderID, w.AccAddress, w.TTCAddress, w.Status))
}

// PersonalOrderRecord is the struct that contains transaction related to accAddress
type PersonalOrderRecord struct {
	AccAddress       sdk.AccAddress `json:"accAddress"`
	DepositOrderIDs  []uint64       `json:"depositOrderIDs"`
	WithdrawOrderIDs []uint64       `json:"withdrawOrderIDs"`
}

// implement fmt.Stringer
func (p PersonalOrderRecord) String() string {
	return strings.TrimSpace(fmt.Sprintf(` AccAddress: %s || DepositIDs: %v || WithdrawIDs: %v`, p.AccAddress, p.DepositOrderIDs, p.WithdrawOrderIDs))
}

// OrderExtra is the struct that contains extra order info during process
type OrderExtra struct {
	OrderID uint64 `json:"orderID"` // ID of DepositOrder
	Step    int    `json:"step"`    // Process step
}

// CurrentOrderRecord is the struct that contains order not finish (CCtxOrder.Status==0)
type CurrentOrderRecord struct {
	MaxWithdraw uint64       `json:"maxWithdraw"`
	MaxDeposit  uint64       `json:"maxDeposit"`
	Deposit     []OrderExtra `json:"currentDeposit"`
	Withdraw    []OrderExtra `json:"currentWithdraw"`
}

// implement fmt.Stringer
func (c CurrentOrderRecord) String() string {
	return strings.TrimSpace(fmt.Sprintf(` MaxDeposit: %d || MaxWithdraw: %d`, c.MaxDeposit, c.MaxWithdraw))
}
