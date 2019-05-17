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
	"encoding/json"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgDeposit defines the Deposit message, deposit the ttc from Cosmos into TTC Mainnet
type MsgDeposit struct {
	From  sdk.AccAddress
	To    string
	Value sdk.Coin
}

// NewMsgDeposit is the constructor function for MsgDeposit
func NewMsgDeposit(from sdk.AccAddress, to string, value sdk.Coin) MsgDeposit {
	return MsgDeposit{
		From:  from,
		To:    to,
		Value: value,
	}
}

// Route should return the name of the module
func (msg MsgDeposit) Route() string { return RouterName }

// Type should return the action
func (msg MsgDeposit) Type() string { return "deposit" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}
	if len(msg.To) == 0 {
		return sdk.ErrUnknownRequest("Target address cannot be empty")
	}
	if !msg.Value.IsPositive() {
		return sdk.ErrInsufficientCoins("Deposit must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeposit) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// MsgWithdrawConfirm defines the Withdraw message, deposit the ttc from TTC Mainnet to Cosmos
type MsgWithdrawConfirm struct {
	From     string
	To       sdk.AccAddress
	Value    *big.Int
	CoinName string
	OrderID  string
}

// NewMsgWithdrawConfirm is the constructor function for MsgWithdraw
func NewMsgWithdrawConfirm(from string, to sdk.AccAddress, value *big.Int, coinName string, orderID string) MsgWithdrawConfirm {
	return MsgWithdrawConfirm{
		From:     from,
		To:       to,
		Value:    value,
		CoinName: coinName,
		OrderID:  orderID,
	}
}

// Route should return the name of the module
func (msg MsgWithdrawConfirm) Route() string { return RouterName }

// Type should return the action
func (msg MsgWithdrawConfirm) Type() string { return "withdraw" }

// ValidateBasic runs stateless checks on the message
func (msg MsgWithdrawConfirm) ValidateBasic() sdk.Error {
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress(msg.To.String())
	}
	if msg.Value.Cmp(big.NewInt(0)) < 0 {
		return sdk.ErrUnknownRequest("Withdraw value must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgWithdrawConfirm) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgWithdrawConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
