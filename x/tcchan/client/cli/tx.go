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

package cli

import (
	"github.com/TTCECO/carnot/x/tcchan"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

// GetCmdDeposit is the CLI command for sending a deposit transaction
func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [target address] [amount]",
		Short: "cross chain deposit, from Cosmos to TTC",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			coin, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			msg := tcchan.NewMsgDeposit(cliCtx.GetFromAddress(), args[0], coin)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
}

// GetCmdWithdrawConfirm is the CLI command for sending a withdraw confirm transaction
func GetCmdWithdrawConfirm(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw [from string] [target address] [amount] [orderID int]",
		Short: "cross chain withdraw transaction confirm, from TTC to Cosmos",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			coin, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			targetAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			msg := tcchan.NewMsgWithdrawConfirm(args[0], targetAddress, coin, args[3], cliCtx.GetFromAddress())

			if err := msg.ValidateBasic() ;err != nil {
				return err
			}
			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
}
