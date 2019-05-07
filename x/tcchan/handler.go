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
)

// NewHandler returns a handler for "tcchan" type messages.
func NewHandler(keeper TCChanKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized tcchan Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to buy name
func handleMsgDeposit(ctx sdk.Context, keeper TCChanKeeper, msg MsgDeposit) sdk.Result {

	_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.From, sdk.NewCoins(msg.Value))
	if err != nil {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}
	if err := keeper.SetOrder(ctx, CCTxOrder{OrderID: 1, AccAddress: msg.From, TTCAddress: msg.To}); err != nil {
		return sdk.ErrInsufficientCoins("Store order fail").Result()
	}
	return sdk.Result{}
}
