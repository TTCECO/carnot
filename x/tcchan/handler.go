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

	"github.com/cosmos/cosmos-sdk/x/staking"
	"math/big"
	"strconv"
)

// NewHandler returns a handler for "tcchan" type messages.
func NewHandler(keeper TCChanKeeper, stakingKeeper staking.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgWithdrawConfirm:
			return handleMsgWithdrawConfirm(ctx, keeper, msg, stakingKeeper)
		default:
			errMsg := fmt.Sprintf("Unrecognized tcchan Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgWithdrawConfirm(ctx sdk.Context, keeper TCChanKeeper, msg MsgWithdrawConfirm, stakingKeeper staking.Keeper) sdk.Result {
	orderID, err := strconv.Atoi(msg.OrderID)
	if err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}

	isDelegation := false
	for _, delegation := range stakingKeeper.GetAllDelegations(ctx) {
		if delegation.ValidatorAddress.Equals(msg.Validator){
			isDelegation = true
		}
	}
	if !isDelegation {
		return sdk.ErrInsufficientFee("not delegation").Result()
	}

	confirm := WithdrawConfirm{
		OrderID:     uint64(orderID),
		BlockNumber: uint64(ctx.BlockHeight()),
		AccAddress:  msg.To,
		TTCAddress:  msg.From,
		Value:       msg.Value,
		Status:      0,
		Confirms:    msg.GetSigners(),
	}
	keeper.SetConfirm(ctx, confirm)

	newOrder := true
	// set current
	currentRecord, err := keeper.GetCurrent(ctx)
	if err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}

	if confirm.OrderID > currentRecord.MaxWithdraw {
		currentRecord.MaxWithdraw = confirm.OrderID
	}
	for _, order := range currentRecord.Withdraw {
		if order.OrderID == confirm.OrderID {
			newOrder = false
		}
	}
	if newOrder {
		currentRecord.Withdraw = append(currentRecord.Withdraw, OrderExtra{OrderID: confirm.OrderID, Step: 0})
		keeper.SetCurrent(ctx, currentRecord)
	}

	// set person
	if newOrder {
		personRecord, err := keeper.GetPerson(ctx, msg.To.String())
		if err != nil {
			return sdk.ErrInsufficientCoins(err.Error()).Result()
		}
		personRecord.WithdrawOrderIDs = append(personRecord.WithdrawOrderIDs, confirm.OrderID)
		keeper.SetPerson(ctx, personRecord)
	}

	return sdk.Result{}
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
	currentRecord.MaxDeposit += 1
	if err := keeper.SetOrder(ctx, DepositOrder{OrderID: currentRecord.MaxDeposit, BlockNumber: uint64(ctx.BlockHeight()), AccAddress: msg.From, TTCAddress: msg.To, Value: msg.Value}); err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}
	currentRecord.Deposit = append(currentRecord.Deposit, OrderExtra{OrderID: currentRecord.MaxDeposit, Step: 0})
	personRecord.DepositOrderIDs = append(personRecord.DepositOrderIDs, currentRecord.MaxDeposit)
	if err := keeper.SetPerson(ctx, personRecord); err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}
	if err := keeper.SetCurrent(ctx, currentRecord); err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}
	if err := keeper.SendConfirmTx(ctx, currentRecord.MaxDeposit, msg.To, msg.Value.Denom, new(big.Int).Mul(big.NewInt(decimals), msg.Value.Amount.BigInt())); err != nil {
		// todo : handle this error later , the fail tx should be record into keep and resend again next block or later during beforeBlock
		return sdk.Result{}
		// can not return err depends on outside cause err, that will break the consensus
	}
	return sdk.Result{}
}
