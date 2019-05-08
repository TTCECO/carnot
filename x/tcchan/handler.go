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

// Handle a message to deposit, from cosmos to ttc
func handleMsgDeposit(ctx sdk.Context, keeper TCChanKeeper, msg MsgDeposit) sdk.Result {
	// get record info
	personRecord, err := keeper.GetPerson(ctx, msg.From.String())
	if err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}
	currentRecord, err := keeper.GetCurrent(ctx)
	if err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}

	// collect coin
	if _, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.From, sdk.NewCoins(msg.Value)); err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}

	// update record info
	currentRecord.MaxOrderNum += 1
	if err := keeper.SetOrder(ctx, CCTxOrder{OrderID: currentRecord.MaxOrderNum, AccAddress: msg.From, TTCAddress: msg.To}); err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}
	personRecord.DepositOrderIDs = append(personRecord.DepositOrderIDs, currentRecord.MaxOrderNum)
	if err := keeper.SetPerson(ctx, personRecord); err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}
	if err := keeper.SetCurrent(ctx, currentRecord); err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}

	return sdk.Result{}
}
